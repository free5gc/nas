package ie

// AdditionalInfo is detailed in 9.11.2.1 Additional information, 24.501
type AdditionalInfo struct {
	// Name, uint8, Bits, Octet
	AdditionalInfoValue uint8 // 8->1,   3->3
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *AdditionalInfo) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "AdditionalInfo",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *AdditionalInfo) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "AdditionalInfo",
	}
	return nil, e
}
