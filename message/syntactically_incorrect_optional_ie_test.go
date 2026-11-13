package message

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/nas/ie"
)

// TS 24.501 section 7.7.1: the network treats a syntactically incorrect optional IE as not
// present, and decodes the rest of the message.

func TestRegReqSyntacticallyIncorrectOptionalIE(t *testing.T) {
	t.Parallel()

	// Initial registration, ngKSI 7, a SUCI; a requested NSSAI whose S-NSSAI has the
	// reserved length 3 (TS 24.501 9.11.2.8); the UE security capability after it.
	packet, err := hex.DecodeString("7e004179" + "000c0102f839f0ff000000004778" + "2f0403010203" + "2e02e060")
	require.NoError(t, err)
	var m RegReq
	require.NoError(t, m.UnmarshalBinary(packet))
	require.Nil(t, m.ReqNSSAI, "the requested NSSAI is not treated as not present")
	require.NotNil(t, m.UESecCapability, "the IE after the incorrect one was not decoded")
}

func TestSvcReqSyntacticallyIncorrectOptionalIE(t *testing.T) {
	t.Parallel()

	// ngKSI / service type, 5G-S-TMSI; an uplink data status of one octet (TS 24.501
	// 9.11.3.57 has two); the PDU session status after it.
	packet, err := hex.DecodeString("7e004c05" + "0007f401b000000091" + "400100" + "50020200")
	require.NoError(t, err)
	var m SvcReq
	require.NoError(t, m.UnmarshalBinary(packet))
	require.Nil(t, m.UplinkDataStatus, "the uplink data status is not treated as not present")
	require.NotNil(t, m.PDUSessStatus, "the IE after the incorrect one was not decoded")
}

func TestULNASTransportSyntacticallyIncorrectOptionalIE(t *testing.T) {
	t.Parallel()

	// Payload container type and payload container; an S-NSSAI of the reserved length 3
	// (TS 24.501 9.11.2.8); the PDU session ID after it.
	packet, err := hex.DecodeString("7e006701" + "00082e0501c1ffff91a1" + "2203010203" + "1205")
	require.NoError(t, err)
	var m ULNASTransport
	require.NoError(t, m.UnmarshalBinary(packet))
	require.Nil(t, m.SNSSAI, "the S-NSSAI is not treated as not present")
	require.NotNil(t, m.PDUSessID, "the IE after the incorrect one was not decoded")
}

func TestSecModeCompleteSyntacticallyIncorrectOptionalIE(t *testing.T) {
	t.Parallel()

	// An IMEISV of one octet (TS 24.501 9.11.3.4 has nine); the NAS message container
	// after it.
	packet, err := hex.DecodeString("7e005e" + "77000115" + "710003" + "7e0041")
	require.NoError(t, err)
	var m SecModeComplete
	require.NoError(t, m.UnmarshalBinary(packet))
	require.Nil(t, m.IMEISV, "the IMEISV is not treated as not present")
	require.NotNil(t, m.NASMsgCntr, "the IE after the incorrect one was not decoded")
}

// TS 24.007 section 11.4.2: a type 4 IE whose length is greater than its value part needs is no
// syntactical error; the network reads it and ignores the extra octets.
func TestRegReqLongerOptionalIE(t *testing.T) {
	t.Parallel()

	// Initial registration, ngKSI 7, a SUCI; requested DRX parameters of two octets (TS 24.501
	// 9.11.3.2A has one), DRX cycle parameter T = 64; the UE security capability after it.
	packet, err := hex.DecodeString("7e004179" + "000c0102f839f0ff000000004778" + "51020200" + "2e02e060")
	require.NoError(t, err)
	var m RegReq
	require.NoError(t, m.UnmarshalBinary(packet))
	require.NotNil(t, m.ReqDRXParams, "the requested DRX parameters are not read")
	require.Equal(t, ie.DRXCycleParameterT64, m.ReqDRXParams.Value)
	require.NotNil(t, m.UESecCapability, "the IE after the longer one was not decoded")
}
