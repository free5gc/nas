package nasType

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQoSRules(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name  string
		buf   []byte
		rules QoSRules
	}{
		{
			name: "Create_OneQoSRule_OnePacketFilter_OneComponent",
			buf: []byte{
				// QoS Rule 1, Length, operation(3) | DQR(1) | num of packet filter(4)
				0x00, 0x00, 0x0E, byte(OperationCodeCreateNewQoSRule)<<5 | 0<<4 | 1,
				// direction(2) | Packet filter 1(4), length
				byte(PacketFilterDirectionBidirectional)<<4 | 0, 0x09,
				0x10, 0xAC, 0x10, 0x06, 0x01, 0xFF, 0xFF, 0xFF, 0xFF,
				// QoS rule precedence, Segregation | QFI(6)
				0xFF, (0 << 6) | 0x00,
			},
			rules: QoSRules{
				{
					Identifier:  0,
					Operation:   OperationCodeCreateNewQoSRule,
					DQR:         false,
					Precedence:  255,
					Segregation: false,
					QFI:         0,
					PacketFilterList: PacketFilterList{
						PacketFilter{
							Identifier: 0,
							Direction:  PacketFilterDirectionBidirectional,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4RemoteAddress{
									Address: net.ParseIP("172.16.6.1").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0xFF},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Create_OneQoSRule_OnePacketFilter_MultipleComponent",
			buf: []byte{
				// QoS Rule 1, Length, operation(3) | DQR(1) | num of packet filter(4)
				0x01, 0x00, 0x13, byte(OperationCodeCreateNewQoSRule)<<5 | 1<<4 | 1,
				// direction(2) | Packet filter 1(4), length
				byte(PacketFilterDirectionBidirectional)<<4 | 1, 0x0E,
				0x11, 0x0A, 0x0A, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x41, 0x03, 0xE8, 0x0B, 0xB8,
				// QoS rule precedence, Segregation | QFI(6)
				0xFF, (0 << 6) | 0x09,
			},
			rules: QoSRules{
				{
					Identifier:  1,
					Operation:   OperationCodeCreateNewQoSRule,
					DQR:         true,
					Precedence:  255,
					Segregation: false,
					QFI:         9,
					PacketFilterList: PacketFilterList{
						PacketFilter{
							Identifier: 1,
							Direction:  PacketFilterDirectionBidirectional,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4LocalAddress{
									Address: net.ParseIP("10.10.0.0").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0x00},
								},
								&PacketFilterLocalPortRange{
									LowLimit:  1000,
									HighLimit: 3000,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Create_OneQoSRule_MultiplePacketFilter_MultipleComponent",
			buf: []byte{
				// QoS Rule 1, Length, operation(3) | DQR(1) | num of packet filter(4)
				0x01, 0x00, 0x21, byte(OperationCodeCreateNewQoSRule)<<5 | 1<<4 | 2,
				// direction(2) | Packet filter 1(4), length
				byte(PacketFilterDirectionUplink)<<4 | 1, 0x0E,
				0x11, 0x0A, 0x0A, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x41, 0x03, 0xE8, 0x0B, 0xB8,
				byte(PacketFilterDirectionDownlink)<<4 | 2, 0x0C,
				0x10, 0x0A, 0x0A, 0x64, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x50, 0x22, 0x65,
				// QoS rule precedence, Segregation | QFI(6)
				0xFF, (0 << 6) | 0x09,
			},
			rules: QoSRules{
				{
					Identifier:  1,
					Operation:   OperationCodeCreateNewQoSRule,
					DQR:         true,
					Precedence:  255,
					Segregation: false,
					QFI:         9,
					PacketFilterList: PacketFilterList{
						PacketFilter{
							Identifier: 1,
							Direction:  PacketFilterDirectionUplink,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4LocalAddress{
									Address: net.ParseIP("10.10.0.0").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0x00},
								},
								&PacketFilterLocalPortRange{
									LowLimit:  1000,
									HighLimit: 3000,
								},
							},
						},
						PacketFilter{
							Identifier: 2,
							Direction:  PacketFilterDirectionDownlink,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4RemoteAddress{
									Address: net.ParseIP("10.10.100.0").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0x00},
								},
								&PacketFilterSingleRemotePort{
									Value: 8805,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Create_MultipleQoSRule_MultiplePacketFilter_MultipleComponent",
			buf: []byte{
				// QoS Rule 1, Length, operation(3) | DQR(1) | num of packet filter(4)
				0x01, 0x00, 0x21, byte(OperationCodeCreateNewQoSRule)<<5 | 1<<4 | 2,
				// direction(2) | Packet filter 1(4), length
				byte(PacketFilterDirectionUplink)<<4 | 1, 0x0E,
				0x11, 0x0A, 0x0A, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x41, 0x03, 0xE8, 0x0B, 0xB8,
				byte(PacketFilterDirectionDownlink)<<4 | 2, 0x0C,
				0x10, 0x0A, 0x0A, 0x64, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x50, 0x22, 0x65,
				// QoS rule precedence, Segregation | QFI(6)
				0xFF, (0 << 6) | 0x09,

				// QoS Rule 2, Length, operation(3) | DQR(1) | num of packet filter(4)
				0x02, 0x00, 0x2a, byte(OperationCodeCreateNewQoSRule)<<5 | 0<<4 | 3,
				// direction(2) | Packet filter 1(4), length
				byte(PacketFilterDirectionUplink)<<4 | 1, 0x0E,
				0x11, 0x0A, 0x0A, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x41, 0x03, 0xE8, 0x0B, 0xB8,
				byte(PacketFilterDirectionDownlink)<<4 | 2, 0x0C,
				0x10, 0x0A, 0x0A, 0x64, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
				0x50, 0x22, 0x65,
				byte(PacketFilterDirectionBidirectional)<<4 | 3, 0x07,
				0x82, 0x11, 0x3E, 0x4D, 0x02, 0x1B, 0xCD,
				// QoS rule precedence, Segregation | QFI(6)
				0xFF, (0 << 6) | 0x05,
			},
			rules: QoSRules{
				{
					Identifier:  1,
					Operation:   OperationCodeCreateNewQoSRule,
					DQR:         true,
					Precedence:  255,
					Segregation: false,
					QFI:         9,
					PacketFilterList: PacketFilterList{
						PacketFilter{
							Identifier: 1,
							Direction:  PacketFilterDirectionUplink,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4LocalAddress{
									Address: net.ParseIP("10.10.0.0").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0x00},
								},
								&PacketFilterLocalPortRange{
									LowLimit:  1000,
									HighLimit: 3000,
								},
							},
						},
						PacketFilter{
							Identifier: 2,
							Direction:  PacketFilterDirectionDownlink,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4RemoteAddress{
									Address: net.ParseIP("10.10.100.0").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0x00},
								},
								&PacketFilterSingleRemotePort{
									Value: 8805,
								},
							},
						},
					},
				},
				{
					Identifier:  2,
					Operation:   OperationCodeCreateNewQoSRule,
					DQR:         false,
					Precedence:  255,
					Segregation: false,
					QFI:         5,
					PacketFilterList: PacketFilterList{
						PacketFilter{
							Identifier: 1,
							Direction:  PacketFilterDirectionUplink,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4LocalAddress{
									Address: net.ParseIP("10.10.0.0").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0x00},
								},
								&PacketFilterLocalPortRange{
									LowLimit:  1000,
									HighLimit: 3000,
								},
							},
						},
						PacketFilter{
							Identifier: 2,
							Direction:  PacketFilterDirectionDownlink,
							Components: PacketFilterComponentList{
								&PacketFilterIPv4RemoteAddress{
									Address: net.ParseIP("10.10.100.0").To4(),
									Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0x00},
								},
								&PacketFilterSingleRemotePort{
									Value: 8805,
								},
							},
						},
						PacketFilter{
							Identifier: 3,
							Direction:  PacketFilterDirectionBidirectional,
							Components: PacketFilterComponentList{
								&PacketFilterSourceMACAddress{
									MAC: net.HardwareAddr{0x11, 0x3E, 0x4D, 0x02, 0x1B, 0xCD},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Modify_OneQoSRule_DeletePacketFilter",
			buf: []byte{
				// QoS Rule 1, Length, operation(3) | DQR(1) | num of packet filter(4)
				0x01, 0x00, 0x08, byte(OperationCodeModifyExistingQoSRuleAndDeletePacketFilters)<<5 | 1<<4 | 5,
				// delete packet filter list
				0x01, 0x02, 0x03, 0x04, 0x05,
				// QoS rule precedence, Segregation | QFI(6)
				0xFF, (0 << 6) | 0x09,
			},
			rules: QoSRules{
				{
					Identifier:  1,
					Operation:   OperationCodeModifyExistingQoSRuleAndDeletePacketFilters,
					DQR:         true,
					Precedence:  255,
					Segregation: false,
					QFI:         9,
					PacketFilterList: PacketFilterList{
						PacketFilter{
							Identifier: 1,
						},
						PacketFilter{
							Identifier: 2,
						},
						PacketFilter{
							Identifier: 3,
						},
						PacketFilter{
							Identifier: 4,
						},
						PacketFilter{
							Identifier: 5,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf, err := tc.rules.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, tc.buf, buf)

			var rules QoSRules
			err = rules.UnmarshalBinary(tc.buf)
			require.NoError(t, err)
			require.Equal(t, tc.rules, rules)
		})
	}
}

func TestPacketFilterComponent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		id   PacketFilterComponentType
		buf  []byte
		pf   PacketFilterComponent
	}{
		{
			name: "MatchAll",
			id:   PacketFilterComponentTypeMatchAll,
			buf:  []byte{},
			pf:   &PacketFilterMatchAll{},
		},
		{
			name: "IPv4LocalAddress",
			id:   PacketFilterComponentTypeIPv4LocalAddress,
			buf:  []byte{0x0A, 0x0A, 0x00, 0x0B, 0xFF, 0xFF, 0xFF, 0xFF},
			pf: &PacketFilterIPv4LocalAddress{
				Address: net.ParseIP("10.10.0.11").To4(),
				Mask:    net.IPMask{0xFF, 0xFF, 0xFF, 0xFF},
			},
		},
		{
			name: "IPv4RemoteAddress",
			id:   PacketFilterComponentTypeIPv4RemoteAddress,
			buf:  []byte{0x0A, 0x0A, 0x00, 0x64, 0xFF, 0xFF, 0x00, 0x00},
			pf: &PacketFilterIPv4RemoteAddress{
				Address: net.ParseIP("10.10.0.100").To4(),
				Mask:    net.IPMask{0xFF, 0xFF, 0x00, 0x00},
			},
		},
		{
			name: "ProtocolID",
			id:   PacketFilterComponentTypeProtocolIdentifierOrNextHeader,
			buf:  []byte{0x0A},
			pf: &PacketFilterProtocolIdentifier{
				Value: 10,
			},
		},
		{
			name: "SingleLocalPort",
			id:   PacketFilterComponentTypeSingleLocalPort,
			buf:  []byte{0x1F, 0x40},
			pf: &PacketFilterSingleLocalPort{
				Value: 8000,
			},
		},
		{
			name: "SingleRemotePort",
			id:   PacketFilterComponentTypeSingleRemotePort,
			buf:  []byte{0xA6, 0x47},
			pf: &PacketFilterSingleRemotePort{
				Value: 42567,
			},
		},
		{
			name: "LocalPortRange",
			id:   PacketFilterComponentTypeLocalPortRange,
			buf:  []byte{0x03, 0xE8, 0x1F, 0x40},
			pf: &PacketFilterLocalPortRange{
				LowLimit:  1000,
				HighLimit: 8000,
			},
		},
		{
			name: "RemotePortRange",
			id:   PacketFilterComponentTypeRemotePortRange,
			buf:  []byte{0x03, 0xE8, 0x1F, 0x40},
			pf: &PacketFilterRemotePortRange{
				LowLimit:  1000,
				HighLimit: 8000,
			},
		},
		{
			name: "SecurityParameterIndex",
			id:   PacketFilterComponentTypeSecurityParameterIndex,
			buf:  []byte{0x00, 0x13, 0xBE, 0x79},
			pf: &PacketFilterSecurityParameterIndex{
				Index: 1293945,
			},
		},
		{
			name: "ServiceClass",
			id:   PacketFilterComponentTypeTypeOfServiceOrTrafficClass,
			buf:  []byte{0x01, 0x7A},
			pf: &PacketFilterServiceClass{
				Class: 1,
				Mask:  122,
			},
		},
		{
			name: "FlowLabel",
			id:   PacketFilterComponentTypeFlowLabel,
			buf:  []byte{0x04, 0x93, 0xE0},
			pf: &PacketFilterFlowLabel{
				Label: 300000,
			},
		},
		{
			name: "DestinationMACAddress",
			id:   PacketFilterComponentTypeDestinationMACAddress,
			buf:  []byte{0x00, 0x1B, 0x44, 0x31, 0xAB, 0xC6},
			pf: &PacketFilterDestinationMACAddress{
				MAC: net.HardwareAddr{0x00, 0x1B, 0x44, 0x31, 0xAB, 0xC6},
			},
		},
		{
			name: "SourceMACAddress",
			id:   PacketFilterComponentTypeSourceMACAddress,
			buf:  []byte{0x31, 0xAB, 0xC6, 0x00, 0x1B, 0x44},
			pf: &PacketFilterSourceMACAddress{
				MAC: net.HardwareAddr{0x31, 0xAB, 0xC6, 0x00, 0x1B, 0x44},
			},
		},
		{
			name: "8021Q_CTAG_VID",
			id:   PacketFilterComponentType8021Q_CTAG_VID,
			buf:  []byte{0x28, 0x35},
			pf: &PacketFilterCTagVID{
				VID: 10293,
			},
		},
		{
			name: "8021Q_STAG_VID",
			id:   PacketFilterComponentType8021Q_STAG_VID,
			buf:  []byte{0x4F, 0xAA},
			pf: &PacketFilterSTagVID{
				VID: 20394,
			},
		},
		{
			name: "8021Q_CTAG_PCPOrDEI",
			id:   PacketFilterComponentType8021Q_CTAG_PCPOrDEI,
			buf:  []byte{0x04},
			pf: &PacketFilterCTagPCPDEI{
				Value: 4,
			},
		},
		{
			name: "8021Q_STAG_PCPOrDEI",
			id:   PacketFilterComponentType8021Q_STAG_PCPOrDEI,
			buf:  []byte{0x01},
			pf: &PacketFilterSTagPCPDEI{
				Value: 1,
			},
		},
		{
			name: "EtherType",
			id:   PacketFilterComponentTypeEthertype,
			buf:  []byte{0xAB, 0x74},
			pf: &PacketFilterEtherType{
				EtherType: 43892,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf, err := tc.pf.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, tc.buf, buf)
			require.Equal(t, tc.pf.Length(), len(buf))

			pfc := newPacketFilterComponent(tc.id)
			require.Equal(t, tc.id, pfc.Type())
			err = pfc.UnmarshalBinary(tc.buf)
			require.NoError(t, err)
			require.Equal(t, tc.pf, pfc)
		})
	}
}
