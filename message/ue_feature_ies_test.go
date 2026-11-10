package message

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/nas/ie"
)

// The optional IEs that a UE sends for a feature that it supports (TS 24.501
// 8.2.6, 8.2.10, 8.2.16), each in a message with the mandatory IEs only: the
// message decodes, and encodes again to the same octets.
func TestOptionalIEsOfUEFeatures(t *testing.T) {
	t.Parallel()

	regReq := func(i ...byte) []byte {
		return append([]byte{
			0x7e, 0x00, 0x41, // EPD, security header type, REGISTRATION REQUEST
			0x79, // ngKSI "no key is available", initial registration with follow-on request
			// 5GS mobile identity: a SUCI of the null scheme, PLMN 208/93, MSIN 0000000001
			0x00, 0x0d, 0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10,
		}, i...)
	}
	svcReq := func(i ...byte) []byte {
		return append([]byte{
			0x7e, 0x00, 0x4c, // EPD, security header type, SERVICE REQUEST
			0x10, // service type "data", ngKSI 0
			// 5G-S-TMSI: AMF set ID 1, AMF pointer 0, 5G-TMSI 1
			0x00, 0x07, 0xf4, 0x00, 0x40, 0x00, 0x00, 0x00, 0x01,
		}, i...)
	}
	ulNASTransport := func(i ...byte) []byte {
		return append([]byte{
			0x7e, 0x00, 0x67, // EPD, security header type, UL NAS TRANSPORT
			0x01,                         // payload container type "N1 SM information"
			0x00, 0x04, 0x2e, 0x09, 0x01, // payload container: a PDU SESSION ESTABLISHMENT REQUEST
			0xc1,
			0x12, 0x09, // PDU session ID 9
		}, i...)
	}
	type nasMessage interface {
		encoding.BinaryMarshaler
		encoding.BinaryUnmarshaler
	}

	for _, tc := range []struct {
		name string
		msg  nasMessage
		wire []byte
		ie   func(nasMessage) any
		want any
	}{
		{
			"registration request, mobile station classmark 2 (9.11.3.31C)",
			&RegReq{}, regReq(0x41, 0x03, 0x57, 0x58, 0xa6),
			func(m nasMessage) any { return m.(*RegReq).MobileStationClassmark2 },
			&ie.MobileStationClassmark2{Contents: []byte{0x57, 0x58, 0xa6}},
		},
		{
			// TS 24.007 11.4.2: a length greater than the value part is no syntactical error.
			"registration request, mobile station classmark 2 of 4 octets",
			&RegReq{}, regReq(0x41, 0x04, 0x57, 0x58, 0xa6, 0x00),
			func(m nasMessage) any { return m.(*RegReq).MobileStationClassmark2 },
			&ie.MobileStationClassmark2{Contents: []byte{0x57, 0x58, 0xa6, 0x00}},
		},
		{
			"registration request, supported codecs (9.11.3.51A)",
			&RegReq{}, regReq(0x42, 0x04, 0x04, 0x02, 0x60, 0x04),
			func(m nasMessage) any { return m.(*RegReq).SupportedCodecs },
			&ie.SupportedCodecList{Contents: []byte{0x04, 0x02, 0x60, 0x04}},
		},
		{
			"registration request, UE radio capability ID (9.11.3.68)",
			&RegReq{}, regReq(0x67, 0x05, 0x01, 0x23, 0x45, 0x67, 0x89),
			func(m nasMessage) any { return m.(*RegReq).UERadioCapabilityID },
			&ie.UERadioCapabilityID{Contents: []byte{0x01, 0x23, 0x45, 0x67, 0x89}},
		},
		{
			"registration request, requested mapped NSSAI (9.11.3.31B)",
			&RegReq{}, regReq(0x35, 0x05, 0x04, 0x01, 0x01, 0x02, 0x03),
			func(m nasMessage) any { return m.(*RegReq).ReqMappedNSSAI },
			&ie.MappedNSSAI{Contents: []byte{0x04, 0x01, 0x01, 0x02, 0x03}},
		},
		{
			"registration request, additional information requested (9.11.3.12A)",
			&RegReq{}, regReq(0x48, 0x01, 0x01),
			func(m nasMessage) any { return m.(*RegReq).AdditionalInfoReq },
			&ie.AdditionalInfoReq{Cipherkey: 1},
		},
		{
			"registration request, requested WUS assistance information (9.11.3.71)",
			&RegReq{}, regReq(0x1a, 0x01, 0x05),
			func(m nasMessage) any { return m.(*RegReq).ReqWUSAssistanceInfo },
			&ie.WUSAssistanceInfo{Contents: []byte{0x05}},
		},
		{
			"registration request, N5GC indication (9.11.3.72)",
			&RegReq{}, regReq(0xa1),
			func(m nasMessage) any { return m.(*RegReq).N5GCInd },
			&ie.N5GCInd{N5GCREG: true},
		},
		{
			"registration request, requested NB-N1 mode DRX parameters (9.11.3.73)",
			&RegReq{}, regReq(0x30, 0x01, 0x02),
			func(m nasMessage) any { return m.(*RegReq).ReqNBN1ModeDRXParams },
			&ie.NBN1ModeDRXParams{NBN1ModeDRXValue: 2},
		},
		{
			"registration request, paging restriction (9.11.3.77)",
			&RegReq{}, regReq(0x28, 0x01, 0x01),
			func(m nasMessage) any { return m.(*RegReq).PagingRestriction },
			&ie.PagingRestriction{PagingRestrictionType: 1},
		},
		{
			"registration request, service-level-AA container (9.11.2.10)",
			&RegReq{}, regReq(0x72, 0x00, 0x04, 0x10, 0x02, 0x01, 0x02),
			func(m nasMessage) any { return m.(*RegReq).SvcLvlAACntr },
			&ie.SvcLvlAACntr{Contents: []byte{0x10, 0x02, 0x01, 0x02}},
		},
		{
			"registration request, NID (9.11.3.79)",
			&RegReq{}, regReq(0x32, 0x06, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba),
			func(m nasMessage) any { return m.(*RegReq).NID },
			&ie.NID{Contents: []byte{0x10, 0x32, 0x54, 0x76, 0x98, 0xba}},
		},
		{
			"registration request, requested PEIPS assistance information (9.11.3.80)",
			&RegReq{}, regReq(0x2a, 0x01, 0x05),
			func(m nasMessage) any { return m.(*RegReq).ReqPEIPSAssistanceInfo },
			&ie.PEIPSAssistanceInfo{Contents: []byte{0x05}},
		},
		{
			"service request, paging restriction with PDU sessions (9.11.3.77)",
			&SvcReq{}, svcReq(0x28, 0x03, 0x03, 0x20, 0x01),
			func(m nasMessage) any { return m.(*SvcReq).PagingRestriction },
			&ie.PagingRestriction{PagingRestrictionType: 3, Psi: ie.Psi{PSI: [16]bool{5: true, 8: true}}},
		},
		{
			"UL NAS transport, MA PDU session information (9.11.3.31A)",
			&ULNASTransport{}, ulNASTransport(0xa1),
			func(m nasMessage) any { return m.(*ULNASTransport).MAPDUSessInfo },
			&ie.MAPDUSessInfo{MAPDUSessInfoValue: 1},
		},
		{
			"UL NAS transport, release assistance indication (9.11.3.46A)",
			&ULNASTransport{}, ulNASTransport(0xf1),
			func(m nasMessage) any { return m.(*ULNASTransport).RelAssistanceInd },
			&ie.RelAssistanceInd{DDX: 1},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.NoError(t, tc.msg.UnmarshalBinary(tc.wire))
			require.Equal(t, tc.want, tc.ie(tc.msg))
			wire, err := tc.msg.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, tc.wire, wire)
		})
	}
}
