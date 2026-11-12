package ie

import "github.com/pkg/errors"

// valuePart is the value part that a decoder reads of a type 4 or type 6 IE whose value part has at most n octets: a
// greater length is no syntactical error, and the extra octets are ignored (TS 24.007 11.4.2).
func valuePart(b []byte, n int) []byte {
	if len(b) > n {
		return b[:n]
	}
	return b
}

// copyIEValue validates the value length and gives the native IE its own buffer.
// Opaque values do not claim to implement the procedures carried inside them.
func copyIEValue(name string, b []byte, minLen, maxLen int) ([]byte, error) {
	if len(b) < minLen || len(b) > maxLen {
		return nil, errors.Errorf("%s IE length %d outside %d..%d", name, len(b), minLen, maxLen)
	}
	return append([]byte(nil), b...), nil
}
