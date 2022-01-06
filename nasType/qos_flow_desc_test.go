package nasType

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQoSFlowDescs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		buf   []byte
		descs QoSFlowDescs
	}{
		{
			name: "OneDesc",
			buf: []byte{
				0x09, 0x20, 1<<6 | 0x05,
				0x01, 0x01, 0x05,
				0x02, 0x03, 0x06, 0x00, 0x64,
				0x03, 0x03, 0x06, 0x00, 0x0A,
				0x04, 0x03, 0x0B, 0x01, 0x2C,
				0x05, 0x03, 0x0B, 0x00, 0x1E,
			},
			descs: QoSFlowDescs{
				QoSFlowDesc{
					QFI:           9,
					OperationCode: OperationCodeCreateNewQoSFlowDescription,
					Parameters: QoSFlowParameterList{
						&QoSFlow5QI{
							FiveQI: 5,
						},
						&QoSFlowGFBRUplink{
							Unit:  QoSFlowBitRateUnit1Mbps,
							Value: 100,
						},
						&QoSFlowGFBRDownlink{
							Unit:  QoSFlowBitRateUnit1Mbps,
							Value: 10,
						},
						&QoSFlowMFBRUplink{
							Unit:  QoSFlowBitRateUnit1Gbps,
							Value: 300,
						},
						&QoSFlowMFBRDownlink{
							Unit:  QoSFlowBitRateUnit1Gbps,
							Value: 30,
						},
					},
				},
			},
		},
		{
			name: "MultipleDesc",
			buf: []byte{
				0x09, 0x20, 1<<6 | 0x05,
				0x01, 0x01, 0x05,
				0x02, 0x03, 0x06, 0x00, 0x64,
				0x03, 0x03, 0x06, 0x00, 0x0A,
				0x04, 0x03, 0x0B, 0x01, 0x2C,
				0x05, 0x03, 0x0B, 0x00, 0x1E,
				0x10, 0x60, 1<<6 | 0x06,
				0x01, 0x01, 0x10,
				0x02, 0x03, 0x06, 0x00, 0x64,
				0x03, 0x03, 0x06, 0x00, 0x0A,
				0x04, 0x03, 0x0B, 0x01, 0x2C,
				0x05, 0x03, 0x0B, 0x00, 0x1E,
				0x06, 0x02, 0x03, 0xE8,
			},
			descs: QoSFlowDescs{
				QoSFlowDesc{
					QFI:           9,
					OperationCode: OperationCodeCreateNewQoSFlowDescription,
					Parameters: QoSFlowParameterList{
						&QoSFlow5QI{
							FiveQI: 5,
						},
						&QoSFlowGFBRUplink{
							Unit:  QoSFlowBitRateUnit1Mbps,
							Value: 100,
						},
						&QoSFlowGFBRDownlink{
							Unit:  QoSFlowBitRateUnit1Mbps,
							Value: 10,
						},
						&QoSFlowMFBRUplink{
							Unit:  QoSFlowBitRateUnit1Gbps,
							Value: 300,
						},
						&QoSFlowMFBRDownlink{
							Unit:  QoSFlowBitRateUnit1Gbps,
							Value: 30,
						},
					},
				},
				QoSFlowDesc{
					QFI:           16,
					OperationCode: OperationCodeModifyExistingQoSFlowDescription,
					Parameters: QoSFlowParameterList{
						&QoSFlow5QI{
							FiveQI: 16,
						},
						&QoSFlowGFBRUplink{
							Unit:  QoSFlowBitRateUnit1Mbps,
							Value: 100,
						},
						&QoSFlowGFBRDownlink{
							Unit:  QoSFlowBitRateUnit1Mbps,
							Value: 10,
						},
						&QoSFlowMFBRUplink{
							Unit:  QoSFlowBitRateUnit1Gbps,
							Value: 300,
						},
						&QoSFlowMFBRDownlink{
							Unit:  QoSFlowBitRateUnit1Gbps,
							Value: 30,
						},
						&QoSFlowAveragingWindow{
							AverageWindow: 1000,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := tc.descs.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, tc.buf, ret)

			var descs QoSFlowDescs
			err = descs.UnmarshalBinary(tc.buf)
			require.NoError(t, err)
			require.Equal(t, tc.descs, descs)
		})
	}
}
