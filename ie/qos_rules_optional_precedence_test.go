package ie

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// TS 24.501 9.11.4.13: a QoS rule that modifies an existing one may leave out its precedence and its
// QFI octet, and one that deletes packet filters lists their identifiers only.
func TestQosRulesOptionalPrecedence(t *testing.T) {
	t.Parallel()

	// QRI 1, "modify existing QoS rule and add packet filters", one bidirectional match-all packet filter
	// with identifier 1, no precedence; QRI 2, "modify existing QoS rule and delete packet filters" 1 and
	// 2, no precedence; QRI 3, the same for packet filter 1, with precedence 255 alone.
	in, err := hex.DecodeString("01000461310101" + "020003a20102" + "030003a101ff")
	require.NoError(t, err)
	var rules QosRules
	require.NoError(t, rules.UnmarshalBinary(in))
	require.Len(t, rules.Rules, 3)

	add := rules.Rules[0]
	require.Equal(t, OpCode_ModifyAddPktFilters, add.OpCode)
	require.True(t, add.NoPrecedence)
	require.True(t, add.NoQFI)
	require.Len(t, add.PktFilterList, 1)
	require.Equal(t, uint8(1), add.PktFilterList[0].Id)

	del := rules.Rules[1]
	require.Equal(t, OpCode_ModifyDelPktFilters, del.OpCode)
	require.True(t, del.NoPrecedence)
	require.Len(t, del.PktFilterList, 2)
	require.Equal(t, uint8(2), del.PktFilterList[1].Id)

	precedenceOnly := rules.Rules[2]
	require.False(t, precedenceOnly.NoPrecedence)
	require.True(t, precedenceOnly.NoQFI)
	require.Equal(t, uint8(255), precedenceOnly.Precedence)

	out, err := rules.MarshalBinary()
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(in), hex.EncodeToString(out))

	// A new rule has its precedence and its QFI.
	create, err := hex.DecodeString("01000421310101")
	require.NoError(t, err)
	require.Error(t, new(QosRules).UnmarshalBinary(create))
}
