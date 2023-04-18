package generator

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strings"
)

// Generate test data
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

		fmt.Fprintf(fOut, "var tests%sMessage = []struct {\n", gmmGsm)
		fmt.Fprintln(fOut, "name string")
		fmt.Fprintln(fOut, "want Message")
		fmt.Fprintln(fOut, "}{")
		for _, msgDef := range msgOrder {
			for _, isMax := range []bool{true, false} {
				minMax := ""
				if isMax {
					minMax = "Max"
				} else {
					minMax = "Min"
				}
				if msgDef != nil && msgDef.isGMM == isGMM && msgDef.structName != "SecurityProtected5GSNASMessage" {
					fmt.Fprintln(fOut, "{")
					fmt.Fprintf(fOut, "name: \"%s%s\",\n", minMax, msgDef.structName)
					fmt.Fprintf(fOut, "want: Message{\n")
					fmt.Fprintf(fOut, "%sMessage: &%sMessage{\n", gmmGsm, gmmGsm)
					if isGMM {
						fmt.Fprintf(fOut, "GmmHeader: GmmHeader{Octet: [3]uint8{0x7e, 0x00, 0x%02x}},\n", msgDef.msgType)
					} else {
						fmt.Fprintf(fOut, "GsmHeader: GsmHeader{Octet: [4]uint8{0x2e, 0x00, 0x00, 0x%02x}},\n", msgDef.msgType)
					}
					fmt.Fprintf(fOut, "%s: &nasMessage.%s{\n", msgDef.structName, msgDef.structName)
					fData, err := os.Create("testdata/" + gmmGsm + "Message/" + minMax + msgDef.structName)
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
						var lenWrite int
						if isMax {
							lenWrite = ie.maxLength
							if ie.maxLength == length7or11or15 {
								lenWrite = 15
							}
						} else {
							lenWrite = ie.minLength
							if ie.minLength == length7or11or15 {
								lenWrite = 7
							}
						}
						if !ie.mandatory {
							if ie.iei < 16 {
								err = binary.Write(fData, binary.BigEndian, uint8(ie.iei<<4))
								if err != nil {
									panic(err)
								}
							} else {
								err = binary.Write(fData, binary.BigEndian, uint8(ie.iei))
								if err != nil {
									panic(err)
								}
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
							err = binary.Write(fData, binary.BigEndian, uint8(lenWrite))
							if err != nil {
								panic(err)
							}
						case 2:
							if lenWrite > math.MaxUint16 {
								lenWrite = math.MaxUint16
							}
							err = binary.Write(fData, binary.BigEndian, uint16(lenWrite))
							if err != nil {
								panic(err)
							}
						default:
							panic(fmt.Sprintf("Invalid lengthSize %d", ie.lengthSize))
						}
						if ie.lengthSize != 0 {
							fmt.Fprintf(fOut, "Len: %d,\n", lenWrite)
						}
						if strings.Contains(ie.typeName, "MessageIdentity") {
							err = binary.Write(fData, binary.BigEndian, msgDef.msgType)
							if err != nil {
								panic(err)
							}
						} else {
							switch ie.typeName {
							case "ExtendedProtocolDiscriminator":
								if isGMM {
									err = binary.Write(fData, binary.BigEndian, uint8(0x7E))
									if err != nil {
										panic(err)
									}
								} else {
									err = binary.Write(fData, binary.BigEndian, uint8(0x2E))
									if err != nil {
										panic(err)
									}
								}
							default:
								err = binary.Write(fData, binary.BigEndian, largeBuf[:lenWrite])
								if err != nil {
									panic(err)
								}
							}
						}
						if dataFieldName != "" {
							fmt.Fprintf(fOut, "%s: ", dataFieldName)
							switch dateFieldType.Kind() {
							case reflect.Uint8:
								if strings.Contains(ie.typeName, "MessageIdentity") {
									fmt.Fprintf(fOut, "MsgType%s", msgDef.structName)
								} else if ie.typeName == "ExtendedProtocolDiscriminator" {
									if isGMM {
										fmt.Fprintf(fOut, "nasMessage.Epd5GSMobilityManagementMessage")
									} else {
										fmt.Fprintf(fOut, "nasMessage.Epd5GSSessionManagementMessage")
									}
								} else if !ie.mandatory && ie.maxLength == 1 {
									fmt.Fprintf(fOut, "nasMessage.%s%sType", msgDef.structName, ie.typeName)
									if ie.iei < 16 {
										fmt.Fprintf(fOut, " << 4")
									}
								} else {
									fmt.Fprintf(fOut, "0x00")
								}
							case reflect.Array:
								fmt.Fprintf(fOut, "[%d]uint8{", dateFieldType.Len())
								for i := 0; i < lenWrite; i++ {
									if i != 0 {
										fmt.Fprintf(fOut, ",")
									}
									fmt.Fprintf(fOut, "0x%02x", i)
								}
								fmt.Fprintf(fOut, "}")
							case reflect.Slice:
								if isMax {
									fmt.Fprintf(fOut, "generateBufferSlice(%d)", lenWrite)
								} else {
									fmt.Fprintf(fOut, "[]uint8{")
									for i := 0; i < lenWrite; i++ {
										if i != 0 {
											fmt.Fprintf(fOut, ",")
										}
										fmt.Fprintf(fOut, "0x%02x", i)
									}
									fmt.Fprintf(fOut, "}")
								}
							}
							fmt.Fprintln(fOut, ",")
						}
						fmt.Fprintln(fOut, "},")
					}

					fmt.Fprintln(fOut, "},")
					fmt.Fprintln(fOut, "},")
					fmt.Fprintln(fOut, "},")
					fmt.Fprintln(fOut, "},")

					err = fData.Close()
					if err != nil {
						panic(err)
					}

					if !isMax {
						fData, err = os.Open("testdata/" + gmmGsm + "Message/" + minMax + msgDef.structName)
						if err != nil {
							panic(err)
						}
						data, err := io.ReadAll(fData)
						if err != nil {
							panic(err)
						}
						err = fData.Close()
						if err != nil {
							panic(err)
						}
						fFuzz, err := os.Create("testdata/fuzz/Fuzz" + gmmGsm + "MessageDecode/msg" + msgDef.structName)
						if err != nil {
							panic(err)
						}
						fmt.Fprintf(fFuzz, "go test fuzz v1\n")
						fmt.Fprintf(fFuzz, "[]byte(\"")
						for _, b := range data {
							fmt.Fprintf(fFuzz, "\\x%02x", b)
						}
						fmt.Fprintf(fFuzz, "\")\n")
						if err := fFuzz.Close(); err != nil {
							panic(err)
						}
					}
				}
			}
		}
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")
	}

	if err := fOut.Close(); err != nil {
		panic(err)
	}
}
