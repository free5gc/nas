package uePolicyContainer

import (
	"bytes"
	"encoding/binary"

	"github.com/free5gc/nas/nasType"
)

type ManageUEPolicyReject struct {
	nasType.PTI
	UePolicyDeliveryServiceMsgType
	UEPolicySectionManagementResult
}

func NewManageUEPolicyReject(msgType uint8) (manageUEPolicyReject *ManageUEPolicyReject) {
	manageUEPolicyReject = &ManageUEPolicyReject{}
	manageUEPolicyReject.UePolicyDeliveryServiceMsgType.SetMessageIdentity(msgType)
	return manageUEPolicyReject
}

func (m *ManageUEPolicyReject) EncodeManageUEPolicyReject(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, m.PTI.Octet)
	binary.Write(buffer, binary.BigEndian, m.UePolicyDeliveryServiceMsgType.Octet)
	// UEPolicySectionManagementResult
	uePolSecMngRsult, _ := m.UEPolicySectionManagementResult.MarshalBinary()
	binary.Write(buffer, binary.BigEndian, uePolSecMngRsult)
}

func (m *ManageUEPolicyReject) DecodeManageUEPolicyReject(byteArray []byte) error {
	buffer := bytes.NewBuffer(byteArray)
	binary.Read(buffer, binary.BigEndian, &m.PTI.Octet)
	binary.Read(buffer, binary.BigEndian, &m.UePolicyDeliveryServiceMsgType.Octet)
	if err := m.UEPolicySectionManagementResult.UnmarshalBinary(buffer); err != nil {
		return err
	}
	return nil
}
