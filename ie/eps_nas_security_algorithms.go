package ie

// EPSNASSecAlgos preserves the bounded IE value without interpreting its contents.
// Contents excludes the IEI and length octets.
type EPSNASSecAlgos struct {
	Contents []byte
}

func (i *EPSNASSecAlgos) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("EPSNASSecAlgos", b, 1, 1)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *EPSNASSecAlgos) MarshalBinary() ([]byte, error) {
	return copyIEValue("EPSNASSecAlgos", i.Contents, 1, 1)
}
