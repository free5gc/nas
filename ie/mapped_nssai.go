package ie

// MappedNSSAI is detailed in 9.11.3.31B Mapped NSSAI, 24.501
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type MappedNSSAI struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	MappedSNSSAIContent1 uint8
	MappedSNSSAIContent2 uint8
	MappedSNSSAIContentN uint8
	Contents             []byte
}

func (i *MappedNSSAI) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("MappedNSSAI", b, 2, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *MappedNSSAI) MarshalBinary() ([]byte, error) {
	return copyIEValue("MappedNSSAI", i.Contents, 2, 255)
}
