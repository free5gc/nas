package ie

// MobileStationClassmark2 is detailed in 9.11.3.31C Mobile station classmark 2, 24.501 (10.5.1.6, 24.008)
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type MobileStationClassmark2 struct {
	Contents []byte
}

func (i *MobileStationClassmark2) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("MobileStationClassmark2", b, 3, 255)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *MobileStationClassmark2) MarshalBinary() ([]byte, error) {
	return copyIEValue("MobileStationClassmark2", i.Contents, 3, 255)
}
