package message

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/nas/ie"
)

var (
	countTestKnasEnc = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	countTestKnasInt = []byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
)

func newCountTestSecCtx(side Side, cipheringAlg ie.AlgCiphering, integrityAlg ie.AlgIntegrity) *SecCtx {
	return NewSecCtx(side, Bearer3GPP, cipheringAlg, integrityAlg, countTestKnasEnc, countTestKnasInt)
}

func setCount(c *Count, count uint32) {
	c.Set(uint16(count>>8), uint8(count))
}

// uplinkMsg protects a REGISTRATION COMPLETE with the UE uplink COUNT count.
func uplinkMsg(t *testing.T, cipheringAlg ie.AlgCiphering, count uint32, st SecHdrType) []byte {
	t.Helper()
	ue := newCountTestSecCtx(UESide, cipheringAlg, AlgIntegrity128NIA2)
	setCount(ue.UplinkCount, count)
	b, err := Marshal(&RegComplete{}, ue, st)
	require.NoError(t, err)
	return b
}

// acceptedAt returns a core network context that has accepted uplink COUNT count.
func acceptedAt(cipheringAlg ie.AlgCiphering, count uint32) *SecCtx {
	sc := newCountTestSecCtx(CoreNetworkSide, cipheringAlg, AlgIntegrity128NIA2)
	setCount(sc.UplinkCount, count)
	sc.UplinkCount.MarkAccepted()
	return sc
}

func requireMACFailure(t *testing.T, msg Message, err error) {
	t.Helper()
	var nasErr *Error
	require.ErrorAs(t, err, &nasErr)
	require.NotNil(t, nasErr.MACFailure)
	require.NotNil(t, msg)
}

func requireAccepted(t *testing.T, c *Count, count uint32) {
	t.Helper()
	require.Equal(t, count, c.Get())
	require.True(t, c.accepted)
}

func TestParseReplayedCountFailsIntegrityCheck(t *testing.T) {
	rx := acceptedAt(AlgCiphering128NEA0, 99)
	p := uplinkMsg(t, AlgCiphering128NEA0, 100, SecHdrTypeIntegrityProtected)

	msg, err := Parse(p, rx)
	require.NoError(t, err)
	require.IsType(t, &RegComplete{}, msg)
	requireAccepted(t, rx.UplinkCount, 100)

	// The replay is estimated at 0x200+100 and fails the MAC check. An
	// integrity-only message is still returned with the MAC failure.
	msg, err = Parse(p, rx)
	requireMACFailure(t, msg, err)
	requireAccepted(t, rx.UplinkCount, 100)

	_, err = Parse(uplinkMsg(t, AlgCiphering128NEA0, 101, SecHdrTypeIntegrityProtected), rx)
	require.NoError(t, err)
	requireAccepted(t, rx.UplinkCount, 101)
}

func TestParseReplayedCipheredCountIsRejected(t *testing.T) {
	rx := acceptedAt(AlgCiphering128NEA2, 99)
	p := uplinkMsg(t, AlgCiphering128NEA2, 100, SecHdrTypeIntegrityProtectedAndCiphered)

	msg, err := Parse(p, rx)
	require.NoError(t, err)
	require.IsType(t, &RegComplete{}, msg)

	_, err = Parse(p, rx)
	require.Error(t, err)
	requireAccepted(t, rx.UplinkCount, 100)
}

func TestParseLowerCountFailsIntegrityCheck(t *testing.T) {
	rx := acceptedAt(AlgCiphering128NEA0, 100)

	msg, err := Parse(uplinkMsg(t, AlgCiphering128NEA0, 99, SecHdrTypeIntegrityProtected), rx)
	requireMACFailure(t, msg, err)
	requireAccepted(t, rx.UplinkCount, 100)
}

func TestParseOverflowEstimate(t *testing.T) {
	rx := acceptedAt(AlgCiphering128NEA0, 0x1ff)

	_, err := Parse(uplinkMsg(t, AlgCiphering128NEA0, 0x200, SecHdrTypeIntegrityProtected), rx)
	require.NoError(t, err)
	requireAccepted(t, rx.UplinkCount, 0x200)

	// Same SQN as the stored COUNT: the estimate takes the next overflow value.
	_, err = Parse(uplinkMsg(t, AlgCiphering128NEA0, 0x300, SecHdrTypeIntegrityProtected), rx)
	require.NoError(t, err)
	requireAccepted(t, rx.UplinkCount, 0x300)
}

func TestParseNewContextAcceptsStartValueOnce(t *testing.T) {
	smcComplete := uplinkMsg(t, AlgCiphering128NEA2, 0, SecHdrTypeIntegrityProtectedAndCipheredWithNew5gNasSecCtx)

	rx := newCountTestSecCtx(CoreNetworkSide, AlgCiphering128NEA2, AlgIntegrity128NIA2)
	_, err := Parse(smcComplete, rx)
	require.NoError(t, err)
	requireAccepted(t, rx.UplinkCount, 0)

	_, err = Parse(smcComplete, rx)
	require.Error(t, err)
	requireAccepted(t, rx.UplinkCount, 0)

	// The new-context header does not reset an accepted receive COUNT.
	rx = acceptedAt(AlgCiphering128NEA2, 10)
	_, err = Parse(smcComplete, rx)
	require.Error(t, err)
	requireAccepted(t, rx.UplinkCount, 10)

	// Set, as done by the AMF when it sends the SECURITY MODE COMMAND, starts
	// a new receive COUNT.
	rx.UplinkCount.Set(0, 0)
	_, err = Parse(smcComplete, rx)
	require.NoError(t, err)
	requireAccepted(t, rx.UplinkCount, 0)
}

func TestParseMarkAccepted(t *testing.T) {
	p := uplinkMsg(t, AlgCiphering128NEA0, 100, SecHdrTypeIntegrityProtected)

	notAccepted := newCountTestSecCtx(CoreNetworkSide, AlgCiphering128NEA0, AlgIntegrity128NIA2)
	setCount(notAccepted.UplinkCount, 100)
	_, err := Parse(p, notAccepted)
	require.NoError(t, err)

	marked := acceptedAt(AlgCiphering128NEA0, 100)
	msg, err := Parse(p, marked)
	requireMACFailure(t, msg, err)
	requireAccepted(t, marked.UplinkCount, 100)
}

func TestParseNIA0HasNoReplayCheck(t *testing.T) {
	p := uplinkMsg(t, AlgCiphering128NEA0, 5, SecHdrTypeIntegrityProtected)
	rx := newCountTestSecCtx(CoreNetworkSide, AlgCiphering128NEA0, AlgIntegrity128NIA0)

	for i := 0; i < 2; i++ {
		_, err := Parse(p, rx)
		require.NoError(t, err)
		require.Equal(t, uint32(5), rx.UplinkCount.Get())
	}

	// With NIA0 the COUNT may wrap (TS 24.501 4.4.3.5).
	setCount(rx.UplinkCount, 0xffffff)
	rx.UplinkCount.MarkAccepted()
	_, err := Parse(uplinkMsg(t, AlgCiphering128NEA0, 3, SecHdrTypeIntegrityProtected), rx)
	require.NoError(t, err)
	require.Equal(t, uint32(3), rx.UplinkCount.Get())
}

func TestParseRejectsCountWrap(t *testing.T) {
	rx := acceptedAt(AlgCiphering128NEA0, 0xffffff)

	for _, count := range []uint32{0, 5} {
		msg, err := Parse(uplinkMsg(t, AlgCiphering128NEA0, count, SecHdrTypeIntegrityProtected), rx)
		requireMACFailure(t, msg, err)
		requireAccepted(t, rx.UplinkCount, 0xffffff)
	}

	// A receive COUNT installed with Set and not yet accepted does not wrap either.
	rx = newCountTestSecCtx(CoreNetworkSide, AlgCiphering128NEA0, AlgIntegrity128NIA2)
	setCount(rx.UplinkCount, 0xffff80)
	msg, err := Parse(uplinkMsg(t, AlgCiphering128NEA0, 0x10, SecHdrTypeIntegrityProtected), rx)
	requireMACFailure(t, msg, err)
	require.Equal(t, uint32(0xffff80), rx.UplinkCount.Get())
	require.False(t, rx.UplinkCount.accepted)

	_, err = Parse(uplinkMsg(t, AlgCiphering128NEA0, 0xffff80, SecHdrTypeIntegrityProtected), rx)
	require.NoError(t, err)
	requireAccepted(t, rx.UplinkCount, 0xffff80)
}

func TestParseKeepsCountOnMacCalculateError(t *testing.T) {
	rx := acceptedAt(AlgCiphering128NEA0, 100)
	rx.IntegrityAlg = 0xff

	_, err := Parse(uplinkMsg(t, AlgCiphering128NEA0, 101, SecHdrTypeIntegrityProtected), rx)
	require.Error(t, err)
	requireAccepted(t, rx.UplinkCount, 100)
}

func TestParseUESideReplay(t *testing.T) {
	amf := newCountTestSecCtx(CoreNetworkSide, AlgCiphering128NEA0, AlgIntegrity128NIA2)
	setCount(amf.DownlinkCount, 7)
	p, err := Marshal(&RegComplete{}, amf, SecHdrTypeIntegrityProtected)
	require.NoError(t, err)

	ue := newCountTestSecCtx(UESide, AlgCiphering128NEA0, AlgIntegrity128NIA2)
	setCount(ue.DownlinkCount, 6)
	ue.DownlinkCount.MarkAccepted()

	_, err = Parse(p, ue)
	require.NoError(t, err)
	requireAccepted(t, ue.DownlinkCount, 7)

	msg, err := Parse(p, ue)
	requireMACFailure(t, msg, err)
	requireAccepted(t, ue.DownlinkCount, 7)
	require.Equal(t, uint32(0), ue.UplinkCount.Get())
}

func TestCloneCopiesAccepted(t *testing.T) {
	clone := acceptedAt(AlgCiphering128NEA0, 100).Clone()

	msg, err := Parse(uplinkMsg(t, AlgCiphering128NEA0, 100, SecHdrTypeIntegrityProtected), clone)
	requireMACFailure(t, msg, err)
}
