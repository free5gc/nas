package ie

import "testing"

func TestPacketFilterContentsRejectsOutOfRangePorts(t *testing.T) {
	tests := []struct {
		name   string
		remote string
		local  string
	}{
		{"remote single", "65536", ""},
		{"remote range", "80-65536", ""},
		{"local single", "", "65536"},
		{"local range", "", "65536-80"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PacketFilterContents{
				RemoteAddr:      "any",
				LocalAddr:       "assigned",
				RemotePortRange: tt.remote,
				LocalPortRange:  tt.local,
			}

			if _, err := p.MarshalBinary(); err == nil {
				t.Fatal("expected error for out-of-range port")
			}
		})
	}
}

func TestPacketFilterContentsAcceptsValidPorts(t *testing.T) {
	tests := []struct {
		name   string
		remote string
		local  string
	}{
		{"remote zero", "0", ""},
		{"remote max", "65535", ""},
		{"remote range", "80-65535", ""},
		{"local zero", "", "0"},
		{"local max", "", "65535"},
		{"local range", "", "80-65535"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PacketFilterContents{
				RemoteAddr:      "any",
				LocalAddr:       "assigned",
				RemotePortRange: tt.remote,
				LocalPortRange:  tt.local,
			}

			if _, err := p.MarshalBinary(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
