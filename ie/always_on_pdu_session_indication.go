package ie

// AlwaysonPDUSessInd is detailed in 9.11.4.3 Always-on PDU session indication, 24.501
type AlwaysonPDUSessInd struct {
	// Name, uint8, Bits, Octet
	APSI bool // 1->1,   1->1
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *AlwaysonPDUSessInd) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "AlwaysonPDUSessInd",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *AlwaysonPDUSessInd) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "AlwaysonPDUSessInd",
	}
	return nil, e
}
