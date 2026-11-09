package ie

// EPSNASMsgCntr is detailed in 9.11.3.24 EPS NAS message container, 24.501
// The bounded value is preserved without interpreting the inner contents:
// TS 24.501 7.5.2 allows no diagnosis other than presence and length.
// Contents excludes the IEI and length octets.
type EPSNASMsgCntr struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	EPSNASMsgCntr uint8
	Contents      []byte
}

func (i *EPSNASMsgCntr) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("EPSNASMsgCntr", b, 1, 65535)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *EPSNASMsgCntr) MarshalBinary() ([]byte, error) {
	return copyIEValue("EPSNASMsgCntr", i.Contents, 1, 65535)
}
