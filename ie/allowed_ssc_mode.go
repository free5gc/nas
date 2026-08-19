package ie

// AllowedSSCMode is detailed in 9.11.4.5 Allowed SSC mode, 24.501
type AllowedSSCMode struct {
	// Name, uint8, Bits, Octet
	SSC3 uint8 // 3->3,   1->1
	SSC2 uint8 // 2->2,   1->1
	SSC1 uint8 // 1->1,   1->1
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *AllowedSSCMode) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "AllowedSSCMode",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *AllowedSSCMode) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "AllowedSSCMode",
	}
	return nil, e
}
