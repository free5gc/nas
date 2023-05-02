// Code generated by generate.sh, DO NOT EDIT.

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/free5gc/nas/nasType"
)

type PDUSessionReleaseReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONRELEASEREJECTMessageIdentity
	nasType.Cause5GSM
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionReleaseReject(iei uint8) (pDUSessionReleaseReject *PDUSessionReleaseReject) {
	pDUSessionReleaseReject = &PDUSessionReleaseReject{}
	return pDUSessionReleaseReject
}

const (
	PDUSessionReleaseRejectExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionReleaseReject) EncodePDUSessionReleaseReject(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/PDUSessionID): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/PTI): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PDUSESSIONRELEASEREJECTMessageIdentity.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/PDUSESSIONRELEASEREJECTMessageIdentity): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.Cause5GSM.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/Cause5GSM): %w", err)
	}
	if a.ExtendedProtocolConfigurationOptions != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/ExtendedProtocolConfigurationOptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/ExtendedProtocolConfigurationOptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.Buffer); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionReleaseReject/ExtendedProtocolConfigurationOptions): %w", err)
		}
	}
	return nil
}

func (a *PDUSessionReleaseReject) DecodePDUSessionReleaseReject(byteArray *[]byte) error {
	buffer := bytes.NewBuffer(*byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/PDUSessionID): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/PTI): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PDUSESSIONRELEASEREJECTMessageIdentity.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/PDUSESSIONRELEASEREJECTMessageIdentity): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.Cause5GSM.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/Cause5GSM): %w", err)
	}
	for buffer.Len() > 0 {
		var ieiN uint8
		var tmpIeiN uint8
		if err := binary.Read(buffer, binary.BigEndian, &ieiN); err != nil {
			return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/iei): %w", err)
		}
		// fmt.Println(ieiN)
		if ieiN >= 0x80 {
			tmpIeiN = (ieiN & 0xf0) >> 4
		} else {
			tmpIeiN = ieiN
		}
		// fmt.Println("type", tmpIeiN)
		switch tmpIeiN {
		case PDUSessionReleaseRejectExtendedProtocolConfigurationOptionsType:
			a.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolConfigurationOptions.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/ExtendedProtocolConfigurationOptions): %w", err)
			}
			if a.ExtendedProtocolConfigurationOptions.Len < 1 {
				return fmt.Errorf("invalid ie length (PDUSessionReleaseReject/ExtendedProtocolConfigurationOptions): %d", a.ExtendedProtocolConfigurationOptions.Len)
			}
			a.ExtendedProtocolConfigurationOptions.SetLen(a.ExtendedProtocolConfigurationOptions.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.Buffer); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionReleaseReject/ExtendedProtocolConfigurationOptions): %w", err)
			}
		default:
		}
	}
	return nil
}
