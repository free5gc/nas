package ie

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/require"
)

// TS 24.007 11.4.2: a type 4 IE whose length is greater than its value part needs is no syntactical error, and the
// extra octets are ignored.
func TestLongerValueIgnoresExtraOctets(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value []byte // the longest value part
		new   func() encoding.BinaryUnmarshaler
	}{
		{
			"5GMM capability",
			[]byte{0x07, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01},
			func() encoding.BinaryUnmarshaler { return new(Capability5GMM) },
		},
		{
			"UE security capability",
			[]byte{0xf0, 0x70, 0xe0, 0xe0, 0, 0, 0, 0},
			func() encoding.BinaryUnmarshaler { return new(UESecCapability) },
		},
		{"UE status", []byte{0x02}, func() encoding.BinaryUnmarshaler { return new(UEStatus) }},
		{"UE's usage setting", []byte{0x01}, func() encoding.BinaryUnmarshaler { return new(UesUsageSetting) }},
		{"5GS DRX parameters", []byte{0x02}, func() encoding.BinaryUnmarshaler { return new(DRXParams5GS) }},
		{"5GS update type", []byte{0x01}, func() encoding.BinaryUnmarshaler { return new(UpdateType5GS) }},
		{"GPRS timer 3", []byte{0x21}, func() encoding.BinaryUnmarshaler { return new(GPRSTimer3) }},
		{"UE request type", []byte{0x01}, func() encoding.BinaryUnmarshaler { return new(UEReqType) }},
		{
			"extended DRX parameters", // octet 4, the Extended Paging Time Window, is not modelled
			[]byte{0x35, 0x01},
			func() encoding.BinaryUnmarshaler { return new(ExtendedDRXParams) },
		},
		{
			"S-NSSAI",
			[]byte{0x01, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x02},
			func() encoding.BinaryUnmarshaler { return new(SNSSAI) },
		},
		{
			"EPS bearer context status",
			[]byte{0x20, 0x00},
			func() encoding.BinaryUnmarshaler { return new(EPSBearerCtxStatus) },
		},
		{
			"S1 UE security capability",
			[]byte{0xe0, 0xe0, 0x00, 0x00, 0x00},
			func() encoding.BinaryUnmarshaler { return new(S1UESecCapability) },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := tc.new()
			require.NoError(t, want.UnmarshalBinary(tc.value))
			got := tc.new()
			require.NoError(t, got.UnmarshalBinary(append(append([]byte{}, tc.value...), 0xff, 0x00)))
			require.Equal(t, want, got)
		})
	}
}
