package ie

import "github.com/pkg/errors"

// SMSInd is detailed in 9.11.3.50A SMS indication, 24.501
// Spare bits are ignored on decode and emitted as zero.
type SMSInd struct {
	SAI bool
}

func (i *SMSInd) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return errors.Errorf("SMSInd IE length %d, expected 1", len(b))
	}
	i.SAI = GetBit1(b[0]) == 1
	return nil
}

func (i *SMSInd) MarshalBinary() ([]byte, error) {
	return []byte{bool2uint8(i.SAI)}, nil
}
