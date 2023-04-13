// Code generated by generate.sh, DO NOT EDIT.

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/free5gc/nas/nasType"
)

type PDUSessionEstablishmentReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONESTABLISHMENTREJECTMessageIdentity
	nasType.Cause5GSM
	*nasType.BackoffTimerValue
	*nasType.AllowedSSCMode
	*nasType.EAPMessage
	*nasType.CongestionReattemptIndicator5GSM
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionEstablishmentReject(iei uint8) (pDUSessionEstablishmentReject *PDUSessionEstablishmentReject) {
	pDUSessionEstablishmentReject = &PDUSessionEstablishmentReject{}
	return pDUSessionEstablishmentReject
}

const (
	PDUSessionEstablishmentRejectBackoffTimerValueType                    uint8 = 0x37
	PDUSessionEstablishmentRejectAllowedSSCModeType                       uint8 = 0x0F
	PDUSessionEstablishmentRejectEAPMessageType                           uint8 = 0x78
	PDUSessionEstablishmentRejectCongestionReattemptIndicator5GSMType     uint8 = 0x61
	PDUSessionEstablishmentRejectExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionEstablishmentReject) EncodePDUSessionEstablishmentReject(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/PDUSessionID): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/PTI): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &a.PDUSESSIONESTABLISHMENTREJECTMessageIdentity.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/PDUSESSIONESTABLISHMENTREJECTMessageIdentity): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &a.Cause5GSM.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/Cause5GSM): %w", err)
	}
	if a.BackoffTimerValue != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.BackoffTimerValue.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/BackoffTimerValue): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.BackoffTimerValue.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/BackoffTimerValue): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, &a.BackoffTimerValue.Octet); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/BackoffTimerValue): %w", err)
		}
	}
	if a.AllowedSSCMode != nil {
		if err := binary.Write(buffer, binary.BigEndian, &a.AllowedSSCMode.Octet); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/AllowedSSCMode): %w", err)
		}
	}
	if a.EAPMessage != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.EAPMessage.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/EAPMessage): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.EAPMessage.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/EAPMessage): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, &a.EAPMessage.Buffer); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/EAPMessage): %w", err)
		}
	}
	if a.CongestionReattemptIndicator5GSM != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.CongestionReattemptIndicator5GSM.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/CongestionReattemptIndicator5GSM): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.CongestionReattemptIndicator5GSM.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/CongestionReattemptIndicator5GSM): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, &a.CongestionReattemptIndicator5GSM.Octet); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/CongestionReattemptIndicator5GSM): %w", err)
		}
	}
	if a.ExtendedProtocolConfigurationOptions != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/ExtendedProtocolConfigurationOptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/ExtendedProtocolConfigurationOptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, &a.ExtendedProtocolConfigurationOptions.Buffer); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionEstablishmentReject/ExtendedProtocolConfigurationOptions): %w", err)
		}
	}
	return nil
}

func (a *PDUSessionEstablishmentReject) DecodePDUSessionEstablishmentReject(byteArray *[]byte) error {
	buffer := bytes.NewBuffer(*byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/PDUSessionID): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/PTI): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PDUSESSIONESTABLISHMENTREJECTMessageIdentity.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/PDUSESSIONESTABLISHMENTREJECTMessageIdentity): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.Cause5GSM.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/Cause5GSM): %w", err)
	}
	for buffer.Len() > 0 {
		var ieiN uint8
		var tmpIeiN uint8
		if err := binary.Read(buffer, binary.BigEndian, &ieiN); err != nil {
			return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/iei): %w", err)
		}
		// fmt.Println(ieiN)
		if ieiN >= 0x80 {
			tmpIeiN = (ieiN & 0xf0) >> 4
		} else {
			tmpIeiN = ieiN
		}
		// fmt.Println("type", tmpIeiN)
		switch tmpIeiN {
		case PDUSessionEstablishmentRejectBackoffTimerValueType:
			a.BackoffTimerValue = nasType.NewBackoffTimerValue(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.BackoffTimerValue.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/BackoffTimerValue): %w", err)
			}
			if a.BackoffTimerValue.Len != 1 {
				return fmt.Errorf("invalid ie length (PDUSessionEstablishmentReject/BackoffTimerValue): %d", a.BackoffTimerValue.Len)
			}
			a.BackoffTimerValue.SetLen(a.BackoffTimerValue.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, &a.BackoffTimerValue.Octet); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/BackoffTimerValue): %w", err)
			}
		case PDUSessionEstablishmentRejectAllowedSSCModeType:
			a.AllowedSSCMode = nasType.NewAllowedSSCMode(ieiN)
			a.AllowedSSCMode.Octet = ieiN
		case PDUSessionEstablishmentRejectEAPMessageType:
			a.EAPMessage = nasType.NewEAPMessage(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.EAPMessage.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/EAPMessage): %w", err)
			}
			if a.EAPMessage.Len < 4 || a.EAPMessage.Len > 1500 {
				return fmt.Errorf("invalid ie length (PDUSessionEstablishmentReject/EAPMessage): %d", a.EAPMessage.Len)
			}
			a.EAPMessage.SetLen(a.EAPMessage.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.EAPMessage.Buffer[:a.EAPMessage.GetLen()]); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/EAPMessage): %w", err)
			}
		case PDUSessionEstablishmentRejectCongestionReattemptIndicator5GSMType:
			a.CongestionReattemptIndicator5GSM = nasType.NewCongestionReattemptIndicator5GSM(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.CongestionReattemptIndicator5GSM.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/CongestionReattemptIndicator5GSM): %w", err)
			}
			if a.CongestionReattemptIndicator5GSM.Len != 1 {
				return fmt.Errorf("invalid ie length (PDUSessionEstablishmentReject/CongestionReattemptIndicator5GSM): %d", a.CongestionReattemptIndicator5GSM.Len)
			}
			a.CongestionReattemptIndicator5GSM.SetLen(a.CongestionReattemptIndicator5GSM.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, &a.CongestionReattemptIndicator5GSM.Octet); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/CongestionReattemptIndicator5GSM): %w", err)
			}
		case PDUSessionEstablishmentRejectExtendedProtocolConfigurationOptionsType:
			a.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolConfigurationOptions.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/ExtendedProtocolConfigurationOptions): %w", err)
			}
			if a.ExtendedProtocolConfigurationOptions.Len < 1 {
				return fmt.Errorf("invalid ie length (PDUSessionEstablishmentReject/ExtendedProtocolConfigurationOptions): %d", a.ExtendedProtocolConfigurationOptions.Len)
			}
			a.ExtendedProtocolConfigurationOptions.SetLen(a.ExtendedProtocolConfigurationOptions.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.Buffer[:a.ExtendedProtocolConfigurationOptions.GetLen()]); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionEstablishmentReject/ExtendedProtocolConfigurationOptions): %w", err)
			}
		default:
		}
	}
	return nil
}
