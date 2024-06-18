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

func (m *ManageUEPolicyComplete) EncodeManageUEPolicyComplete(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, m.PTI.Octet)
	binary.Write(buffer, binary.BigEndian, m.UePolicyDeliveryServiceMsgType.Octet)
}

func (m *ManageUEPolicyComplete) DecodeManageUEPolicyComplete(byteArray []byte) error {
	buffer := bytes.NewBuffer(byteArray)
	binary.Read(buffer, binary.BigEndian, &m.PTI.Octet)
	binary.Read(buffer, binary.BigEndian, &m.UePolicyDeliveryServiceMsgType.Octet)
	return nil
}
