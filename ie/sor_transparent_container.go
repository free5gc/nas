package ie

// SORTransparentCntr is detailed in 9.11.3.51 SOR transparent container, 24.501
type SORTransparentCntr struct {
	// Name, uint8, Bits, Octet
	SORHdr                        uint8 // 8->1,   4->4
	SORMACIAUSF                   uint8 // 8->1,   5->520
	Countersor                    uint8 // 8->1,   521->2122
	SecuredPkt                    uint8 // 8->1,   2123->23
	PLMNIDAndAccessTechnologyList uint8 // 8->1,   2123->23102
}

// UnmarshalBinary handles the value part, not including IEI and length.
func (i *SORTransparentCntr) UnmarshalBinary(b []byte) error {
	var e error = &IEToDo{
		IEName: "SORTransparentCntr",
	}
	return e
}

// MarshalBinary returns the value part, not including IEI and length.
func (i *SORTransparentCntr) MarshalBinary() ([]byte, error) {
	var e error = &IEToDo{
		IEName: "SORTransparentCntr",
	}
	return nil, e
}
