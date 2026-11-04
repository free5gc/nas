package message

import (
	"bytes"
	"fmt"
	"io"
)

// UnknownRequiredIEError identifies an unknown comprehension-required IE.
// The caller can apply TS 24.501 section 7.5.1 rather than ignore the IE.
type UnknownRequiredIEError struct {
	IEI uint8
}

func (e *UnknownRequiredIEError) Error() string {
	return fmt.Sprintf("unknown comprehension-required IEI[0x%02x]", e.IEI)
}

// skipUnknownOptionalIE skips an unknown IE after its IEI octet has been read.
// These rules are for 5GMM/5GSM, per TS 24.007 V17.3.0 sections 11.2.4 and
// 11.2.5. In particular, this is not the EPS TLV-E classification.
func skipUnknownOptionalIE(reader *bytes.Buffer, iei uint8) error {
	if iei < 0x10 || (iei >= 0x7c && iei <= 0x7f) {
		return &UnknownRequiredIEError{IEI: iei}
	}
	if iei >= 0x80 {
		// A type 1/2 IE occupies only the already consumed octet.
		return nil
	}

	lengthType := IELen8Bits
	if iei >= 0x70 {
		lengthType = IELen16Bits
	}
	length, err := getIeLen(reader, lengthType)
	if err != nil {
		// No reliable boundary remains from which decoding could resume.
		return io.ErrUnexpectedEOF
	}
	reader.Next(int(length))
	return nil
}
