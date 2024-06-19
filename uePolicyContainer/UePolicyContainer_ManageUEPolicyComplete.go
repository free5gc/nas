package uePolicyContainer

import (
	"bytes"
	"encoding/binary"

	"github.com/free5gc/nas/nasType"
)

type ManageUEPolicyComplete struct {
	nasType.PTI
	UePolicyDeliveryServiceMsgType
}

func NewManageUEPolicyComplete(msgType uint8) (manageUEPolicyComplete *ManageUEPolicyComplete) {
	manageUEPolicyComplete = &ManageUEPolicyComplete{}
	manageUEPolicyComplete.UePolicyDeliveryServiceMsgType.SetMessageIdentity(msgType)
	return manageUEPolicyComplete
}

func (m *ManageUEPolicyComplete) EncodeManageUEPolicyComplete(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, m.PTI.Octet); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, m.UePolicyDeliveryServiceMsgType.Octet); err != nil {
		return err
	}
	return nil
}

func (m *ManageUEPolicyComplete) DecodeManageUEPolicyComplete(byteArray []byte) error {
	buffer := bytes.NewBuffer(byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &m.PTI.Octet); err != nil {
		return err
	}
	if err := binary.Read(buffer, binary.BigEndian, &m.UePolicyDeliveryServiceMsgType.Octet); err != nil {
		return err
	}
	return nil
}
