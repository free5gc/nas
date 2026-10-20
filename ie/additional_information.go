package ie

// AdditionalInfo is detailed in 9.11.2.1 Additional information, 24.501
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type AdditionalInfo struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	AdditionalInfoValue uint8
	Contents            []byte
}

func (i *AdditionalInfo) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("AdditionalInfo", b, 1, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *AdditionalInfo) MarshalBinary() ([]byte, error) {
	return copyIEValue("AdditionalInfo", i.Contents, 1, 255)
}
