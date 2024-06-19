package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/free5gc/nas/nasType"
)

// TS24.501 v17.7.1, sec D.5.1
type ManageUEPolicyCommand struct {
	nasType.PTI
	UePolicyDeliveryServiceMsgType
	UEPolicySectionManagementList
	*UEPolicyNetworkClassmark
}

func NewManageUEPolicyCommand(msgType uint8) (manageUEPolicyCommand *ManageUEPolicyCommand) {
	manageUEPolicyCommand = &ManageUEPolicyCommand{}
	manageUEPolicyCommand.UePolicyDeliveryServiceMsgType.SetMessageIdentity(msgType)
	return manageUEPolicyCommand
}

func (m *ManageUEPolicyCommand) EncodeManageUEPolicyCommand(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, m.PTI.Octet); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, m.UePolicyDeliveryServiceMsgType.Octet); err != nil {
		return err
	}
	// structure: UEPolicySectionManagementList
	{
		if err := binary.Write(buffer, binary.BigEndian, m.UEPolicySectionManagementList.GetIei()); err != nil {
			return err
		}
		if err := binary.Write(buffer, binary.BigEndian, m.UEPolicySectionManagementList.GetLen()); err != nil {
			return err
		}
		if err := binary.Write(buffer, binary.BigEndian,
			m.UEPolicySectionManagementList.GetUEPolicySectionManagementListContent()); err != nil {
			return err
		}
	}
	// Optinal structure: UEPolicyNetworkClassmark
	if m.UEPolicyNetworkClassmark != nil {
		if err := binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetIei()); err != nil {
			return err
		}
		if err := binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetLen()); err != nil {
			return err
		}
		if err := binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetNSSUI()); err != nil {
			return err
		}
		if err := binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetSpare()); err != nil {
			return err
		}
	}
	return nil
}

func (m *ManageUEPolicyCommand) DecodeManageUEPolicyCommand(byteArray []byte) error {
	buffer := bytes.NewBuffer(byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &m.PTI.Octet); err != nil {
		return err
	}
	if err := binary.Read(buffer, binary.BigEndian, &m.UePolicyDeliveryServiceMsgType.Octet); err != nil {
		return err
	}

	{ // structure: UEPolicySectionManagementList
		if err := binary.Read(buffer, binary.BigEndian, &m.UEPolicySectionManagementList.Iei); err != nil {
			return err
		}
		if err := binary.Read(buffer, binary.BigEndian, &m.UEPolicySectionManagementList.Len); err != nil {
			return err
		}
		m.UEPolicySectionManagementList.Buffer = make([]uint8, m.UEPolicySectionManagementList.Len)
		if err := binary.Read(buffer, binary.BigEndian,
			m.UEPolicySectionManagementList.Buffer[:m.UEPolicySectionManagementList.GetLen()]); err != nil {
			return err
		}
	}
	// optinal structure: UEPolicyNetworkClassmark
	if buffer.Len() > 0 {
		// initial pointer type element
		m.UEPolicyNetworkClassmark = NewUEPolicyNetworkClassmark()
		if err := binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.Iei); err != nil {
			return err
		}
		if err := binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.Len); err != nil {
			return err
		}
		if err := binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.NSSUI); err != nil {
			return err
		}
		if err := binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.Spare); err != nil {
			return err
		}
		if buffer.Len() > 0 {
			return errors.New("deecode [Manage UE Policy Command] Error: nas msg out of range")
		}
	}
	return nil
}
