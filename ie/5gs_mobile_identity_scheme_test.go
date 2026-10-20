package ie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// 3GPP TS 33.501 Annex C standardises three SUCI protection schemes: the null
// scheme (0x0), ECIES profile A (0x1) and ECIES profile B (0x2). 0x3-0xB are
// reserved for future standardised schemes and 0xC-0xF for home-operator
// proprietary schemes, so a SUCI carrying one of them cannot be parsed and must be
// rejected while the identity is being decoded rather than being forwarded.
//
// The buffer below is the ECIES profile A vector from the existing positive cases
// with the protection scheme nibble replaced, so only the scheme ID differs
// between cases.
func TestMobileId5GS_SUCIReservedProtectionScheme(t *testing.T) {
	t.Parallel()

	// Same SUCI as the profile A positive case, parameterised by scheme ID.
	newSUCI := func(schemeID uint8) []byte {
		b := []byte{
			0x01, 0x12, 0x34, 0x56, 0x78, 0x90, 0x01, 0x12,
			0x2b, 0xc3, 0xc1, 0xba, 0xef, 0xf3, 0xcb, 0x71, 0x58, 0xca, 0x0e, 0x35, 0x70, 0x35, 0x9f, 0xdf,
			0x9b, 0x11, 0xcc, 0xe0, 0x59, 0xa9, 0x0c, 0xef, 0xef, 0xda, 0x6b, 0x9d, 0xa1, 0x7a, 0x5f, 0x0b,
			0xb8, 0x6f, 0xb0, 0xfe, 0x25, 0x29, 0xd7, 0x2d, 0x14, 0x20, 0x0a, 0x1d, 0x8a,
		}
		// byte 6 packs routing indicator digits, the scheme ID and the
		// home network public key identifier. Set only the scheme ID nibble.
		b[6] = (b[6] & 0xf0) | (schemeID & 0x0f)
		return b
	}

	// 0x3-0xB reserved for future standardised schemes.
	reservedStandard := []uint8{0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b}
	// 0xC-0xF reserved for home-operator proprietary schemes.
	reservedProprietary := []uint8{0x0c, 0x0d, 0x0e, 0x0f}

	for _, id := range reservedStandard {
		t.Run("reserved-standard-"+hexID(id), func(t *testing.T) {
			t.Parallel()
			ie := new(MobileId5GS)
			err := ie.UnmarshalBinary(newSUCI(id))
			require.Error(t, err)
			require.Contains(t, err.Error(), "unsupported protection scheme ID")
		})
	}

	for _, id := range reservedProprietary {
		t.Run("reserved-proprietary-"+hexID(id), func(t *testing.T) {
			t.Parallel()
			ie := new(MobileId5GS)
			err := ie.UnmarshalBinary(newSUCI(id))
			require.Error(t, err)
			require.Contains(t, err.Error(), "unsupported protection scheme ID")
		})
	}
}

// The three standardised scheme IDs must keep parsing: a guard against the
// rejection being too broad.
func TestMobileId5GS_SUCIStandardProtectionSchemesStillParse(t *testing.T) {
	t.Parallel()

	// Null scheme, from the existing positive cases: SchemeOutput is an MSIN in BCD.
	nullSchemeSUCI := []byte{
		0x01, 0x12, 0x34, 0x56, 0x78, 0x90, 0x00, 0x12,
		0x01, 0x23, 0x45, 0x67, 0x89,
	}

	// ECIES profile A, from the existing positive cases.
	profileASUCI := []byte{
		0x01, 0x12, 0x34, 0x56, 0x78, 0x90, 0x01, 0x12,
		0x2b, 0xc3, 0xc1, 0xba, 0xef, 0xf3, 0xcb, 0x71, 0x58, 0xca, 0x0e, 0x35, 0x70, 0x35, 0x9f, 0xdf,
		0x9b, 0x11, 0xcc, 0xe0, 0x59, 0xa9, 0x0c, 0xef, 0xef, 0xda, 0x6b, 0x9d, 0xa1, 0x7a, 0x5f, 0x0b,
		0xb8, 0x6f, 0xb0, 0xfe, 0x25, 0x29, 0xd7, 0x2d, 0x14, 0x20, 0x0a, 0x1d, 0x8a,
	}

	// ECIES profile B needs at least 41 octets of scheme output; the same header with
	// the scheme ID set to 2 and a long enough body.
	profileBSUCI := append([]byte{
		0x01, 0x12, 0x34, 0x56, 0x78, 0x90, 0x02, 0x12,
	}, make([]byte, 48)...)

	for _, tc := range []struct {
		name   string
		id     uint8
		buffer []byte
	}{
		{"null", NullScheme, nullSchemeSUCI},
		{"profile-a", ECIESSchemeProfileA, profileASUCI},
		{"profile-b", ECIESSchemeProfileB, profileBSUCI},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ie := new(MobileId5GS)
			require.NoError(t, ie.UnmarshalBinary(tc.buffer))
			require.Equal(t, tc.id, ie.ProtectionSchemeId)
		})
	}
}

func hexID(id uint8) string {
	const hexDigits = "0123456789abcdef"
	return string([]byte{hexDigits[id>>4], hexDigits[id&0x0f]})
}
