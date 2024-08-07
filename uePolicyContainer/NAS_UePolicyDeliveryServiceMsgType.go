package uePolicyContainer

// TS24.501 v17.7.1, D.6.1
type UePolicyDeliveryServiceMsgType struct {
	Octet uint8
}

func (m *UePolicyDeliveryServiceMsgType) SetMessageIdentity(MessageIdentity uint8) {
	m.Octet = MessageIdentity
}

func (m *UePolicyDeliveryServiceMsgType) GetMessageIdentity(MessageIdentity uint8) uint8 {
	return m.Octet
}
