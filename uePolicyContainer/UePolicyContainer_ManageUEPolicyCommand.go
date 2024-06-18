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

func (m *ManageUEPolicyCommand) EncodeManageUEPolicyCommand(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, m.PTI.Octet)
	binary.Write(buffer, binary.BigEndian, m.UePolicyDeliveryServiceMsgType.Octet)
	// structure: UEPolicySectionManagementList
	binary.Write(buffer, binary.BigEndian, m.UEPolicySectionManagementList.GetIei())
	binary.Write(buffer, binary.BigEndian, m.UEPolicySectionManagementList.GetLen())
	binary.Write(buffer, binary.BigEndian, m.UEPolicySectionManagementList.GetUEPolicySectionManagementListContent())
	// Optinal structure: UEPolicyNetworkClassmark
	if m.UEPolicyNetworkClassmark != nil {
		binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetIei())
		binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetLen())
		binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetNSSUI())
		binary.Write(buffer, binary.BigEndian, m.UEPolicyNetworkClassmark.GetSpare())
	}
}

func (m *ManageUEPolicyCommand) DecodeManageUEPolicyCommand(byteArray []byte) error {
	buffer := bytes.NewBuffer(byteArray)
	binary.Read(buffer, binary.BigEndian, &m.PTI.Octet)
	binary.Read(buffer, binary.BigEndian, &m.UePolicyDeliveryServiceMsgType.Octet)
	// structure: UEPolicySectionManagementList
	binary.Read(buffer, binary.BigEndian, &m.UEPolicySectionManagementList.Iei)
	binary.Read(buffer, binary.BigEndian, &m.UEPolicySectionManagementList.Len)
	m.UEPolicySectionManagementList.Buffer = make([]uint8, m.UEPolicySectionManagementList.Len)
	binary.Read(buffer, binary.BigEndian, m.UEPolicySectionManagementList.Buffer[:m.UEPolicySectionManagementList.GetLen()])

	// optinal structure: UEPolicyNetworkClassmark
	if buffer.Len() > 0 {
		// initial pointer type element
		m.UEPolicyNetworkClassmark = NewUEPolicyNetworkClassmark()
		binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.Iei)
		binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.Len)
		binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.NSSUI)
		binary.Read(buffer, binary.BigEndian, &m.UEPolicyNetworkClassmark.Spare)
		if buffer.Len() > 0 {
			return errors.New("deecode [Manage UE Policy Command] Error: nas msg out of range")
		}
	}
	return nil
}
