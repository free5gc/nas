package ie

// OperatorDefinedAccessCategoryDefs is detailed in 9.11.3.38 Operator-defined access category definitions, 24.501
type OperatorDefinedAccessCategoryDefs struct {
	// Name, uint8, Bits, Octet
	OperatorDefinedAccessCategoryDef1 uint8 // 8->1,   4->995
	OperatorDefinedAccessCategoryDef2 uint8 // 8->1,   996->1
	OperatorDefinedAccessCategoryDefN uint8 // 8->1,   2->1
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *OperatorDefinedAccessCategoryDefs) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "OperatorDefinedAccessCategoryDefs",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *OperatorDefinedAccessCategoryDefs) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "OperatorDefinedAccessCategoryDefs",
	}
	return nil, e
}
