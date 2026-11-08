package ie

import "github.com/pkg/errors"

// MAPDUSessInfo is detailed in 9.11.3.31A MA PDU session information, 24.501
type MAPDUSessInfo struct {
	MAPDUSessInfoValue uint8
}

func (i *MAPDUSessInfo) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return errors.Errorf("MAPDUSessInfo IE length %d, expected 1", len(b))
	}
	i.MAPDUSessInfoValue = Get4Bits41(b[0])
	return nil
}

func (i *MAPDUSessInfo) MarshalBinary() ([]byte, error) {
	if i.MAPDUSessInfoValue > 15 {
		return nil, errors.Errorf("invalid MA PDU session information value %d", i.MAPDUSessInfoValue)
	}
	return []byte{i.MAPDUSessInfoValue}, nil
}
