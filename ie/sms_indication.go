package ie

// SMSInd is detailed in 9.11.3.50A SMS indication, 24.501
type SMSInd struct {
	// Name, uint8, Bits, Octet
	SAI bool // 1->1,   1->1
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *SMSInd) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "SMSInd",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *SMSInd) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "SMSInd",
	}
	return nil, e
}
