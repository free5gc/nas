package ie

type S1UESecCapability struct{}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *S1UESecCapability) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "S1UESecCapability",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *S1UESecCapability) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "S1UESecCapability",
	}
	return nil, e
}
