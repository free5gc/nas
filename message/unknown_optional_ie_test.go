package message

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
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
		tlvE      string // the IEI of the unknown TLV-E IE, when 0x7a is known
	}{
		{
			// Initial registration, a SUCI; the UE security capability.
			"RegReq", "7e004179" + "000c" + "0102f839f0ff000000004778", "2e02e060",
			func(b []byte) (bool, error) {
				var m RegReq
				err := m.UnmarshalBinary(b)
				return m.UESecCapability != nil, err
			}, "",
		},
		{
			// ngKSI / service type, 5G-S-TMSI; the PDU session status.
			"SvcReq", "7e004c05" + "0007f401b000000091", "50020200",
			func(b []byte) (bool, error) {
				var m SvcReq
				err := m.UnmarshalBinary(b)
				return m.PDUSessStatus != nil, err
			}, "",
		},
		{
			// Payload container type and payload container; the PDU session ID.
			"ULNASTransport", "7e006701" + "0008" + "2e0501c1ffff91a1", "1205",
			func(b []byte) (bool, error) {
				var m ULNASTransport
				err := m.UnmarshalBinary(b)
				return m.PDUSessID != nil, err
			}, "",
		},
		{
			// No mandatory IE; the IMEISV.
			"SecModeComplete", "7e005e", "770009" + "15325476980000" + "01f1",
			func(b []byte) (bool, error) {
				var m SecModeComplete
				err := m.UnmarshalBinary(b)
				return m.IMEISV != nil, err
			}, "",
		},
		{
			// No mandatory IE; the authentication response parameter.
			"AuthRsp", "7e0057", "2d10" + strings.Repeat("00", 16),
			func(b []byte) (bool, error) {
				var m AuthRsp
				err := m.UnmarshalBinary(b)
				return m.AuthRspParam != nil, err
			}, "",
		},
		{
			// 5GMM cause #21; the authentication failure parameter.
			"AuthFailure", "7e0059" + "15", "300e" + strings.Repeat("00", 14),
			func(b []byte) (bool, error) {
				var m AuthFailure
				err := m.UnmarshalBinary(b)
				return m.AuthFailureParam != nil, err
			}, "",
		},
		{
			// No mandatory IE; the SOR transparent container.
			"RegComplete", "7e0043", "730011" + strings.Repeat("00", 17),
			func(b []byte) (bool, error) {
				var m RegComplete
				err := m.UnmarshalBinary(b)
				return m.SORTransparentCntr != nil, err
			}, "",
		},
		{
			// Control plane service type and ngKSI; the PDU session status.
			"CtrlPlaneSvcReq", "7e004f" + "00", "50020000",
			func(b []byte) (bool, error) {
				var m CtrlPlaneSvcReq
				err := m.UnmarshalBinary(b)
				return m.PDUSessStatus != nil, err
			}, "",
		},
		{
			// Integrity protection maximum data rate; the PDU session type.
			"PDUSessEstReq", "2e0101c1" + "ffff", "91",
			func(b []byte) (bool, error) {
				var m PDUSessEstReq
				err := m.UnmarshalBinary(b)
				return m.PDUSessType != nil, err
			}, "",
		},
		{
			// No mandatory IE; the 5GSM capability.
			"PDUSessModReq", "2e0101c9", "280100",
			func(b []byte) (bool, error) {
				var m PDUSessModReq
				err := m.UnmarshalBinary(b)
				return m.Capability5GSM != nil, err
			}, "78",
		},
		{
			// 5GSM cause; the extended protocol configuration options.
			"PDUSessModCmdRej", "2e0101cd" + "1a", "7b000180",
			func(b []byte) (bool, error) {
				var m PDUSessModCmdRej
				err := m.UnmarshalBinary(b)
				return m.ExtendedProtCfgOpts != nil, err
			}, "",
		},
		{
			// No mandatory IE; the extended protocol configuration options.
			"PDUSessModComplete", "2e0101cc", "7b000180",
			func(b []byte) (bool, error) {
				var m PDUSessModComplete
				err := m.UnmarshalBinary(b)
				return m.ExtendedProtCfgOpts != nil, err
			}, "",
		},
		{
			// No mandatory IE; the 5GSM cause.
			"PDUSessRelReq", "2e0101d1", "5924",
			func(b []byte) (bool, error) {
				var m PDUSessRelReq
				err := m.UnmarshalBinary(b)
				return m.Cause5GSM != nil, err
			}, "",
		},
		{
			// No mandatory IE; the 5GSM cause.
			"PDUSessRelComplete", "2e0101d4", "5924",
			func(b []byte) (bool, error) {
				var m PDUSessRelComplete
				err := m.UnmarshalBinary(b)
				return m.Cause5GSM != nil, err
			}, "",
		},
		{
			// An EAP message; the extended protocol configuration options.
			"PDUSessAuthComplete", "2e0101c6" + "0004" + "02010004", "7b000180",
			func(b []byte) (bool, error) {
				var m PDUSessAuthComplete
				err := m.UnmarshalBinary(b)
				return m.ExtendedProtCfgOpts != nil, err
			}, "",
		},
		{
			// No mandatory IE; the remote UE context connected, which the decoder reports as not
			// implemented: it reaches it.
			"RemoteUEReport", "2e0101da", "760010" + strings.Repeat("00", 16),
			func(b []byte) (bool, error) {
				var m RemoteUEReport
				err := m.UnmarshalBinary(b)
				var todo *Error
				if errors.As(err, &todo) && bytes.Equal(todo.IEToDoList, []byte{0x76}) {
					return true, nil
				}
				return false, err
			}, "",
		},
	}
	for _, msg := range messages {
		tlvE := msg.tlvE
		if tlvE == "" {
			tlvE = "7a"
		}
		cases := []struct{ name, before, after string }{
			{"unknown TV", "e1", ""},
			{"unknown TLV", "5f022e00", ""},
			{"unknown TLV-E", tlvE + "00022e00", ""},
			{"all three", "e15f022e00" + tlvE + "00022e00", ""},
			{"incomplete unknown IE at the end", "", tlvE + "00"},
		}
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
