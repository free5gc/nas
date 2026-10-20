package ie

import "github.com/pkg/errors"

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
	if len(b) < 1 || len(b) > 65535 {
		return errors.Errorf("PayloadCntr IE length %d outside 1..65535", len(b))
	}
	i.Pct = pct
	i.Contents = append([]byte(nil), b...)
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
	if ieLen < 1 || ieLen > 65535 {
		return nil, errors.Errorf("PayloadCntr IE length %d outside 1..65535", ieLen)
	}
	b := make([]byte, ieLen)
	copy(b, i.Contents)
	return b, nil
}
