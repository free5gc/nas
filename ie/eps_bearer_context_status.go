package ie

// EPSBearerCtxStatus preserves the bounded IE value without interpreting its contents.
// Contents excludes the IEI and length octets.
type EPSBearerCtxStatus struct {
	Contents []byte
}

func (i *EPSBearerCtxStatus) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("EPSBearerCtxStatus", valuePart(b, 2), 2, 2)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *EPSBearerCtxStatus) MarshalBinary() ([]byte, error) {
	return copyIEValue("EPSBearerCtxStatus", i.Contents, 2, 2)
}
