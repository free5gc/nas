package ie

// CongestionReattemptIndicator5GSM is detailed in 9.11.4.21 5GSM congestion re-attempt indicator, 24.501
type CongestionReattemptIndicator5GSM struct {
	// Name, uint8, Bits, Octet
	ABO uint8 // 1->1,   3->3
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *CongestionReattemptIndicator5GSM) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "CongestionReattemptIndicator5GSM",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *CongestionReattemptIndicator5GSM) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "CongestionReattemptIndicator5GSM",
	}
	return nil, e
}
