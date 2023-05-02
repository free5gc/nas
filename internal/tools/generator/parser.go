package generator

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ieEntry struct {
	iei        int
	typeName   string
	mandatory  bool
	lengthSize int
	minLength  int
	maxLength  int
}

const length7or11or15 = -1

type msgEntry struct {
	structName string
	section    string
	isGMM      bool
	msgType    uint8
	IEs        []ieEntry
}

var (
	msgDefs  map[string]*msgEntry
	type2Msg []*msgEntry
	msgOrder []*msgEntry
)

// Parse spec.csv and construct table of NAS messages and these IEs
func ParseSpecs() {
	msgDefs = make(map[string]*msgEntry)
	type2Msg = make([]*msgEntry, 256)

	fCsv, err := os.Open("spec.csv")
	if err != nil {
		panic(err)
	}
	defer func() {
		errClose := fCsv.Close()
		if errClose != nil {
			panic(errClose)
		}
	}()

	regMessageContent := regexp.MustCompile(`Table\pZ+(8\..+): (.*) message content`)
	regMessageType := regexp.MustCompile(`Table\pZ+9\.7\..+: Message types for (.*)`)
	csv := csv.NewReader(fCsv)
	csv.FieldsPerRecord = -1
	deregFlag := false
	var prevFields []string
	for {
		var fields []string
		// CSV parse and scan title of table
		if prevFields == nil {
			fields, err = csv.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
		} else {
			fields = prevFields
			prevFields = nil
		}
		if len(fields) == 1 {
			if m := regMessageContent.FindStringSubmatch(fields[0]); m != nil {
				// Table for message content
				sectionNumber := m[1]
				messageName := m[2]
				section := strings.Join(strings.Split(sectionNumber, ".")[:3], ".")
				if section == "8.2.2341" {
					// XXX
					section = "8.2.24"
				}
				// Convert message name to struct name
				structName := convertMessageName(messageName, &deregFlag)

				topFields, err := csv.Read()
				if err != nil {
					panic(err)
				}
				if topFields[0] != "IEI" ||
					topFields[1] != "Information Element" ||
					topFields[2] != "Type/Reference" ||
					topFields[3] != "Presence" ||
					topFields[4] != "Format" ||
					topFields[5] != "Length" {
					panic("Invalid fields")
				}
				var ies []ieEntry
				prevHalf := false
				// Read IEs in table
			skipIE:
				for {
					ieFields, err := csv.Read()
					if err != nil {
						panic(err)
					}
					if len(ieFields) == 1 {
						// End of table
						prevFields = ieFields
						break
					}

					// Parse column in table
					var ie ieEntry
					iei := ieFields[0]
					ieName := ieFields[1]
					// typeRef := fields[2]
					presence := ieFields[3]
					format := ieFields[4]
					length := ieFields[5]

					switch presence {
					case "M":
						// mandatory IE
						if iei != "" {
							panic("IEI must be empty")
						}
						ie.mandatory = true
						// parse format value
						switch format {
						case "V":
						case "LV":
							ie.lengthSize = 1
						case "LV-E":
							ie.lengthSize = 2
						default:
							panic(fmt.Sprintf("Invalid format %s", format))
						}
					case "O", "C":
						// not mandatory IE
						if iei == "" {
							panic("IEI must not be empty")
						}
						// parse IEI value
						if len(iei) > 1 && iei[1] == '-' {
							if i, err := strconv.ParseUint(iei[0:1], 16, 4); err != nil {
								panic(err)
							} else {
								ie.iei = int(i)
							}
						} else {
							if i, err := strconv.ParseUint(iei, 16, 8); err != nil {
								panic(err)
							} else {
								ie.iei = int(i)
							}
						}
						// parse format value
						switch format {
						case "TV":
						case "TLV":
							ie.lengthSize = 1
						case "TLV-E":
							ie.lengthSize = 2
						default:
							panic(fmt.Sprintf("Invalid format %s", format))
						}
					default:
						panic(fmt.Sprintf("Invalid presence %s", presence))
					}

					// parse length field
					half := false
					lenSplit := strings.Split(length, "-")
					if len(lenSplit) == 1 {
						// Fixed length
						if lenSplit[0] == "1/2" {
							// half octet IE
							half = true
						} else if lenSplit[0] == "7, 11 or 15" {
							// Special case (PDU address IE)
							ie.minLength = length7or11or15
							ie.maxLength = length7or11or15
						} else {
							if i, err := strconv.ParseInt(lenSplit[0], 10, strconv.IntSize); err != nil {
								panic(err)
							} else {
								ie.minLength = int(i)
								ie.maxLength = int(i)
							}
						}
					} else {
						// Length range
						if i, err := strconv.ParseInt(lenSplit[0], 10, strconv.IntSize); err != nil {
							panic(err)
						} else {
							ie.minLength = int(i)
						}
						if lenSplit[1] == "n" {
							// length is not limited
							ie.maxLength = math.MaxInt
						} else {
							if i, err := strconv.ParseInt(lenSplit[1], 10, strconv.IntSize); err != nil {
								panic(err)
							} else {
								ie.maxLength = int(i)
							}
						}
					}

					// Convert IE name text to go type for IE
					ieCell := strings.TrimSpace(ieName)
					words := strings.Split(ieCell, " ")
					typeName := ""
					if words[0][0] == '5' {
						words = append(words[1:], words[0])
					}
					if strings.HasPrefix(ieCell, "PDU SESSION ") && strings.HasSuffix(ieCell, " message identity") {
						typeName = strings.ReplaceAll(strings.ReplaceAll(ieCell, " ", ""), "messageidentity", "MessageIdentity")
					} else {
						switch ieCell {
						case "Payload container type":
							if ie.mandatory {
								typeName = "PayloadContainerType"
							} else {
								continue skipIE
							}
						case "5GMM STATUS message identity":
							typeName = "STATUSMessageIdentity5GMM"
						case "5GSM STATUS message identity":
							typeName = "STATUSMessageIdentity5GSM"
						case "5G-GUTI":
							typeName = "GUTI5G"
						case "5G-S-TMSI":
							typeName = "TMSI5GS"
						case "Authentication parameter RAND (5G authentication challenge)":
							typeName = "AuthenticationParameterRAND"
						case "Authentication parameter AUTN (5G authentication challenge)":
							typeName = "AuthenticationParameterAUTN"
						case "PDU session ID":
							if ie.iei == 0x12 {
								typeName = "PduSessionID2Value"
							} else {
								typeName = "PDUSessionID"
							}
						default:
							for _, word := range words {
								switch word {
								case "NAS", "ABBA", "EAP", "TAI", "NSSAI", "LADN", "MICO", "DL", "UL", "SMS",
									"DNN", "TRANSPORT", "ID", "5G", "5GS", "5GSM", "5GMM", "PDU", "PTI", "SSC",
									"AMBR", "RQ", "EPS", "SM", "DN", "SOR", "DRX", "UE", "GUTI", "IMEISV":
									typeName += word
								case "S-NSSAI":
									typeName += "SNSSAI"
								case "Non-3GPP":
									typeName += "Non3Gpp"
								default:
									typeName += titleCase(strings.ReplaceAll(strings.ReplaceAll(word, "'", ""), "-", ""))
								}
							}
						}
					}
					ie.typeName = typeName

					if half && prevHalf {
						// Merge IEs have half octet size
						prevIe := &ies[len(ies)-1]
						if prevIe.minLength != 0 {
							panic("Merge non half IEs")
						}
						if !prevIe.mandatory || !ie.mandatory {
							panic("Merge non mandatory IEs")
						}
						if prevIe.lengthSize != 0 || ie.lengthSize != 0 {
							panic("Merge IEs has length")
						}
						prevIe.typeName = ie.typeName + "And" + prevIe.typeName
						prevIe.minLength = 1
						prevIe.maxLength = 1
						prevHalf = false
					} else {
						ies = append(ies, ie)
						prevHalf = half
					}
				}
				msg := &msgEntry{
					structName: structName,
					section:    section,
					IEs:        ies,
				}
				msgDefs[structName] = msg
				msgOrder = append(msgOrder, msg)
			} else if m := regMessageType.FindStringSubmatch(fields[0]); m != nil {
				// Parse message ID table
				isGMM := false
				if m[1] == "5GS mobility management" {
					isGMM = true
				}
				for {
					idFields, err := csv.Read()
					if err != nil {
						panic(err)
					}
					if len(idFields) == 1 {
						prevFields = idFields
						break
					}
					if len(idFields) == 10 {
						ok := true
						var msgType uint8
						for i := 0; i < 8; i++ {
							switch idFields[i] {
							case "0":
							case "1":
								msgType += (1 << (7 - i))
							default:
								ok = false
							}
						}
						if ok {
							msgName := convertMessageName(idFields[9], nil)
							msgDefs[msgName].isGMM = isGMM
							msgDefs[msgName].msgType = msgType
							type2Msg[msgType] = msgDefs[msgName]
						}
					}
				}
			}
		}
	}

	// XXX
	msgDefs["SecurityProtected5GSNASMessage"].isGMM = true
}

// covert message name in TS document to struct name in generated code
func convertMessageName(msgNameInDoc string, deregFlag *bool) string {
	words := strings.Split(msgNameInDoc, " ")
	msgName := ""
	switch msgNameInDoc {
	case "5GMM STATUS", "5GMM status":
		msgName = "Status5GMM"
	case "5GSM STATUS", "5GSM status":
		msgName = "Status5GSM"
	case "DEREGISTRATION REQUEST":
		if *deregFlag {
			msgName = "DeregistrationRequestUETerminatedDeregistration"
		} else {
			msgName = "DeregistrationRequestUEOriginatingDeregistration"
		}
	case "Deregistration request (UE terminated)":
		msgName = "DeregistrationRequestUETerminatedDeregistration"
	case "Deregistration request (UE originating)":
		msgName = "DeregistrationRequestUEOriginatingDeregistration"
	case "DEREGISTRATION ACCEPT":
		if *deregFlag {
			msgName = "DeregistrationAcceptUETerminatedDeregistration"
		} else {
			msgName = "DeregistrationAcceptUEOriginatingDeregistration"
			*deregFlag = true
		}
	case "Deregistration accept (UE terminated)":
		msgName = "DeregistrationAcceptUETerminatedDeregistration"
	case "Deregistration accept (UE originating)":
		msgName = "DeregistrationAcceptUEOriginatingDeregistration"
	default:
		for _, word := range words {
			switch word {
			case "PDU", "5GMM", "5GSM", "5GS", "UL", "DL", "NAS":
				msgName += word
			default:
				msgName += titleCase(word)
			}
		}
	}
	return msgName
}

func titleCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
