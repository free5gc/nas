package generator

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

// Generate NAS decoder and encoder in nasMessage package
func GenerateNasMessage() {
	for msgName, msgDef := range msgDefs {
		ies := msgDef.IEs
		lmsgName := strings.ToLower(msgName[0:1]) + msgName[1:]

		fOut := NewOutputFile("nasMessage/NAS_"+msgName+".go", "nasMessage", []string{
			"\"bytes\"",
			"\"encoding/binary\"",
			"\"fmt\"",
			"",
			"\"github.com/free5gc/nas/nasType\"",
		})

		// struct definition
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

		// NewXXX function
		fmt.Fprintf(fOut, "func New%s(iei uint8) (%s *%s) {\n", msgName, lmsgName, msgName)
		fmt.Fprintf(fOut, "%s = &%s{}\n", lmsgName, msgName)
		fmt.Fprintf(fOut, "return %s\n", lmsgName)
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")

		// constant definition for IEI values
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

		// Encoder
		fmt.Fprintf(fOut, "func (a *%s) Encode%s(buffer *bytes.Buffer) error {\n", msgName, msgName)
		for _, ie := range ies {
			ieType := nasTypeTable[ie.typeName]
			if ieType == nil {
				panic(fmt.Sprintf("Type %s is not exist", ie.typeName))
			}
			dataFieldName := ""
			var dateFieldType reflect.Type
			for _, n := range []string{"Buffer", "Octet"} {
				if field, exist := ieType.FieldByName(n); exist {
					dataFieldName = n
					dateFieldType = field.Type
					break
				}
			}
			// encode IEI
			if !ie.mandatory {
				fmt.Fprintf(fOut, "if a.%s != nil {\n", ie.typeName)
				if ie.iei >= 16 {
					putReadWrite(fOut, true, msgName, ie.typeName, fmt.Sprintf("a.%s.GetIei()", ie.typeName))
				}
			}
			// encode length
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
			// encode value
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

		// decoder
		fmt.Fprintf(fOut, "func (a *%s) Decode%s(byteArray *[]byte) error {\n", msgName, msgName)
		fmt.Fprintln(fOut, "buffer := bytes.NewBuffer(*byteArray)")
		for _, mandatoryPart := range []bool{true, false} {
			// parse IEI in top of optional part
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
					headLen := 0
					if !ie.mandatory {
						// allocate optional IE
						fmt.Fprintf(fOut, "case %s%sType:\n", msgName, ie.typeName)
						fmt.Fprintf(fOut, "a.%s = nasType.New%s(ieiN)\n", ie.typeName, ie.typeName)
						headLen++
					}
					if ie.lengthSize != 0 {
						// Read and check length
						putReadWrite(fOut, false, msgName, ie.typeName, fmt.Sprintf("&a.%s.Len", ie.typeName))
						headLen += ie.lengthSize
						if ie.minLength == length7or11or15 {
							fmt.Fprintf(fOut, "if a.%s.Len != %d && a.%s.Len != %d && a.%s.Len != %d {\n",
								ie.typeName, 7-headLen, ie.typeName, 11-headLen, ie.typeName, 15-headLen)
							fmt.Fprintf(fOut, "return fmt.Errorf(\"invalid ie length (%s/%s): %%d\", a.%s.Len)\n", msgName, ie.typeName, ie.typeName)
							fmt.Fprintln(fOut, "}")
						} else {
							// Calculate minimal and maximum length without type and length field
							minLength := ie.minLength - headLen
							if minLength < 0 {
								panic(fmt.Sprintf("Invalid minimal length %s/%s", msgName, ie.typeName))
							}
							maxLength := ie.maxLength - headLen
							bufMaxLength := math.MaxInt
							if dataFieldName != "" {
								// Limit value type size
								switch dateFieldType.Kind() {
								case reflect.Uint8:
									bufMaxLength = 1
								case reflect.Array:
									bufMaxLength = dateFieldType.Len()
								}
							}
							if maxLength > bufMaxLength {
								maxLength = bufMaxLength
							}
							if minLength > maxLength {
								panic(fmt.Sprintf("Invalid length %s/%s", msgName, ie.typeName))
							}
							var check []string
							// generate length check code
							if minLength == maxLength {
								check = append(check, fmt.Sprintf("a.%s.Len != %d", ie.typeName, maxLength))
							} else {
								if minLength > 0 {
									check = append(check, fmt.Sprintf("a.%s.Len < %d", ie.typeName, minLength))
								}
								typeMax := math.MaxInt8
								if ie.lengthSize == 2 {
									typeMax = math.MaxInt16
								}
								if maxLength <= typeMax {
									check = append(check, fmt.Sprintf("a.%s.Len > %d", ie.typeName, maxLength))
								}
							}
							if len(check) != 0 {
								fmt.Fprintf(fOut, "if %s {\n", strings.Join(check, " || "))
								fmt.Fprintf(fOut, "return fmt.Errorf(\"invalid ie length (%s/%s): %%d\", a.%s.Len)\n", msgName, ie.typeName, ie.typeName)
								fmt.Fprintln(fOut, "}")
							}
						}
						fmt.Fprintf(fOut, "a.%s.SetLen(a.%s.GetLen())\n", ie.typeName, ie.typeName)
					}
					// Read value
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
						} else if ie.minLength == ie.maxLength && ie.maxLength != length7or11or15 &&
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

		if err := fOut.Close(); err != nil {
			panic(err)
		}
	}
}

// Generate code for read and write with error check
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
