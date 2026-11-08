package ie

// NID is detailed in 9.11.3.79 NID, 24.501 (9.2.7, 24.502)
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type NID struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	NIDValueDigit1  uint8
	AssignmentMode  uint8
	NIDValueDigit3  uint8
	NIDValueDigit2  uint8
	NIDValueDigit5  uint8
	NIDValueDigit4  uint8
	NIDValueDigit7  uint8
	NIDValueDigit6  uint8
	NIDValueDigit9  uint8
	NIDValueDigit8  uint8
	NIDValueDigit10 uint8
	Contents        []byte
}

func (i *NID) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("NID", b, 6, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *NID) MarshalBinary() ([]byte, error) {
	return copyIEValue("NID", i.Contents, 6, 255)
}
