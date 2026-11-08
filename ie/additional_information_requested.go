package ie

import "github.com/pkg/errors"

// AdditionalInfoReq is detailed in 9.11.3.12A Additional information requested, 24.501
// Spare bits and extra octets (24.007 11.4.2) are ignored on decode; spare bits are emitted as zero.
type AdditionalInfoReq struct {
	Cipherkey uint8
}

func (i *AdditionalInfoReq) UnmarshalBinary(b []byte) error {
	if len(b) < 1 {
		return errors.Errorf("AdditionalInfoReq IE length %d, expected at least 1", len(b))
	}
	i.Cipherkey = GetBit1(b[0])
	return nil
}

func (i *AdditionalInfoReq) MarshalBinary() ([]byte, error) {
	if i.Cipherkey > 1 {
		return nil, errors.Errorf("invalid additional information requested CipherKey %d", i.Cipherkey)
	}
	return []byte{i.Cipherkey}, nil
}
