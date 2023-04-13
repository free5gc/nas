package generator

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
)

func GenerateTestLarge() {
	largeBuf := make([]byte, 256)
	for i := 0; i < 256; i++ {
		largeBuf[i] = byte(i)
	}
	for len(largeBuf) < math.MaxUint16 {
		largeBuf = append(largeBuf, largeBuf...)
	}

	fOut := NewOutputFile("nas_generated_test.go", "nas", []string{
		"\"github.com/free5gc/nas/nasMessage\"",
		"\"github.com/free5gc/nas/nasType\"",
	})

	for _, isGMM := range []bool{true, false} {
		gmmGsm := ""
		if isGMM {
			gmmGsm = "Gmm"
		} else {
			gmmGsm = "Gsm"
		}

		fmt.Fprintf(fOut, "var tests%sMessageLarge = []struct {\n", gmmGsm)
		fmt.Fprintln(fOut, "name string")
		fmt.Fprintln(fOut, "want Message")
		fmt.Fprintln(fOut, "}{")
		for _, msgDef := range msgOrder {
			if msgDef != nil && msgDef.isGMM == isGMM && msgDef.msgName != "SecurityProtected5GSNASMessage" {
				fmt.Fprintln(fOut, "{")
				fmt.Fprintf(fOut, "name: \"%s\",\n", msgDef.msgName)
				fmt.Fprintf(fOut, "want: Message{\n")
				fmt.Fprintf(fOut, "%sMessage: &%sMessage{\n", gmmGsm, gmmGsm)
				if isGMM {
					fmt.Fprintf(fOut, "GmmHeader: GmmHeader{Octet: [3]uint8{0x7e, 0x00, 0x%02x}},\n", msgDef.msgType)
				} else {
					fmt.Fprintf(fOut, "GsmHeader: GsmHeader{Octet: [4]uint8{0x2e, 0x00, 0x00, 0x%02x}},\n", msgDef.msgType)
				}
				fmt.Fprintf(fOut, "%s: &nasMessage.%s{\n", msgDef.msgName, msgDef.msgName)
				fData, err := os.Create("testdata/" + gmmGsm + "MessageLarge/" + msgDef.msgName)
				if err != nil {
					panic(err)
				}

				for _, ie := range msgDef.IEs {
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

					ptrMark := ""
					if !ie.mandatory {
						ptrMark = "&"
					}
					fmt.Fprintf(fOut, "%s: %snasType.%s{\n", ie.typeName, ptrMark, ie.typeName)
					if !ie.mandatory && ie.maxLength != 1 {
						fmt.Fprintf(fOut, "Iei: 0x%02x,\n", ie.iei)
					}
					lenWrite := ie.maxLength
					if !ie.mandatory {
						if ie.iei < 16 {
							binary.Write(fData, binary.BigEndian, uint8(ie.iei<<4))
						} else {
							binary.Write(fData, binary.BigEndian, uint8(ie.iei))
						}
						lenWrite--
					}
					lenWrite -= ie.lengthSize
					switch ie.lengthSize {
					case 0:
					case 1:
						if lenWrite > math.MaxUint8 {
							lenWrite = math.MaxUint8
						}
						binary.Write(fData, binary.BigEndian, uint8(lenWrite))
					case 2:
						if lenWrite > math.MaxUint16 {
							lenWrite = math.MaxUint16
						}
						binary.Write(fData, binary.BigEndian, uint16(lenWrite))
					default:
						panic(fmt.Sprintf("Invalid lengthSize %d", ie.lengthSize))
					}
					if ie.lengthSize != 0 {
						fmt.Fprintf(fOut, "Len: %d,\n", lenWrite)
					}
					if strings.Contains(ie.typeName, "MessageIdentity") {
						binary.Write(fData, binary.BigEndian, msgDef.msgType)
					} else {
						switch ie.typeName {
						case "ExtendedProtocolDiscriminator":
							if isGMM {
								binary.Write(fData, binary.BigEndian, uint8(0x7E))
							} else {
								binary.Write(fData, binary.BigEndian, uint8(0x2E))
							}
						default:
							binary.Write(fData, binary.BigEndian, largeBuf[:lenWrite])
						}
					}
					if (lenWrite != 0 && dataFieldName != "") || (!ie.mandatory && ie.maxLength == 1) {
						fmt.Fprintf(fOut, "%s: ", dataFieldName)
						switch dateFieldType.Kind() {
						case reflect.Uint8:
							if strings.Contains(ie.typeName, "MessageIdentity") {
								fmt.Fprintf(fOut, "MsgType%s", msgDef.msgName)
							} else if ie.typeName == "ExtendedProtocolDiscriminator" {
								if isGMM {
									fmt.Fprintf(fOut, "nasMessage.Epd5GSMobilityManagementMessage")
								} else {
									fmt.Fprintf(fOut, "nasMessage.Epd5GSSessionManagementMessage")
								}
							} else if !ie.mandatory && ie.maxLength == 1 {
								fmt.Fprintf(fOut, "nasMessage.%s%sType", msgDef.msgName, ie.typeName)
								if ie.iei < 16 {
									fmt.Fprintf(fOut, " << 4")
								}
							} else {
								fmt.Fprintf(fOut, "0x00")
							}
						case reflect.Array:
							fmt.Fprintf(fOut, "[%d]uint8{", lenWrite)
							for i := 0; i < lenWrite; i++ {
								if i != 0 {
									fmt.Fprintf(fOut, ",")
								}
								fmt.Fprintf(fOut, "0x%02x", i)
							}
							fmt.Fprintf(fOut, "}")
						case reflect.Slice:
							fmt.Fprintf(fOut, "generateBufferSlice(%d)", lenWrite)
						}
						fmt.Fprintln(fOut, ",")
					}
					fmt.Fprintln(fOut, "},")
				}

				fmt.Fprintln(fOut, "},")
				fmt.Fprintln(fOut, "},")
				fmt.Fprintln(fOut, "},")
				fmt.Fprintln(fOut, "},")
				fData.Close()
			}
		}
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")
	}

	fOut.Close()
}
