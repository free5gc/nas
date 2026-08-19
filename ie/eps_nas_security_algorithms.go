package ie

type EPSNASSecAlgos struct{}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *EPSNASSecAlgos) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "EPSNASSecAlgos",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *EPSNASSecAlgos) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "EPSNASSecAlgos",
	}
	return nil, e
}
