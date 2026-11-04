package message

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/nas/ie"
)

func TestNotifRspUnmarshalBinary(t *testing.T) {
	t.Parallel()

	pduSessStatus := &ie.PDUSessStatus{
		Psi: ie.Psi{
			PSI: [16]bool{
				true, false, true, false, true, false, true, false,
				true, false, true, false, true, false, true, false,
			},
		},
	}

	testCases := []struct {
		name      string
		input     []byte
		expected  *NotifRsp
		expectErr bool
	}{
		{
			name: "Valid Notification Response with PDU Session Status",
			input: []byte{
				byte(Epd5GSMobilityMgmtMsg), // EPD
				byte(SecHdrTypePlainNas),    // Security Header
				byte(MsgTypeNotifRsp),       // Message Type
				NotifRspIEIPDUSessStatus,    // IEI
				0x02,                        // Length
				0x55, 0x55,                  // PDU Session Status content
			},
			expected: &NotifRsp{
				PDUSessStatus: pduSessStatus,
			},
			expectErr: false,
		},
		{
			name: "Valid Notification Response without optional IEs",
			input: []byte{
				byte(Epd5GSMobilityMgmtMsg),
				byte(SecHdrTypePlainNas),
				byte(MsgTypeNotifRsp),
			},
			expected:  &NotifRsp{},
			expectErr: false,
		},
		{
			name: "Invalid input - too short",
			input: []byte{
				byte(Epd5GSMobilityMgmtMsg),
				byte(SecHdrTypePlainNas),
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Unknown single octet IE is ignored",
			input: []byte{
				byte(Epd5GSMobilityMgmtMsg),
				byte(SecHdrTypePlainNas),
				byte(MsgTypeNotifRsp),
				0x99, // Unknown IEI
			},
			expected:  &NotifRsp{},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var notifRsp NotifRsp
			err := notifRsp.UnmarshalBinary(tc.input)

			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, &notifRsp)
		})
	}
}

func TestNotifRspOptionalIEHandling(t *testing.T) {
	t.Parallel()

	// R17: TS 24.501 sections 7.6.1, 7.6.3 and 7.7.1, and TS 24.007
	// V17.3.0 sections 11.2.4 and 11.4.2. Status is the expected value part.
	cases := []struct {
		name   string
		body   string
		status string
	}{
		{"issue 56 reproduction", "6e010050020000", "0000"},
		{"unknown TLV after status", "500201006e0100", "0100"},
		{"unknown TV before status", "9950020100", "0100"},
		{"unknown TV upper boundary", "ff50020100", "0100"},
		{"unknown TV lower boundary", "8050020100", "0100"},
		{"unknown TLV lower boundary", "10010050020100", "0100"},
		{"unknown TLV upper boundary", "6f010050020100", "0100"},
		{"unknown TLV-E lower boundary", "7000010050020100", "0100"},
		{"unknown TLV-E upper boundary", "7b00010050020100", "0100"},
		{"TLV-E two byte length", "710100" + strings.Repeat("6e", 256) + "50020100", "0100"},
		{"empty unknown TLV", "6e0050020100", "0100"},
		{"empty unknown TLV-E", "71000050020100", "0100"},
		{"mixed unknown IEs", "996e010071000100ff50020100", "0100"},
		{"empty optional status", "5000", ""},
		{"short optional status", "500100", ""},
		{"short status before unknown IE", "5001006e0100", ""},
		{"duplicate status keeps first", "5002010050020200", "0100"},
		{"invalid duplicate keeps first", "50020100500100", "0100"},
		{"invalid first status ignores duplicate", "50010050020200", ""},
		{"extended status ignores extra octets", "5021" + "0100" + strings.Repeat("ff", 31), "0100"},
		{"missing status length", "50", ""},
		{"truncated status value", "500201", ""},
		{"missing unknown TLV length", "6e", ""},
		{"truncated unknown TLV value", "6e0200", ""},
		{"missing unknown TLV-E length", "71", ""},
		{"short unknown TLV-E length", "7100", ""},
		{"truncated unknown TLV-E value", "71000200", ""},
		{"valid status before truncated duplicate", "5002010050", "0100"},
		{"valid status before truncated unknown", "500201007100", "0100"},
		{"do not rescan truncated value as IE", "6eff50020100", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := hex.DecodeString("7e0066" + tc.body)
			require.NoError(t, err)
			var direct NotifRsp
			require.NoError(t, direct.UnmarshalBinary(packet))
			parsed, err := Parse(packet, nil)
			require.NoError(t, err)
			require.IsType(t, &NotifRsp{}, parsed)
			require.Equal(t, &direct, parsed)
			if tc.status == "" {
				require.Nil(t, direct.PDUSessStatus)
				return
			}
			require.NotNil(t, direct.PDUSessStatus)
			value, err := direct.PDUSessStatus.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, tc.status, hex.EncodeToString(value))
		})
	}
}

func TestNotifRspUnknownRequiredIE(t *testing.T) {
	t.Parallel()

	// All R17 comprehension-required IEIs, including the boundary at 0x7c.
	for _, iei := range []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x7c, 0x7d, 0x7e, 0x7f,
	} {
		t.Run(fmt.Sprintf("IEI_%02x", iei), func(t *testing.T) {
			packet := []byte{0x7e, 0x00, 0x66, 0x50, 0x02, 0x01, 0x00, iei, 0x00, 0x00}
			var direct NotifRsp
			var required *UnknownRequiredIEError
			require.ErrorAs(t, direct.UnmarshalBinary(packet), &required)
			require.Equal(t, iei, required.IEI)
			parsed, err := Parse(packet, nil)
			require.Nil(t, parsed)
			require.ErrorAs(t, err, &required)
			require.Equal(t, iei, required.IEI)
		})
	}
}

func TestNotifRspUnmarshalReusedMessage(t *testing.T) {
	var response NotifRsp
	require.NoError(t, response.UnmarshalBinary([]byte{0x7e, 0x00, 0x66, 0x50, 0x02, 0x01, 0x00}))
	require.NotNil(t, response.PDUSessStatus)
	require.NoError(t, response.UnmarshalBinary([]byte{0x7e, 0x00, 0x66, 0x50, 0x01, 0x00}))
	require.Nil(t, response.PDUSessStatus)
}

func TestNotifRspMarshalBinary(t *testing.T) {
	t.Parallel()

	pduSessStatus := &ie.PDUSessStatus{
		Psi: ie.Psi{
			PSI: [16]bool{
				true, false, true, false, true, false, true, false,
				true, false, true, false, true, false, true, false,
			},
		},
	}

	testCases := []struct {
		name      string
		input     *NotifRsp
		expected  []byte
		expectErr bool
	}{
		{
			name: "Marshal Notification Response with PDU Session Status",
			input: &NotifRsp{
				PDUSessStatus: pduSessStatus,
			},
			expected: []byte{
				byte(Epd5GSMobilityMgmtMsg), // EPD
				byte(SecHdrTypePlainNas),    // Security Header
				byte(MsgTypeNotifRsp),       // Message Type
				NotifRspIEIPDUSessStatus,    // IEI
				0x02,                        // Length
				0x55, 0x55,                  // PDU Session Status content
			},
			expectErr: false,
		},
		{
			name:  "Marshal Notification Response without optional IEs",
			input: &NotifRsp{},
			expected: []byte{
				byte(Epd5GSMobilityMgmtMsg),
				byte(SecHdrTypePlainNas),
				byte(MsgTypeNotifRsp),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := tc.input.MarshalBinary()

			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, output)
		})
	}
}
