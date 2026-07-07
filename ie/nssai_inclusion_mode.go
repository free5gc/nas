package ie

// NSSAIInclusionMode is detailed in 9.11.3.37A NSSAI inclusion mode, 24.501
type NSSAIInclusionMode struct {
	// Name, uint8, Bits, Octet
	NSSAIInclusionMode uint8 // 2->1,   1->1
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *NSSAIInclusionMode) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "NSSAIInclusionMode",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *NSSAIInclusionMode) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "NSSAIInclusionMode",
	}
	return nil, e
}
