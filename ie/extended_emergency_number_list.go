package ie

type ExtendedEmergNumList struct{}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *ExtendedEmergNumList) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "ExtendedEmergNumList",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *ExtendedEmergNumList) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "ExtendedEmergNumList",
	}
	return nil, e
}
