// Code generated by generate.sh, DO NOT EDIT.

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/free5gc/nas/nasType"
)

type Status5GSM struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.STATUSMessageIdentity5GSM
	nasType.Cause5GSM
}

func NewStatus5GSM(iei uint8) (status5GSM *Status5GSM) {
	status5GSM = &Status5GSM{}
	return status5GSM
}

func (a *Status5GSM) EncodeStatus5GSM(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS encode error (Status5GSM/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS encode error (Status5GSM/PDUSessionID): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS encode error (Status5GSM/PTI): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.STATUSMessageIdentity5GSM.Octet); err != nil {
		return fmt.Errorf("NAS encode error (Status5GSM/STATUSMessageIdentity5GSM): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.Cause5GSM.Octet); err != nil {
		return fmt.Errorf("NAS encode error (Status5GSM/Cause5GSM): %w", err)
	}
	return nil
}

func (a *Status5GSM) DecodeStatus5GSM(byteArray *[]byte) error {
	buffer := bytes.NewBuffer(*byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS decode error (Status5GSM/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS decode error (Status5GSM/PDUSessionID): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS decode error (Status5GSM/PTI): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.STATUSMessageIdentity5GSM.Octet); err != nil {
		return fmt.Errorf("NAS decode error (Status5GSM/STATUSMessageIdentity5GSM): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.Cause5GSM.Octet); err != nil {
		return fmt.Errorf("NAS decode error (Status5GSM/Cause5GSM): %w", err)
	}
	for buffer.Len() > 0 {
		var ieiN uint8
		var tmpIeiN uint8
		if err := binary.Read(buffer, binary.BigEndian, &ieiN); err != nil {
			return fmt.Errorf("NAS decode error (Status5GSM/iei): %w", err)
		}
		// fmt.Println(ieiN)
		if ieiN >= 0x80 {
			tmpIeiN = (ieiN & 0xf0) >> 4
		} else {
			tmpIeiN = ieiN
		}
		// fmt.Println("type", tmpIeiN)
		switch tmpIeiN {
		default:
		}
	}
	return nil
}
