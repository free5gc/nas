package ie

import "github.com/pkg/errors"

// AlwaysonPDUSessInd is detailed in 9.11.4.3 Always-on PDU session indication, 24.501
// Spare bits are ignored on decode and emitted as zero.
type AlwaysonPDUSessInd struct {
	APSI bool
}

func (i *AlwaysonPDUSessInd) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return errors.Errorf("AlwaysonPDUSessInd IE length %d, expected 1", len(b))
	}
	i.APSI = GetBit1(b[0]) == 1
	return nil
}

func (i *AlwaysonPDUSessInd) MarshalBinary() ([]byte, error) {
	return []byte{bool2uint8(i.APSI)}, nil
}
