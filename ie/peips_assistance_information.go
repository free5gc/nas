package ie

// PEIPSAssistanceInfo is detailed in 9.11.3.80 PEIPS assistance information, 24.501
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type PEIPSAssistanceInfo struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	PEIPSAssistanceInfoType1 uint8
	PEIPSAssistanceInfoType2 uint8
	PEIPSAssistanceInfoTypeP uint8
	Contents                 []byte
}

func (i *PEIPSAssistanceInfo) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("PEIPSAssistanceInfo", b, 1, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *PEIPSAssistanceInfo) MarshalBinary() ([]byte, error) {
	return copyIEValue("PEIPSAssistanceInfo", i.Contents, 1, 255)
}
