package ie

// SupportedCodecList is detailed in 9.11.3.51A Supported codec list, 24.501 (10.5.4.32, 24.008)
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type SupportedCodecList struct {
	Contents []byte
}

func (i *SupportedCodecList) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("SupportedCodecList", b, 3, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *SupportedCodecList) MarshalBinary() ([]byte, error) {
	return copyIEValue("SupportedCodecList", i.Contents, 3, 255)
}
