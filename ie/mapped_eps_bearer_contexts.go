package ie

type MappedEPSBearerCtxs struct{}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *MappedEPSBearerCtxs) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "MappedEPSBearerCtxs",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *MappedEPSBearerCtxs) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "MappedEPSBearerCtxs",
	}
	return nil, e
}
