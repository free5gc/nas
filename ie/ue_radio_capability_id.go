package ie

// UERadioCapabilityID is detailed in 9.11.3.68 UE radio capability ID, 24.501
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type UERadioCapabilityID struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	UERadioCapabilityID uint8
	Contents            []byte
}

func (i *UERadioCapabilityID) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("UERadioCapabilityID", b, 1, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *UERadioCapabilityID) MarshalBinary() ([]byte, error) {
	return copyIEValue("UERadioCapabilityID", i.Contents, 1, 255)
}
