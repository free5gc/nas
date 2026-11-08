package ie

import "github.com/pkg/errors"

// N5GCInd is detailed in 9.11.3.72 N5GC indication, 24.501
// Spare bits are ignored on decode and emitted as zero.
type N5GCInd struct {
	N5GCREG bool
}

func (i *N5GCInd) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return errors.Errorf("N5GCInd IE length %d, expected 1", len(b))
	}
	i.N5GCREG = GetBit1(b[0]) == 1
	return nil
}

func (i *N5GCInd) MarshalBinary() ([]byte, error) {
	return []byte{SetBit1(0, bool2uint8(i.N5GCREG))}, nil
}
