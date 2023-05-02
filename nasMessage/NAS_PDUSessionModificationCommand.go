// Code generated by generate.sh, DO NOT EDIT.

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/free5gc/nas/nasType"
)

type PDUSessionModificationCommand struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONMODIFICATIONCOMMANDMessageIdentity
	*nasType.Cause5GSM
	*nasType.SessionAMBR
	*nasType.RQTimerValue
	*nasType.AlwaysonPDUSessionIndication
	*nasType.AuthorizedQosRules
	*nasType.MappedEPSBearerContexts
	*nasType.AuthorizedQosFlowDescriptions
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionModificationCommand(iei uint8) (pDUSessionModificationCommand *PDUSessionModificationCommand) {
	pDUSessionModificationCommand = &PDUSessionModificationCommand{}
	return pDUSessionModificationCommand
}

const (
	PDUSessionModificationCommandCause5GSMType                            uint8 = 0x59
	PDUSessionModificationCommandSessionAMBRType                          uint8 = 0x2A
	PDUSessionModificationCommandRQTimerValueType                         uint8 = 0x56
	PDUSessionModificationCommandAlwaysonPDUSessionIndicationType         uint8 = 0x08
	PDUSessionModificationCommandAuthorizedQosRulesType                   uint8 = 0x7A
	PDUSessionModificationCommandMappedEPSBearerContextsType              uint8 = 0x75
	PDUSessionModificationCommandAuthorizedQosFlowDescriptionsType        uint8 = 0x79
	PDUSessionModificationCommandExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionModificationCommand) EncodePDUSessionModificationCommand(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/PDUSessionID): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/PTI): %w", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, a.PDUSESSIONMODIFICATIONCOMMANDMessageIdentity.Octet); err != nil {
		return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/PDUSESSIONMODIFICATIONCOMMANDMessageIdentity): %w", err)
	}
	if a.Cause5GSM != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.Cause5GSM.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/Cause5GSM): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.Cause5GSM.Octet); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/Cause5GSM): %w", err)
		}
	}
	if a.SessionAMBR != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.SessionAMBR.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/SessionAMBR): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.SessionAMBR.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/SessionAMBR): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.SessionAMBR.Octet[:]); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/SessionAMBR): %w", err)
		}
	}
	if a.RQTimerValue != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.RQTimerValue.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/RQTimerValue): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.RQTimerValue.Octet); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/RQTimerValue): %w", err)
		}
	}
	if a.AlwaysonPDUSessionIndication != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.AlwaysonPDUSessionIndication.Octet); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/AlwaysonPDUSessionIndication): %w", err)
		}
	}
	if a.AuthorizedQosRules != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.AuthorizedQosRules.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/AuthorizedQosRules): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.AuthorizedQosRules.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/AuthorizedQosRules): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.AuthorizedQosRules.Buffer); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/AuthorizedQosRules): %w", err)
		}
	}
	if a.MappedEPSBearerContexts != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.MappedEPSBearerContexts.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/MappedEPSBearerContexts): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.MappedEPSBearerContexts.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/MappedEPSBearerContexts): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.MappedEPSBearerContexts.Buffer); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/MappedEPSBearerContexts): %w", err)
		}
	}
	if a.AuthorizedQosFlowDescriptions != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.AuthorizedQosFlowDescriptions.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/AuthorizedQosFlowDescriptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.AuthorizedQosFlowDescriptions.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/AuthorizedQosFlowDescriptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.AuthorizedQosFlowDescriptions.Buffer); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/AuthorizedQosFlowDescriptions): %w", err)
		}
	}
	if a.ExtendedProtocolConfigurationOptions != nil {
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetIei()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/ExtendedProtocolConfigurationOptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.GetLen()); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/ExtendedProtocolConfigurationOptions): %w", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.Buffer); err != nil {
			return fmt.Errorf("NAS encode error (PDUSessionModificationCommand/ExtendedProtocolConfigurationOptions): %w", err)
		}
	}
	return nil
}

func (a *PDUSessionModificationCommand) DecodePDUSessionModificationCommand(byteArray *[]byte) error {
	buffer := bytes.NewBuffer(*byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolDiscriminator.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/ExtendedProtocolDiscriminator): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PDUSessionID.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/PDUSessionID): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PTI.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/PTI): %w", err)
	}
	if err := binary.Read(buffer, binary.BigEndian, &a.PDUSESSIONMODIFICATIONCOMMANDMessageIdentity.Octet); err != nil {
		return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/PDUSESSIONMODIFICATIONCOMMANDMessageIdentity): %w", err)
	}
	for buffer.Len() > 0 {
		var ieiN uint8
		var tmpIeiN uint8
		if err := binary.Read(buffer, binary.BigEndian, &ieiN); err != nil {
			return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/iei): %w", err)
		}
		// fmt.Println(ieiN)
		if ieiN >= 0x80 {
			tmpIeiN = (ieiN & 0xf0) >> 4
		} else {
			tmpIeiN = ieiN
		}
		// fmt.Println("type", tmpIeiN)
		switch tmpIeiN {
		case PDUSessionModificationCommandCause5GSMType:
			a.Cause5GSM = nasType.NewCause5GSM(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.Cause5GSM.Octet); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/Cause5GSM): %w", err)
			}
		case PDUSessionModificationCommandSessionAMBRType:
			a.SessionAMBR = nasType.NewSessionAMBR(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.SessionAMBR.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/SessionAMBR): %w", err)
			}
			if a.SessionAMBR.Len != 6 {
				return fmt.Errorf("invalid ie length (PDUSessionModificationCommand/SessionAMBR): %d", a.SessionAMBR.Len)
			}
			a.SessionAMBR.SetLen(a.SessionAMBR.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.SessionAMBR.Octet[:]); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/SessionAMBR): %w", err)
			}
		case PDUSessionModificationCommandRQTimerValueType:
			a.RQTimerValue = nasType.NewRQTimerValue(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.RQTimerValue.Octet); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/RQTimerValue): %w", err)
			}
		case PDUSessionModificationCommandAlwaysonPDUSessionIndicationType:
			a.AlwaysonPDUSessionIndication = nasType.NewAlwaysonPDUSessionIndication(ieiN)
			a.AlwaysonPDUSessionIndication.Octet = ieiN
		case PDUSessionModificationCommandAuthorizedQosRulesType:
			a.AuthorizedQosRules = nasType.NewAuthorizedQosRules(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.AuthorizedQosRules.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/AuthorizedQosRules): %w", err)
			}
			if a.AuthorizedQosRules.Len < 4 {
				return fmt.Errorf("invalid ie length (PDUSessionModificationCommand/AuthorizedQosRules): %d", a.AuthorizedQosRules.Len)
			}
			a.AuthorizedQosRules.SetLen(a.AuthorizedQosRules.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.AuthorizedQosRules.Buffer); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/AuthorizedQosRules): %w", err)
			}
		case PDUSessionModificationCommandMappedEPSBearerContextsType:
			a.MappedEPSBearerContexts = nasType.NewMappedEPSBearerContexts(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.MappedEPSBearerContexts.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/MappedEPSBearerContexts): %w", err)
			}
			if a.MappedEPSBearerContexts.Len < 4 {
				return fmt.Errorf("invalid ie length (PDUSessionModificationCommand/MappedEPSBearerContexts): %d", a.MappedEPSBearerContexts.Len)
			}
			a.MappedEPSBearerContexts.SetLen(a.MappedEPSBearerContexts.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.MappedEPSBearerContexts.Buffer); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/MappedEPSBearerContexts): %w", err)
			}
		case PDUSessionModificationCommandAuthorizedQosFlowDescriptionsType:
			a.AuthorizedQosFlowDescriptions = nasType.NewAuthorizedQosFlowDescriptions(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.AuthorizedQosFlowDescriptions.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/AuthorizedQosFlowDescriptions): %w", err)
			}
			if a.AuthorizedQosFlowDescriptions.Len < 3 {
				return fmt.Errorf("invalid ie length (PDUSessionModificationCommand/AuthorizedQosFlowDescriptions): %d", a.AuthorizedQosFlowDescriptions.Len)
			}
			a.AuthorizedQosFlowDescriptions.SetLen(a.AuthorizedQosFlowDescriptions.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.AuthorizedQosFlowDescriptions.Buffer); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/AuthorizedQosFlowDescriptions): %w", err)
			}
		case PDUSessionModificationCommandExtendedProtocolConfigurationOptionsType:
			a.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(ieiN)
			if err := binary.Read(buffer, binary.BigEndian, &a.ExtendedProtocolConfigurationOptions.Len); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/ExtendedProtocolConfigurationOptions): %w", err)
			}
			if a.ExtendedProtocolConfigurationOptions.Len < 1 {
				return fmt.Errorf("invalid ie length (PDUSessionModificationCommand/ExtendedProtocolConfigurationOptions): %d", a.ExtendedProtocolConfigurationOptions.Len)
			}
			a.ExtendedProtocolConfigurationOptions.SetLen(a.ExtendedProtocolConfigurationOptions.GetLen())
			if err := binary.Read(buffer, binary.BigEndian, a.ExtendedProtocolConfigurationOptions.Buffer); err != nil {
				return fmt.Errorf("NAS decode error (PDUSessionModificationCommand/ExtendedProtocolConfigurationOptions): %w", err)
			}
		default:
		}
	}
	return nil
}
