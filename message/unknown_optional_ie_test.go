package message

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnknownOptionalIEs(t *testing.T) {
	t.Parallel()

	// R17: TS 24.501 sections 7.6.1 and 7.7.1, and TS 24.007 V17.3.0 sections 11.2.4
	// and 11.2.5. Each message is its header and mandatory IEs, then an optional IE
	// that the decoder knows.
	messages := []struct {
		name      string
		mandatory string
		known     string
		decode    func(b []byte) (knownDecoded bool, err error)
	}{
		{
			// Initial registration, a SUCI; the UE security capability.
			"RegReq", "7e004179" + "000c" + "0102f839f0ff000000004778", "2e02e060",
			func(b []byte) (bool, error) {
				var m RegReq
				err := m.UnmarshalBinary(b)
				return m.UESecCapability != nil, err
			},
		},
		{
			// ngKSI / service type, 5G-S-TMSI; the PDU session status.
			"SvcReq", "7e004c05" + "0007f401b000000091", "50020200",
			func(b []byte) (bool, error) {
				var m SvcReq
				err := m.UnmarshalBinary(b)
				return m.PDUSessStatus != nil, err
			},
		},
		{
			// Payload container type and payload container; the PDU session ID.
			"ULNASTransport", "7e006701" + "0008" + "2e0501c1ffff91a1", "1205",
			func(b []byte) (bool, error) {
				var m ULNASTransport
				err := m.UnmarshalBinary(b)
				return m.PDUSessID != nil, err
			},
		},
		{
			// No mandatory IE; the IMEISV.
			"SecModeComplete", "7e005e", "770009" + "15325476980000" + "01f1",
			func(b []byte) (bool, error) {
				var m SecModeComplete
				err := m.UnmarshalBinary(b)
				return m.IMEISV != nil, err
			},
		},
	}
	cases := []struct{ name, before, after string }{
		{"unknown TV", "e1", ""},
		{"unknown TLV", "5f022e00", ""},
		{"unknown TLV-E", "7a00022e00", ""},
		{"all three", "e15f022e007a00022e00", ""},
		{"incomplete unknown IE at the end", "", "7a00"},
	}
	for _, msg := range messages {
		for _, tc := range cases {
			t.Run(msg.name+" "+tc.name, func(t *testing.T) {
				packet, err := hex.DecodeString(msg.mandatory + tc.before + msg.known + tc.after)
				require.NoError(t, err)
				decoded, err := msg.decode(packet)
				require.NoError(t, err)
				require.True(t, decoded, "the known IE was not decoded")
			})
		}
		t.Run(msg.name+" unknown comprehension-required IE", func(t *testing.T) {
			packet, err := hex.DecodeString(msg.mandatory + "0f0100" + msg.known)
			require.NoError(t, err)
			_, err = msg.decode(packet)
			var required *UnknownRequiredIEError
			require.ErrorAs(t, err, &required)
			require.Equal(t, byte(0x0f), required.IEI)
		})
	}
}
