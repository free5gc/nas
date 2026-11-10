package ie

// WUSAssistanceInfo is detailed in 9.11.3.71 WUS assistance information, 24.501 (9.9.3.62, 24.301)
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type WUSAssistanceInfo struct {
	Contents []byte
}

func (i *WUSAssistanceInfo) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("WUSAssistanceInfo", b, 1, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *WUSAssistanceInfo) MarshalBinary() ([]byte, error) {
	return copyIEValue("WUSAssistanceInfo", i.Contents, 1, 255)
}
