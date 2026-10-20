package ie

import "github.com/pkg/errors"

// NSSAIInclusionMode is detailed in 9.11.3.37A NSSAI inclusion mode, 24.501
// Spare bits are ignored on decode and emitted as zero.
type NSSAIInclusionMode struct {
	NSSAIInclusionMode uint8
}

func (i *NSSAIInclusionMode) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return errors.Errorf("NSSAIInclusionMode IE length %d, expected 1", len(b))
	}
	i.NSSAIInclusionMode = b[0] & 3
	return nil
}

func (i *NSSAIInclusionMode) MarshalBinary() ([]byte, error) {
	if i.NSSAIInclusionMode > 3 {
		return nil, errors.Errorf("invalid NSSAI inclusion mode %d", i.NSSAIInclusionMode)
	}
	return []byte{i.NSSAIInclusionMode}, nil
}
