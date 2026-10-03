package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIELengthRequiresCompleteValue(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		size IELengthType
		wire []byte
	}{
		{"LV/TLV", IELen8Bits, []byte{2, 0x11}},
		{"LVE/TLVE", IELen16Bits, []byte{0, 2, 0x11}},
		{"empty value", IELen8Bits, []byte{1}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getIeLen(bytes.NewBuffer(tc.wire), tc.size)
			require.Error(t, err)
		})
	}
	for _, size := range []IELengthType{IELen8Bits, IELen16Bits} {
		wire := []byte{1, 0x11}
		if size == IELen16Bits {
			wire = append([]byte{0}, wire...)
		}
		reader := bytes.NewBuffer(wire)
		length, err := getIeLen(reader, size)
		require.NoError(t, err)
		require.Equal(t, uint16(1), length)
		require.Equal(t, []byte{0x11}, reader.Bytes())
	}
}

func TestDLNASTransportRejectsTruncatedDeclaredPayload(t *testing.T) {
	t.Parallel()
	var message DLNASTransport
	// The mandatory payload declares two value octets but contains only one.
	err := message.UnmarshalBinary([]byte{0x7e, 0x00, 0x68, 0x01, 0x00, 0x02, 0x2e})
	require.ErrorContains(t, err, "IE length 2 exceeds remaining data 1")
}

func TestDLNASTransportRejectsTruncatedRepeatedIE(t *testing.T) {
	t.Parallel()
	var message DLNASTransport
	// A complete payload and timer are followed by a truncated second timer.
	err := message.UnmarshalBinary([]byte{
		0x7e, 0x00, 0x68, 0x01, 0x00, 0x01, 0x2e,
		DLNASTransportIEIBackoffTimerValue, 0x01, 0x01,
		DLNASTransportIEIBackoffTimerValue, 0x02, 0x01,
	})
	require.ErrorContains(t, err, "IE length 2 exceeds remaining data 1")
}
