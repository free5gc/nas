package ie

import "github.com/pkg/errors"

// NBN1ModeDRXParams is detailed in 9.11.3.73 NB-N1 mode DRX parameters, 24.501
// Spare bits and extra octets (24.007 11.4.2) are ignored on decode; spare bits are emitted as zero.
type NBN1ModeDRXParams struct {
	NBN1ModeDRXValue uint8
}

func (i *NBN1ModeDRXParams) UnmarshalBinary(b []byte) error {
	if len(b) < 1 {
		return errors.Errorf("NBN1ModeDRXParams IE length %d, expected at least 1", len(b))
	}
	// Values 6 and 8 to 15 "shall be interpreted as DRX value not specified".
	if i.NBN1ModeDRXValue = Get4Bits41(b[0]); i.NBN1ModeDRXValue == 6 || i.NBN1ModeDRXValue > 7 {
		i.NBN1ModeDRXValue = 0
	}
	return nil
}

func (i *NBN1ModeDRXParams) MarshalBinary() ([]byte, error) {
	if i.NBN1ModeDRXValue > 15 {
		return nil, errors.Errorf("invalid NB-N1 mode DRX value %d", i.NBN1ModeDRXValue)
	}
	return []byte{i.NBN1ModeDRXValue}, nil
}
