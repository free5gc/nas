package ie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPayloadContainerBoundsAndReuse(t *testing.T) {
	t.Parallel()
	var payload PayloadCntr
	for _, length := range []int{0, 65536} {
		require.Error(t, payload.UnmarshalBinary(make([]byte, length), PayloadCntrType_N1SMInfo))
		payload.Contents = make([]byte, length)
		_, err := payload.MarshalBinary()
		require.Error(t, err)
	}
	for _, length := range []int{1, 65535} {
		input := make([]byte, length)
		input[0] = 0x5a
		require.NoError(t, payload.UnmarshalBinary(input, PayloadCntrType_N1SMInfo))
		input[0] = 0
		wire, err := payload.MarshalBinary()
		require.NoError(t, err)
		require.Len(t, wire, length)
		require.Equal(t, byte(0x5a), wire[0])
		wire[0] = 0
		require.Equal(t, byte(0x5a), payload.Contents[0])
	}
	require.NoError(t, payload.UnmarshalBinary([]byte{3}, PayloadCntrType_N1SMInfo))
	require.Equal(t, []byte{3}, payload.Contents)
	for _, pct := range []uint8{PayloadCntrType_MultiplePayloads, PayloadCntrType_EventNotif} {
		require.Error(t, payload.UnmarshalBinary([]byte{1}, pct))
		payload.Pct = pct
		_, err := payload.MarshalBinary()
		require.Error(t, err)
	}
}
