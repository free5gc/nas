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
								ie.iei = int(i) << 4
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
							break skipIE
						case "Non-3GPP NW policies":
							break skipIE
						case "EPS bearer context status":
							break skipIE
						case "5GSM congestion re-attempt indicator":
							break skipIE
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

func main() {
	parseSpecs()

	for msgName, ies := range msgDefs {
		fBuf := new(bytes.Buffer)
		fmt.Fprint(fBuf, `
package nasMessage

import (
		// "bytes"
		// "encoding/binary"

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