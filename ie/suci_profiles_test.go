package ie

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSUCIProfilesPreserveSchemeOutput(t *testing.T) {
	t.Parallel()
	for _, scheme := range []uint8{ECIESSchemeProfileA, ECIESSchemeProfileB} {
		t.Run(fmt.Sprint(scheme), func(t *testing.T) {
			keyLen := 32
			if scheme == ECIESSchemeProfileB {
				keyLen = 33
			}
			identity := MobileId5GS{
				TypeOfId: IdType_5GS_SUCI, SUPIFormat: uint8(SupiIMSI),
				PlmnId:             PlmnId{MCC: "208", MNC: "93"},
				RoutingIndDigit:    [4]uint8{1, 2, 3, 4},
				ProtectionSchemeId: scheme, HomeNwPubKeyId: 7,
				ECCEphPubKey: bytes.Repeat([]byte{0x42}, keyLen),
				CipherVal:    []byte{0x11, 0x22, 0x33, 0x44, 0x55},
				MACTag:       []byte{1, 2, 3, 4, 5, 6, 7, 8},
			}
			wire, err := identity.MarshalBinary()
			require.NoError(t, err)
			want := []byte{0x01, 0x02, 0xf8, 0x39, 0x21, 0x43, scheme, 7}
			want = append(want, identity.ECCEphPubKey...)
			want = append(want, identity.CipherVal...)
			want = append(want, identity.MACTag...)
			require.Equal(t, want, wire, "SUCI must not contain trailing padding")
			var decoded MobileId5GS
			require.NoError(t, decoded.UnmarshalBinary(wire))
			require.Equal(t, identity.ECCEphPubKey, decoded.ECCEphPubKey)
			require.Equal(t, identity.CipherVal, decoded.CipherVal)
			require.Equal(t, identity.MACTag, decoded.MACTag)
			for _, length := range []int{keyLen - 1, keyLen + 1} {
				identity.ECCEphPubKey = make([]byte, length)
				_, err = identity.MarshalBinary()
				require.Error(t, err)
			}
			identity.ECCEphPubKey = make([]byte, keyLen)
			for _, length := range []int{7, 9} {
				identity.MACTag = make([]byte, length)
				_, err = identity.MarshalBinary()
				require.Error(t, err)
			}
		})
	}
}
