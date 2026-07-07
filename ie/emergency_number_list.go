package ie

type EmergNumList struct{}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *EmergNumList) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "EmergNumList",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *EmergNumList) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "EmergNumList",
	}
	return nil, e
}
