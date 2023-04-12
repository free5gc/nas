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

type msgEntry struct {
	msgName string
	section string
	isGMM   bool
	msgType uint8
	IEs     []ieEntry
}

var msgDefs map[string]*msgEntry
var type2Msg []*msgEntry
var msgOrder []*msgEntry

func parseSpecs() {
	msgDefs = make(map[string]*msgEntry)
	type2Msg = make([]*msgEntry, 256)

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

	regMessageContent := regexp.MustCompile(`Table\pZ+(8\..+): (.*) message content`)
	regMessageType := regexp.MustCompile(`Table\pZ+9\.7\..+: Message types for (.*)`)
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
			if m := regMessageContent.FindStringSubmatch(fields[0]); m != nil {
				section := strings.Join(strings.Split(m[1], ".")[:3], ".")
				if section == "8.2.2341" {
					// XXX
					section = "8.2.24"
				}
				msgName := convertMessageName(m[2], &deregFlag)

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
				msg := &msgEntry{
					msgName: msgName,
					section: section,
					IEs:     ies,
				}
				msgDefs[msgName] = msg
				msgOrder = append(msgOrder, msg)
			} else if m := regMessageType.FindStringSubmatch(fields[0]); m != nil {
				isGMM := false
				if m[1] == "5GS mobility management" {
					isGMM = true
				}
				for {
					fields, err := csv.Read()
					if err != nil {
						panic(err)
					}
					if len(fields) == 1 {
						prevFields = fields
						break
					}
					if len(fields) == 10 {
						ok := true
						var msgType uint8
						for i := 0; i < 8; i++ {
							switch fields[i] {
							case "0":
							case "1":
								msgType += (1 << (7 - i))
							default:
								ok = false
							}
						}
						if ok {
							msgName := convertMessageName(fields[9], nil)
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
				msgName += strings.Title(strings.ToLower(word))
			}
		}
	}
	return msgName
}

func generateMsgEncDec() {
	for msgName, msgDef := range msgDefs {
		ies := msgDef.IEs
		lmsgName := strings.ToLower(msgName[0:1]) + msgName[1:]

		fOut := newOutputFile("nasMessage/NAS_"+msgName+".go", "nasMessage", []string{
			"\"bytes\"",
			"\"encoding/binary\"",
			"\"fmt\"",
			"",
			"\"github.com/free5gc/nas/nasType\"",
		})

		fmt.Fprintf(fOut, "type %s struct {\n", msgName)
		for _, ie := range ies {
			if ie.mandatory {
				fmt.Fprintf(fOut, "nasType.%s\n", ie.typeName)
			} else {
				fmt.Fprintf(fOut, "*nasType.%s\n", ie.typeName)
			}
		}
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")

		fmt.Fprintf(fOut, "func New%s(iei uint8) (%s *%s) {\n", msgName, lmsgName, msgName)
		fmt.Fprintf(fOut, "%s = &%s{}\n", lmsgName, msgName)
		fmt.Fprintf(fOut, "return %s\n", lmsgName)
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")

		hasOptionalIEs := false
		for _, ie := range ies {
			if !ie.mandatory {
				hasOptionalIEs = true
			}
		}
		if hasOptionalIEs {
			fmt.Fprintln(fOut, "const (")
			for _, ie := range ies {
				if !ie.mandatory {
					fmt.Fprintf(fOut, "%s%sType uint8 = 0x%02X\n", msgName, ie.typeName, ie.iei)
				}
			}
			fmt.Fprintln(fOut, ")")
			fmt.Fprintln(fOut, "")
		}

		fmt.Fprintf(fOut, "func (a *%s) Encode%s(buffer *bytes.Buffer) error {\n", msgName, msgName)
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
				fmt.Fprintf(fOut, "if a.%s != nil {\n", ie.typeName)
				if ie.iei >= 16 {
					putReadWrite(fOut, true, msgName, ie.typeName, fmt.Sprintf("a.%s.GetIei()", ie.typeName))
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
				putReadWrite(fOut, true, msgName, ie.typeName, fmt.Sprintf("a.%s.GetLen()", ie.typeName))
			}
			if dataFieldName == "" {
				putReadWrite(fOut, true, msgName, ie.typeName, fmt.Sprintf("&a.%s", ie.typeName))
			} else if dataFieldName == "Buffer" || dateFieldType.Kind() == reflect.Uint8 || ie.lengthSize == 0 ||
				(ie.typeName == "SessionAMBR" && ie.mandatory) || ie.typeName == "TMSI5GS" {
				putReadWrite(fOut, true, msgName, ie.typeName, fmt.Sprintf("&a.%s.%s", ie.typeName, dataFieldName))
			} else {
				putReadWrite(fOut, true, msgName, ie.typeName, fmt.Sprintf("a.%s.%s[:a.%s.GetLen()]", ie.typeName, dataFieldName, ie.typeName))
			}
			if !ie.mandatory {
				fmt.Fprintln(fOut, "}")
			}
		}
		fmt.Fprintln(fOut, "return nil")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")

		fmt.Fprintf(fOut, "func (a *%s) Decode%s(byteArray *[]byte) error {\n", msgName, msgName)
		fmt.Fprintln(fOut, "buffer := bytes.NewBuffer(*byteArray)")
		for _, mandatoryPart := range []bool{true, false} {
			if !mandatoryPart {
				fmt.Fprintln(fOut, "for buffer.Len() > 0 {")
				fmt.Fprintln(fOut, "var ieiN uint8")
				fmt.Fprintln(fOut, "var tmpIeiN uint8")
				putReadWrite(fOut, false, msgName, "iei", "&ieiN")
				fmt.Fprintln(fOut, "// fmt.Println(ieiN)")
				fmt.Fprintln(fOut, "if ieiN >= 0x80 {")
				fmt.Fprintln(fOut, "tmpIeiN = (ieiN & 0xf0) >> 4")
				fmt.Fprintln(fOut, "} else {")
				fmt.Fprintln(fOut, "tmpIeiN = ieiN")
				fmt.Fprintln(fOut, "}")
				fmt.Fprintln(fOut, "// fmt.Println(\"type\", tmpIeiN)")
				fmt.Fprintln(fOut, "switch tmpIeiN {")
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
						fmt.Fprintf(fOut, "case %s%sType:\n", msgName, ie.typeName)
						fmt.Fprintf(fOut, "a.%s = nasType.New%s(ieiN)\n", ie.typeName, ie.typeName)
					}
					if ie.lengthSize != 0 {
						putReadWrite(fOut, false, msgName, ie.typeName, fmt.Sprintf("&a.%s.Len", ie.typeName))
						fmt.Fprintf(fOut, "a.%s.SetLen(a.%s.GetLen())\n", ie.typeName, ie.typeName)
					}
					if dataFieldName == "" {
						putReadWrite(fOut, false, msgName, ie.typeName, fmt.Sprintf("&a.%s", ie.typeName))
					} else if dataFieldName == "Buffer" {
						if ie.mandatory {
							putReadWrite(fOut, false, msgName, ie.typeName, fmt.Sprintf("&a.%s.%s", ie.typeName, dataFieldName))
						} else {
							putReadWrite(fOut, false, msgName, ie.typeName, fmt.Sprintf("a.%s.%s[:a.%s.GetLen()]", ie.typeName, dataFieldName, ie.typeName))
						}
					} else {
						if dateFieldType.Kind() == reflect.Uint8 && ie.iei < 16 && !ie.mandatory {
							fmt.Fprintf(fOut, "a.%s.Octet = ieiN\n", ie.typeName)
						} else if ie.minLength == ie.maxLength &&
							ie.typeName != "IMEISV" &&
							ie.typeName != "AdditionalGUTI" &&
							ie.typeName != "GUTI5G" &&
							ie.typeName != "AuthenticationFailureParameter" &&
							ie.typeName != "AuthenticationParameterAUTN" &&
							ie.typeName != "AuthenticationResponseParameter" &&
							!(ie.typeName == "SessionAMBR" && !ie.mandatory) {
							putReadWrite(fOut, false, msgName, ie.typeName, fmt.Sprintf("&a.%s.%s", ie.typeName, dataFieldName))
						} else {
							putReadWrite(fOut, false, msgName, ie.typeName, fmt.Sprintf("a.%s.%s[:a.%s.GetLen()]", ie.typeName, dataFieldName, ie.typeName))
						}
					}
				}
			}
			if !mandatoryPart {
				fmt.Fprintln(fOut, "default:")
				fmt.Fprintln(fOut, "}")
				fmt.Fprintln(fOut, "}")
			}
		}
		fmt.Fprintln(fOut, "return nil")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")

		fOut.Close()
	}
}

func putReadWrite(f io.Writer, write bool, msgName string, ieName string, value string) {
	var rw string
	var encDec string
	if write {
		rw = "Write"
		encDec = "encode"
	} else {
		rw = "Read"
		encDec = "decode"
	}
	fmt.Fprintf(f, "if err := binary.%s(buffer, binary.BigEndian, %s) ; err != nil {\n", rw, value)
	fmt.Fprintf(f, "return fmt.Errorf(\"NAS %s error (%s/%s): %%w\", err)\n", encDec, msgName, ieName)
	fmt.Fprintln(f, "}")
}

func generateNasEncDec() {
	fOut := newOutputFile("nas_generated.go", "nas", []string{
		"\"bytes\"",
		"\"encoding/binary\"",
		"\"fmt\"",
		"",
		"\"github.com/free5gc/nas/nasMessage\"",
	})

	for _, isGMM := range []bool{true, false} {
		gmmGsm := ""
		if isGMM {
			gmmGsm = "Gmm"
		} else {
			gmmGsm = "Gsm"
		}

		fmt.Fprintf(fOut, "type %sMessage struct {\n", gmmGsm)
		fmt.Fprintf(fOut, "%sHeader\n", gmmGsm)
		for _, msgDef := range msgOrder {
			if msgDef != nil && msgDef.isGMM == isGMM {
				fmt.Fprintf(fOut, "*nasMessage.%s // %s\n", msgDef.msgName, msgDef.section)
			}
		}
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")

		fmt.Fprintln(fOut, "const (")
		for msgType, msgDef := range type2Msg {
			if msgDef != nil && msgDef.isGMM == isGMM {
				fmt.Fprintf(fOut, "MsgType%s uint8 = %d\n", msgDef.msgName, msgType)
			}
		}
		fmt.Fprintln(fOut, ")")
		fmt.Fprintln(fOut, "")

		fmt.Fprintf(fOut, "func (a *Message) %sMessageDecode(byteArray *[]byte) error {\n", gmmGsm)
		fmt.Fprintf(fOut, "buffer := bytes.NewBuffer(*byteArray)\n")
		fmt.Fprintf(fOut, "a.%sMessage = New%sMessage()\n", gmmGsm, gmmGsm)
		fmt.Fprintf(fOut, "if err := binary.Read(buffer, binary.BigEndian, &a.%sMessage.%sHeader); err != nil {\n", gmmGsm, gmmGsm)
		fmt.Fprintf(fOut, "return fmt.Errorf(\"%s NAS decode Fail: read fail - %%+v\", err)\n", strings.ToUpper(gmmGsm))
		fmt.Fprintf(fOut, "}\n")
		fmt.Fprintf(fOut, "switch a.%sMessage.%sHeader.GetMessageType() {\n", gmmGsm, gmmGsm)
		for _, msgDef := range type2Msg {
			if msgDef != nil && msgDef.isGMM == isGMM {
				fmt.Fprintf(fOut, "case MsgType%s:\n", msgDef.msgName)
				fmt.Fprintf(fOut, "a.%sMessage.%s = nasMessage.New%s(MsgType%s)\n", gmmGsm, msgDef.msgName, msgDef.msgName, msgDef.msgName)
				fmt.Fprintf(fOut, "return a.%sMessage.Decode%s(byteArray)\n", gmmGsm, msgDef.msgName)
			}
		}
		fmt.Fprintf(fOut, "default:\n")
		fmt.Fprintf(fOut, "return fmt.Errorf(\"NAS decode Fail: MsgType[%%d] doesn't exist in %s Message\",\n", strings.ToUpper(gmmGsm))
		fmt.Fprintf(fOut, "a.%sMessage.%sHeader.GetMessageType())\n", gmmGsm, gmmGsm)
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")

		fmt.Fprintf(fOut, "func (a *Message) %sMessageEncode(buffer *bytes.Buffer) error {\n", gmmGsm)
		fmt.Fprintf(fOut, "switch a.%sMessage.%sHeader.GetMessageType() {\n", gmmGsm, gmmGsm)
		for _, msgDef := range type2Msg {
			if msgDef != nil && msgDef.isGMM == isGMM {
				fmt.Fprintf(fOut, "case MsgType%s:\n", msgDef.msgName)
				fmt.Fprintf(fOut, "return a.%sMessage.Encode%s(buffer)\n", gmmGsm, msgDef.msgName)
			}
		}
		fmt.Fprintf(fOut, "default:\n")
		fmt.Fprintf(fOut, "return fmt.Errorf(\"NAS Encode Fail: MsgType[%%d] doesn't exist in %s Message\",\n", strings.ToUpper(gmmGsm))
		fmt.Fprintf(fOut, "a.%sMessage.%sHeader.GetMessageType())\n", gmmGsm, gmmGsm)
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")
	}

	fOut.Close()
}

func main() {
	parseSpecs()

	generateMsgEncDec()

	generateNasEncDec()
}

type outputFile struct {
	*bytes.Buffer
	name string
}

func newOutputFile(name string, pkgName string, imports []string) *outputFile {
	o := outputFile{
		Buffer: new(bytes.Buffer),
		name:   name,
	}

	fmt.Fprintln(o, "// Code generated by generate.sh, DO NOT EDIT.")
	fmt.Fprintln(o, "")
	fmt.Fprintf(o, "package %s\n", pkgName)
	fmt.Fprintln(o, "")
	fmt.Fprintf(o, "import (\n\n%s\n)\n", strings.Join(imports, "\n"))
	fmt.Fprintln(o, "")

	return &o
}

func (o *outputFile) Close() {
	// Output to file
	if false {
		os.Stdout.Write(o.Bytes())
		return
	}
	out, err := format.Source(o.Bytes())
	if err != nil {
		fmt.Println(string(o.Bytes()))
		panic(err)
	}
	fWrite, err := os.Create(o.name)
	if err != nil {
		panic(err)
	}
	defer fWrite.Close()
	_, err = fWrite.Write(out)
	if err != nil {
		panic(err)
	}
}
