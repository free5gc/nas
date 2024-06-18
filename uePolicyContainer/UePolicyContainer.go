package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// UePolicyContainer is a kind of NAS_DLTransport message
// UePolicyContainer maps to UE_policy_delivery_service
// TS 124 501 V17.7.1, page.760, Table 9.11.3.39.1: Payload container information element
type UePolicyContainer = UePolDeliverySer

// TS 124 501 V17.7.1, page.893, UE policy delivery service
type UePolDeliverySer struct {
	UePolDeliveryHeader
	*ManageUEPolicyCommand  // D.2.1.2 Network-requested UE policy management procedure initiation
	*ManageUEPolicyComplete // D.2.1.3 Network-requested UE policy management procedure accepted by the UE
	*ManageUEPolicyReject   // D.2.1.4 Network-requested UE policy management procedure not accepted by the UE
	// TODO: add below encoding way
	// *nasMessage.UEStateIndication      // D.2.2 UE-initiated UE state indication procedure
}

func NewUePolDeliverySer() *UePolDeliverySer {
	uePolDeliverySer := &UePolDeliverySer{}
	return uePolDeliverySer
}

type UePolDeliveryHeader struct {
	Octet [2]uint8
}

// PTI(Procedure transaction identity) is assign by PCF
func (u *UePolDeliveryHeader) SetHeaderPTI(PTI uint8) {
	u.Octet[0] = PTI
}

// PTI(Procedure transaction identity)
func (u *UePolDeliveryHeader) GetHeaderPTI() (pTI uint8) {
	pTI = u.Octet[0]
	return pTI
}

func (u *UePolDeliveryHeader) GetHeaderMessageType() (messageType uint8) {
	messageType = u.Octet[1]
	return messageType
}

func (u *UePolDeliveryHeader) SetHeaderMessageType(messageType uint8) {
	u.Octet[1] = messageType
}

const (
	MsgTypeManageUEPolicyCommand       uint8 = 1
	MsgTypeManageUEPolicyComplete      uint8 = 2
	MsgTypeManageUEPolicyReject        uint8 = 3
	MsgTypeUEStateIndication           uint8 = 4
	MsgTypeUEPolicyProvisioningRequest uint8 = 5
	MsgTypeUEPolicyProvisioningReject  uint8 = 6
)

// TODO: add other Decoding processes
func (u *UePolDeliverySer) UePolDeliverySerDecode(byteArray []byte) error {
	buffer := bytes.NewBuffer(byteArray)
	if err := binary.Read(buffer, binary.BigEndian, &u.UePolDeliveryHeader); err != nil {
		return fmt.Errorf("uePolicyContainer-UePolDeliveryHeader decode Fail:%+v", err)
	}

	switch u.GetHeaderMessageType() {
	case MsgTypeManageUEPolicyCommand:
		u.ManageUEPolicyCommand = NewManageUEPolicyCommand(MsgTypeManageUEPolicyCommand)
		if err := u.DecodeManageUEPolicyCommand(byteArray); err != nil {
			return err
		}
	case MsgTypeManageUEPolicyComplete:
		u.ManageUEPolicyComplete = NewManageUEPolicyComplete(MsgTypeManageUEPolicyComplete)
		if err := u.DecodeManageUEPolicyComplete(byteArray); err != nil {
			return err
		}
	case MsgTypeManageUEPolicyReject:
		u.ManageUEPolicyReject = NewManageUEPolicyReject(MsgTypeManageUEPolicyReject)
		if err := u.DecodeManageUEPolicyReject(byteArray); err != nil {
			return err
		}
	case MsgTypeUEStateIndication:
		// TODO
	case MsgTypeUEPolicyProvisioningRequest:
		// TODO
	case MsgTypeUEPolicyProvisioningReject:
		// TODO
	default:
		return fmt.Errorf("ue Policy Delivery Service decode Fail: MsgType[%d] doesn't exist",
			u.GetHeaderMessageType())
	}
	return nil
}

// TODO: add other Encoding processes
func (u *UePolDeliverySer) UePolDeliverySerEncode() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	switch u.GetHeaderMessageType() {
	case MsgTypeManageUEPolicyCommand:
		u.EncodeManageUEPolicyCommand(buf)
	case MsgTypeManageUEPolicyComplete:
		u.EncodeManageUEPolicyComplete(buf)
	case MsgTypeManageUEPolicyReject:
		u.EncodeManageUEPolicyReject(buf)
	case MsgTypeUEStateIndication:
		// TODO
	case MsgTypeUEPolicyProvisioningRequest:
		// TODO
	case MsgTypeUEPolicyProvisioningReject:
		// TODO
	default:
		return nil, fmt.Errorf("ue Policy Delivery Service Encode Fail: MsgType[%d] doesn't exist",
			u.GetHeaderMessageType())
	}
	return buf.Bytes(), nil
}
