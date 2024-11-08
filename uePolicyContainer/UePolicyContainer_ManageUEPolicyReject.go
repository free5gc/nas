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

func (m *ManageUEPolicyReject) EncodeManageUEPolicyReject(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, m.PTI.Octet); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, m.UePolicyDeliveryServiceMsgType.Octet); err != nil {
		return err
	}
	// UEPolicySectionManagementResult
	uePolSecMngRsult, err := m.UEPolicySectionManagementResult.MarshalBinary()
	if err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, uePolSecMngRsult); err != nil {
		return err
	}
	return nil
}

func (m *ManageUEPolicyReject) DecodeManageUEPolicyReject(byteArray []byte) error {
	buffer := bytes.NewBuffer(byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &m.PTI.Octet); err != nil {
		return err
	}
	if err := binary.Read(buffer, binary.BigEndian, &m.UePolicyDeliveryServiceMsgType.Octet); err != nil {
		return err
	}
	if err := m.UEPolicySectionManagementResult.UnmarshalBinary(buffer); err != nil {
		return err
	}
	return nil
}
