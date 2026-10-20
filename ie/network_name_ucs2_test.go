package ie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetworkNameUCS2(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		text string
		wire []byte
	}{
		{"", []byte{0x90}},
		{"FRé", []byte{0x90, 0x00, 0x46, 0x00, 0x52, 0x00, 0xe9}},
		{"網路", []byte{0x90, 0x7d, 0xb2, 0x8d, 0xef}},
	} {
		t.Run(tc.text, func(t *testing.T) {
			var name NwName
			require.NoError(t, name.UnmarshalBinary(tc.wire))
			require.Equal(t, tc.text, name.TextStr)
			wire, err := name.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, tc.wire, wire)
		})
	}
	for _, wire := range [][]byte{
		{0x90, 0x00}, {0x91, 0x00, 0x46}, {0x90, 0xd8, 0x00}, {0x90, 0xdf, 0xff},
	} {
		var name NwName
		require.Error(t, name.UnmarshalBinary(wire))
	}
	for _, text := range []string{"😀", string([]byte{0xff})} {
		name := NwName{CodeScheme: CodeScheme_UCS2, TextStr: text}
		_, err := name.MarshalBinary()
		require.Error(t, err)
	}
}
