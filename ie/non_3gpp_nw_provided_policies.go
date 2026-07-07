package ie

type Non3GppNWProvidedPolicies struct{}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *Non3GppNWProvidedPolicies) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "Non3GppNWProvidedPolicies",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *Non3GppNWProvidedPolicies) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "Non3GppNWProvidedPolicies",
	}
	return nil, e
}
