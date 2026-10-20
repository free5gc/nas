package ie

// OperatorDefinedAccessCategoryDefs is detailed in 9.11.3.38 Operator-defined access category definitions, 24.501
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type OperatorDefinedAccessCategoryDefs struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	OperatorDefinedAccessCategoryDef1 uint8
	OperatorDefinedAccessCategoryDef2 uint8
	OperatorDefinedAccessCategoryDefN uint8
	Contents                          []byte
}

func (i *OperatorDefinedAccessCategoryDefs) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("OperatorDefinedAccessCategoryDefs", b, 0, 8320)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *OperatorDefinedAccessCategoryDefs) MarshalBinary() ([]byte, error) {
	return copyIEValue("OperatorDefinedAccessCategoryDefs", i.Contents, 0, 8320)
}
