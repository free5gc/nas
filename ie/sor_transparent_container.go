package ie

// SORTransparentCntr is detailed in 9.11.3.51 SOR transparent container, 24.501
// The bounded value is preserved without interpreting the inner contents.
// Contents excludes the IEI and length octets.
type SORTransparentCntr struct {
	// Deprecated: generated placeholder fields were never encoded or decoded.
	// Use Contents for the complete IE value.
	SORHdr                        uint8
	SORMACIAUSF                   uint8
	Countersor                    uint8
	SecuredPkt                    uint8
	PLMNIDAndAccessTechnologyList uint8
	Contents                      []byte
}

func (i *SORTransparentCntr) UnmarshalBinary(b []byte) error {
	value, err := copyIEValue("SORTransparentCntr", b, 17, 65535)
	if err != nil {
		return err
	}
	i.Contents = value
	return nil
}

func (i *SORTransparentCntr) MarshalBinary() ([]byte, error) {
	return copyIEValue("SORTransparentCntr", i.Contents, 17, 65535)
}
