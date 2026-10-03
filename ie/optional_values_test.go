package ie

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/require"
)

type binaryValue interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func TestOptionalOpaqueValues(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		value    binaryValue
		min, max int
	}{
		{"emergency numbers", &EmergNumList{}, 3, 48},
		{"extended emergency numbers", &ExtendedEmergNumList{}, 4, 65535},
		{"SOR", &SORTransparentCntr{}, 17, 65535},
		{"operator categories", &OperatorDefinedAccessCategoryDefs{}, 0, 8320},
		{"EPS bearer status", &EPSBearerCtxStatus{}, 2, 2},
		{"EPS NAS algorithms", &EPSNASSecAlgos{}, 1, 1},
		{"S1 UE security", &S1UESecCapability{}, 2, 5},
		{"additional information", &AdditionalInfo{}, 1, 255},
		{"mapped EPS bearers", &MappedEPSBearerCtxs{}, 4, 65535},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.min > 0 {
				require.Error(t, tc.value.UnmarshalBinary(make([]byte, tc.min-1)))
			}
			require.Error(t, tc.value.UnmarshalBinary(make([]byte, tc.max+1)))
			for _, length := range []int{tc.min, tc.max} {
				input := make([]byte, length)
				if length > 0 {
					input[0] = 0x5a
				}
				require.NoError(t, tc.value.UnmarshalBinary(input))
				wire, err := tc.value.MarshalBinary()
				require.NoError(t, err)
				require.Len(t, wire, length)
				if length > 0 {
					require.Equal(t, input, wire)
					input[0], wire[0] = 0, 0
					again, err := tc.value.MarshalBinary()
					require.NoError(t, err)
					require.Equal(t, byte(0x5a), again[0], "decoder and encoder must own their buffers")
				}
			}
		})
	}
}

func TestOptionalSingleOctetValues(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value binaryValue
		mask  byte
	}{
		{"MICO", &MICOInd{}, 3},
		{"NSSAI inclusion", &NSSAIInclusionMode{}, 3},
		{"non-3GPP policies", &Non3GppNWProvidedPolicies{}, 1},
		{"SMS", &SMSInd{}, 1},
		{"always-on PDU", &AlwaysonPDUSessInd{}, 1},
		{"SSC modes", &AllowedSSCMode{}, 7},
		{"congestion re-attempt", &CongestionReattemptIndicator5GSM{}, 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Error(t, tc.value.UnmarshalBinary(nil))
			require.Error(t, tc.value.UnmarshalBinary([]byte{1, 2}))
			for bits := 0; bits < 256; bits++ {
				require.NoError(t, tc.value.UnmarshalBinary([]byte{byte(bits)}))
				wire, err := tc.value.MarshalBinary()
				require.NoError(t, err)
				require.Equal(t, []byte{byte(bits) & tc.mask}, wire, "emit zero spare bits")
			}
		})
	}
}

func TestOptionalValuesRejectInvalidEncoding(t *testing.T) {
	t.Parallel()
	for _, value := range []encoding.BinaryMarshaler{
		&EmergNumList{Contents: make([]byte, 2)},
		&EmergNumList{Contents: make([]byte, 49)},
		&ExtendedEmergNumList{Contents: make([]byte, 3)},
		&ExtendedEmergNumList{Contents: make([]byte, 65536)},
		&SORTransparentCntr{Contents: make([]byte, 16)},
		&SORTransparentCntr{Contents: make([]byte, 65536)},
		&OperatorDefinedAccessCategoryDefs{Contents: make([]byte, 8321)},
		&EPSBearerCtxStatus{Contents: make([]byte, 1)},
		&EPSBearerCtxStatus{Contents: make([]byte, 3)},
		&EPSNASSecAlgos{Contents: nil},
		&EPSNASSecAlgos{Contents: make([]byte, 2)},
		&S1UESecCapability{Contents: make([]byte, 1)},
		&S1UESecCapability{Contents: make([]byte, 6)},
		&AdditionalInfo{Contents: nil},
		&AdditionalInfo{Contents: make([]byte, 256)},
		&MappedEPSBearerCtxs{Contents: make([]byte, 3)},
		&MappedEPSBearerCtxs{Contents: make([]byte, 65536)},
		&NSSAIInclusionMode{NSSAIInclusionMode: 4},
		&Non3GppNWProvidedPolicies{Value: 2},
		&AllowedSSCMode{SSC1: 2},
		&AllowedSSCMode{SSC2: 2},
		&AllowedSSCMode{SSC3: 2},
		&CongestionReattemptIndicator5GSM{ABO: 2},
		&CongestionReattemptIndicator5GSM{CATBO: 2},
	} {
		wire, err := value.MarshalBinary()
		require.Error(t, err)
		require.Nil(t, wire)
	}
}
