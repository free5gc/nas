package ie

// MICOInd is detailed in 9.11.3.31 MICO indication, 24.501
type MICOInd struct {
	// Name, uint8, Bits, Octet
	SPRTI bool // 2->2,   1->1
	RAAI  bool // 1->1,   1->1
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *MICOInd) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "MICOInd",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *MICOInd) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "MICOInd",
	}
	return nil, e
}
