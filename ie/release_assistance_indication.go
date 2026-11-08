package ie

import "github.com/pkg/errors"

// RelAssistanceInd is detailed in 9.11.3.46A Release assistance indication, 24.501 (9.9.4.25, 24.301)
// Spare bits are ignored on decode and emitted as zero.
type RelAssistanceInd struct {
	DDX uint8
}

func (i *RelAssistanceInd) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return errors.Errorf("RelAssistanceInd IE length %d, expected 1", len(b))
	}
	i.DDX = b[0] & 3
	return nil
}

func (i *RelAssistanceInd) MarshalBinary() ([]byte, error) {
	if i.DDX > 3 {
		return nil, errors.Errorf("invalid release assistance indication DDX %d", i.DDX)
	}
	return []byte{i.DDX}, nil
}
