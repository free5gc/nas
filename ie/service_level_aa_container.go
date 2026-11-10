package ie

// SvcLvlAACntr is detailed in 9.11.2.10 Service-level-AA container, 24.501
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type SvcLvlAACntr struct {
	Contents []byte
}

func (i *SvcLvlAACntr) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("SvcLvlAACntr", b, 1, 65535)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *SvcLvlAACntr) MarshalBinary() ([]byte, error) {
	return copyIEValue("SvcLvlAACntr", i.Contents, 1, 65535)
}
