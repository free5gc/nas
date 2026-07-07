package ie

// PayloadCntr is detailed in 9.11.3.39 Payload container, 24.501
type PayloadCntr struct {
	// Name, uint8, Bits, Octet
	Pct      uint8 // ConstPayloadCntrType
	Contents []byte
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *PayloadCntr) UnmarshalBinary(b []byte, pct uint8) error {
	switch pct {
	case PayloadCntrType_MultiplePayloads, PayloadCntrType_EventNotif:
		var e error = &IEToDo{
			IEName: "(Multiple)PayloadCntr",
		}
		return e
	}
	i.Pct = pct
	i.Contents = append(i.Contents, b...)
	return nil
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *PayloadCntr) MarshalBinary() ([]byte, error) {
	switch i.Pct {
	case PayloadCntrType_MultiplePayloads, PayloadCntrType_EventNotif:
		var e error = &IEToDo{
			IEName: "(Multiple)PayloadCntr",
		}
		return nil, e
	}
	ieLen := len(i.Contents)
	b := make([]byte, ieLen)
	copy(b, i.Contents)
	return b, nil
}
