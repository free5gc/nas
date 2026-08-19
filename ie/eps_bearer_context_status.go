package ie

type EPSBearerCtxStatus struct{}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *EPSBearerCtxStatus) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "EPSBearerCtxStatus",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *EPSBearerCtxStatus) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "EPSBearerCtxStatus",
	}
	return nil, e
}
