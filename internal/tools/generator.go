//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go/format"
	"io"
	"math"
	"os"
	"reflect"
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

var msgDefs map[string][]ieEntry

func parseSpecs() {
	msgDefs = make(map[string][]ieEntry)

	fCsv, err := os.Open("spec.csv")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := fCsv.Close()
		if err != nil {
			panic(err)
		}
	}()

	regTable := regexp.MustCompile(`Table\pZ+8\..+: (.*) message content`)
	csv := csv.NewReader(fCsv)
	csv.FieldsPerRecord = -1
	deregFlag := false
	var prevFields []string
	for {
		var fields []string
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
			if m := regTable.FindStringSubmatch(fields[0]); m != nil {
				words := strings.Split(m[1], " ")
				msgName := ""
				switch m[1] {
				case "5GMM STATUS":
					msgName = "Status5GMM"
				case "5GSM STATUS":
					msgName = "Status5GSM"
				case "DEREGISTRATION REQUEST":
					if deregFlag {
						msgName = "DeregistrationRequestUETerminatedDeregistration"
					} else {
						msgName = "DeregistrationRequestUEOriginatingDeregistration"
					}
				case "DEREGISTRATION ACCEPT":
					if deregFlag {
						msgName = "DeregistrationAcceptUETerminatedDeregistration"
					} else {
						msgName = "DeregistrationAcceptUEOriginatingDeregistration"
						deregFlag = true
					}
				default:
					for _, word := range words {
						switch word {
						case "PDU", "5GMM", "5GSM", "5GS", "UL", "DL", "NAS":
							msgName += word
						default:
							msgName += strings.Title(strings.ToLower(word))
						}
					}
				}

				fields, err := csv.Read()
				if err != nil {
					panic(err)
				}
				if fields[0] != "IEI" ||
					fields[1] != "Information Element" ||
					fields[2] != "Type/Reference" ||
					fields[3] != "Presence" ||
					fields[4] != "Format" ||
					fields[5] != "Length" {
					panic("Invalid fields")
				}
				var ies []ieEntry
				prevHalf := false
			skipIE:
				for {
					fields, err := csv.Read()
					if err != nil {
						panic(err)
					}
					if len(fields) == 1 {
						prevFields = fields
						break
					}
					var ie ieEntry
					switch fields[3] {
					case "M":
						if fields[0] != "" {
							panic("IEI must be empty")
						}
						ie.mandatory = true
						switch fields[4] {
						case "V":
						case "LV":
							ie.lengthSize = 1
						case "LV-E":
							ie.lengthSize = 2
						default:
							panic(fmt.Sprintf("Invalid format %s", fields[4]))
						}
					case "O", "C":
						iei := fields[0]
						if iei == "" {
							panic("IEI must not be empty")
						}
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
						switch fields[4] {
						case "TV":
						case "TLV":
							ie.lengthSize = 1
						case "TLV-E":
							ie.lengthSize = 2
						default:
							panic(fmt.Sprintf("Invalid format %s", fields[4]))
						}
					default:
						panic(fmt.Sprintf("Invalid presence %s", fields[3]))
					}
					half := false
					lenSplit := strings.Split(fields[5], "-")
					if len(lenSplit) == 1 {
						if lenSplit[0] == "1/2" {
							half = true
						} else if lenSplit[0] == "7, 11 or 15" {
							// XXX
							ie.minLength = 7
							ie.maxLength = 15
						} else {
							if i, err := strconv.ParseInt(lenSplit[0], 10, strconv.IntSize); err != nil {
								panic(err)
							} else {
								ie.minLength = int(i)
								ie.maxLength = int(i)
							}
						}
					} else {
						if i, err := strconv.ParseInt(lenSplit[0], 10, strconv.IntSize); err != nil {
							panic(err)
						} else {
							ie.minLength = int(i)
						}
						if lenSplit[1] == "n" {
							ie.maxLength = math.MaxInt
						} else {
							if i, err := strconv.ParseInt(lenSplit[1], 10, strconv.IntSize); err != nil {
								panic(err)
							} else {
								ie.maxLength = int(i)
							}
						}
					}

					ieCell := strings.TrimSpace(fields[1])
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
						case "Non-3GPP NW policies":
							// XXX
							continue skipIE
						case "EPS bearer context status":
							// XXX
							continue skipIE
						case "5GSM congestion re-attempt indicator":
							// XXX
							continue skipIE
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
								case "NAS", "ABBA", "EAP", "TAI", "NSSAI", "LADN", "MICO", "DL", "UL", "SMS", "DNN", "TRANSPORT", "ID", "5G", "5GS", "5GSM", "5GMM", "PDU", "PTI", "SSC", "AMBR", "RQ", "EPS", "SM", "DN", "SOR", "DRX", "UE", "GUTI", "IMEISV":
									typeName += word
								case "S-NSSAI":
									typeName += "SNSSAI"
								case "Non-3GPP":
									typeName += "Non3Gpp"
								default:
									typeName += strings.Title(strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(word, "'", ""), "-", "")))
								}
							}
						}
					}
					ie.typeName = typeName
					// XXX
					if typeName == "MappedEPSBearerContexts" && msgName != "PDUSessionEstablishmentAccept" {
						ie.iei = 0x7F
					}

					if half && prevHalf {
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
				msgDefs[msgName] = ies
			}
		}
	}
}

func generate() {
	for msgName, ies := range msgDefs {
		lmsgName := strings.ToLower(msgName[0:1]) + msgName[1:]

		fBuf := new(bytes.Buffer)
		fmt.Fprint(fBuf, `
// Code generated by generate.sh, DO NOT EDIT.

package nasMessage

import (
		"bytes"
		"encoding/binary"

		"github.com/free5gc/nas/nasType"
)

`)

		fmt.Fprintf(fBuf, "type %s struct {\n", msgName)
		for _, ie := range ies {
			if ie.mandatory {
				fmt.Fprintf(fBuf, "nasType.%s\n", ie.typeName)
			} else {
				fmt.Fprintf(fBuf, "*nasType.%s\n", ie.typeName)
			}
		}
		fmt.Fprintln(fBuf, "}")
		fmt.Fprintln(fBuf, "")

		fmt.Fprintf(fBuf, "func New%s(iei uint8) (%s *%s) {\n", msgName, lmsgName, msgName)
		fmt.Fprintf(fBuf, "%s = &%s{}\n", lmsgName, msgName)
		fmt.Fprintf(fBuf, "return %s\n", lmsgName)
		fmt.Fprintln(fBuf, "}")
		fmt.Fprintln(fBuf, "")

		hasOptionalIEs := false
		for _, ie := range ies {
			if !ie.mandatory {
				hasOptionalIEs = true
			}
		}
		if hasOptionalIEs {
			fmt.Fprintln(fBuf, "const (")
			for _, ie := range ies {
				if !ie.mandatory {
					fmt.Fprintf(fBuf, "%s%sType uint8 = 0x%02X\n", msgName, ie.typeName, ie.iei)
				}
			}
			fmt.Fprintln(fBuf, ")")
			fmt.Fprintln(fBuf, "")
		}

		fmt.Fprintf(fBuf, "func (a *%s) Encode%s(buffer *bytes.Buffer) {\n", msgName, msgName)
		for _, ie := range ies {
			ieType := nasTypeTable[ie.typeName]
			dataFieldName := ""
			var dateFieldType reflect.Type
			for _, n := range []string{"Buffer", "Octet"} {
				if field, exist := ieType.FieldByName(n); exist {
					dataFieldName = n
					dateFieldType = field.Type
					break
				}
			}
			if !ie.mandatory {
				fmt.Fprintf(fBuf, "if a.%s != nil {\n", ie.typeName)
				if ie.iei >= 16 {
					fmt.Fprintf(fBuf, "binary.Write(buffer, binary.BigEndian, a.%s.GetIei())\n", ie.typeName)
				}
			}
			if ie.lengthSize != 0 {
				if lenField, exist := ieType.FieldByName("Len"); !exist {
					panic(fmt.Sprintf("Len is not exist %s", ie.typeName))
				} else {
					if lenField.Type.Size() != uintptr(ie.lengthSize) {
						panic(fmt.Sprintf("Size of length mismatch %d, %d", lenField.Type.Size(), ie.lengthSize))
					}
				}
				fmt.Fprintf(fBuf, "binary.Write(buffer, binary.BigEndian, a.%s.GetLen())\n", ie.typeName)
			}
			if dataFieldName == "" {
				fmt.Fprintf(fBuf, "binary.Write(buffer, binary.BigEndian, &a.%s)\n", ie.typeName)
			} else if dataFieldName == "Buffer" || dateFieldType.Kind() == reflect.Uint8 || ie.lengthSize == 0 ||
				(ie.typeName == "SessionAMBR" && ie.mandatory) || ie.typeName == "TMSI5GS" {
				fmt.Fprintf(fBuf, "binary.Write(buffer, binary.BigEndian, &a.%s.%s)\n", ie.typeName, dataFieldName)
			} else {
				fmt.Fprintf(fBuf, "binary.Write(buffer, binary.BigEndian, a.%s.%s[:a.%s.GetLen()])\n", ie.typeName, dataFieldName, ie.typeName)
			}
			if !ie.mandatory {
				fmt.Fprintln(fBuf, "}")
			}
		}
		fmt.Fprintln(fBuf, "}")
		fmt.Fprintln(fBuf, "")

		fmt.Fprintf(fBuf, "func (a *%s) Decode%s(byteArray *[]byte) {\n", msgName, msgName)
		fmt.Fprintln(fBuf, "buffer := bytes.NewBuffer(*byteArray)")
		for _, mandatoryPart := range []bool{true, false} {
			if !mandatoryPart {
				fmt.Fprintln(fBuf, "for buffer.Len() > 0 {")
				fmt.Fprintln(fBuf, "var ieiN uint8")
				fmt.Fprintln(fBuf, "var tmpIeiN uint8")
				fmt.Fprintln(fBuf, "binary.Read(buffer, binary.BigEndian, &ieiN)")
				fmt.Fprintln(fBuf, "// fmt.Println(ieiN)")
				fmt.Fprintln(fBuf, "if ieiN >= 0x80 {")
				fmt.Fprintln(fBuf, "tmpIeiN = (ieiN & 0xf0) >> 4")
				fmt.Fprintln(fBuf, "} else {")
				fmt.Fprintln(fBuf, "tmpIeiN = ieiN")
				fmt.Fprintln(fBuf, "}")
				fmt.Fprintln(fBuf, "// fmt.Println(\"type\", tmpIeiN)")
				fmt.Fprintln(fBuf, "switch tmpIeiN {")
			}
			for _, ie := range ies {
				ieType := nasTypeTable[ie.typeName]
				dataFieldName := ""
				var dateFieldType reflect.Type
				for _, n := range []string{"Buffer", "Octet"} {
					if field, exist := ieType.FieldByName(n); exist {
						dataFieldName = n
						dateFieldType = field.Type
						break
					}
				}

				if ie.mandatory == mandatoryPart {
					if !ie.mandatory {
						fmt.Fprintf(fBuf, "case %s%sType:\n", msgName, ie.typeName)
						fmt.Fprintf(fBuf, "a.%s = nasType.New%s(ieiN)\n", ie.typeName, ie.typeName)
					}
					if ie.lengthSize != 0 {
						fmt.Fprintf(fBuf, "binary.Read(buffer, binary.BigEndian, &a.%s.Len)\n", ie.typeName)
						fmt.Fprintf(fBuf, "a.%s.SetLen(a.%s.GetLen())\n", ie.typeName, ie.typeName)
					}
					if dataFieldName == "" {
						fmt.Fprintf(fBuf, "binary.Read(buffer, binary.BigEndian, &a.%s)\n", ie.typeName)
					} else if dataFieldName == "Buffer" {
						if ie.mandatory {
							fmt.Fprintf(fBuf, "binary.Read(buffer, binary.BigEndian, &a.%s.%s)\n", ie.typeName, dataFieldName)
						} else {
							fmt.Fprintf(fBuf, "binary.Read(buffer, binary.BigEndian, a.%s.%s[:a.%s.GetLen()])\n", ie.typeName, dataFieldName, ie.typeName)
						}
					} else {
						if dateFieldType.Kind() == reflect.Uint8 && ie.iei < 16 && !ie.mandatory {
							fmt.Fprintf(fBuf, "a.%s.Octet = ieiN\n", ie.typeName)
						} else if ie.minLength == ie.maxLength &&
							ie.typeName != "IMEISV" &&
							ie.typeName != "AdditionalGUTI" &&
							ie.typeName != "GUTI5G" &&
							ie.typeName != "AuthenticationFailureParameter" &&
							ie.typeName != "AuthenticationParameterAUTN" &&
							ie.typeName != "AuthenticationResponseParameter" &&
							!(ie.typeName == "SessionAMBR" && !ie.mandatory) {
							fmt.Fprintf(fBuf, "binary.Read(buffer, binary.BigEndian, &a.%s.%s)\n", ie.typeName, dataFieldName)
						} else {
							fmt.Fprintf(fBuf, "binary.Read(buffer, binary.BigEndian, a.%s.%s[:a.%s.GetLen()])\n", ie.typeName, dataFieldName, ie.typeName)
						}
					}
				}
			}
			if !mandatoryPart {
				fmt.Fprintln(fBuf, "default:")
				fmt.Fprintln(fBuf, "}")
				fmt.Fprintln(fBuf, "}")
			}
		}
		fmt.Fprintln(fBuf, "}")
		fmt.Fprintln(fBuf, "")

		out, err := format.Source(fBuf.Bytes())
		if err != nil {
			fmt.Println(string(fBuf.Bytes()))
			panic(err)
		}
		fOut, err := os.Create("nasMessage/NAS_" + msgName + ".go")
		if err != nil {
			panic(err)
		}
		_, err = fOut.Write(out)
		if err != nil {
			panic(err)
		}
		err = fOut.Close()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	parseSpecs()

	generate()
}
