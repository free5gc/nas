package ie

import "github.com/pkg/errors"

// copyIEValue validates the value length and gives the native IE its own buffer.
// Opaque values do not claim to implement the procedures carried inside them.
func copyIEValue(name string, b []byte, minLen, maxLen int) ([]byte, error) {
	if len(b) < minLen || len(b) > maxLen {
		return nil, errors.Errorf("%s IE length %d outside %d..%d", name, len(b), minLen, maxLen)
	}
	return append([]byte(nil), b...), nil
}
