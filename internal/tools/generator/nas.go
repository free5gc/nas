package generator

import (
	"fmt"
	"strings"
)

func GenerateNasEncDec() {
	fOut := NewOutputFile("nas_generated.go", "nas", []string{
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
