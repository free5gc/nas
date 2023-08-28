package nasConvert_test

import (
	"testing"
	"time"

	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasType"
	"github.com/stretchr/testify/require"
)

var (
	CST *time.Location
	IST *time.Location
	EST *time.Location
	CET *time.Location
)

func TestMain(m *testing.M) {
	CST, _ = time.LoadLocation("Asia/Taipei")
	IST, _ = time.LoadLocation("Asia/Kolkata")
	EST, _ = time.LoadLocation("America/New_York") // If using DST, EST would be EDT.
	CET, _ = time.LoadLocation("Europe/Berlin")    // If using DST, CET would be CEST.

	m.Run()
}

func TestUniversalTimeAndLocalTimeZoneToNas(t *testing.T) {
	tests := []struct {
		in  nasType.UniversalTimeAndLocalTimeZone
		out nasType.UniversalTimeAndLocalTimeZone
	}{
		{
			in: nasConvert.EncodeUniversalTimeAndLocalTimeZoneToNas(time.Date(2023, time.July, 13, 12, 27, 39, 0, CST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x32), uint8(0x70), uint8(0x31), uint8(0x21), uint8(0x72), uint8(0x93), uint8(0x23),
				},
			},
		},
		{
			in: nasConvert.EncodeUniversalTimeAndLocalTimeZoneToNas(time.Date(2019, time.December, 15, 16, 55, 46, 0, IST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x91), uint8(0x21), uint8(0x51), uint8(0x61), uint8(0x55), uint8(0x64), uint8(0x22),
				},
			},
		},
		{
			in: nasConvert.EncodeUniversalTimeAndLocalTimeZoneToNas(time.Date(2001, time.February, 2, 9, 3, 6, 0, EST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x10), uint8(0x20), uint8(0x20), uint8(0x90), uint8(0x30), uint8(0x60), uint8(0x0A),
				},
			},
		},
		{
			in: nasConvert.EncodeUniversalTimeAndLocalTimeZoneToNas(time.Date(2023, time.August, 24, 9, 18, 43, 0, CET)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x32), uint8(0x80), uint8(0x42), uint8(0x90), uint8(0x81), uint8(0x34), uint8(0x80),
				},
			},
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.out, tc.in)
	}
}

func TestDecodeUniversalTimeAndLocalTimeZone(t *testing.T) {
	tests := []struct {
		in  time.Time
		out time.Time
	}{
		{
			in: nasConvert.DecodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x32), uint8(0x70), uint8(0x31), uint8(0x21), uint8(0x72), uint8(0x93), uint8(0x23),
				},
			}),
			out: time.Date(2023, time.July, 13, 12, 27, 39, 0, CST),
		},
		{
			in: nasConvert.DecodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x91), uint8(0x21), uint8(0x51), uint8(0x61), uint8(0x55), uint8(0x64), uint8(0x22),
				},
			}),
			out: time.Date(2019, time.December, 15, 16, 55, 46, 0, IST),
		},
		{
			in: nasConvert.DecodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x10), uint8(0x20), uint8(0x20), uint8(0x90), uint8(0x30), uint8(0x60), uint8(0x0A),
				},
			}),
			out: time.Date(2001, time.February, 2, 9, 3, 6, 0, EST),
		},
		{
			in: nasConvert.DecodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x32), uint8(0x80), uint8(0x42), uint8(0x90), uint8(0x81), uint8(0x34), uint8(0x80),
				},
			}),
			out: time.Date(2023, time.August, 24, 9, 18, 43, 0, CET),
		},
	}

	for _, testData := range tests {
		require.Equal(t, testData.out.Format(time.RFC822Z), testData.in.Format(time.RFC822Z))
	}
}

func TestLocalTimeZoneToNas(t *testing.T) {
	tests := []struct {
		in  nasType.LocalTimeZone
		out nasType.LocalTimeZone
	}{
		{
			in: nasConvert.EncodeLocalTimeZoneToNas("+08:30"),
			out: nasType.LocalTimeZone{
				Octet: uint8(0x43),
			},
		},
		{
			in: nasConvert.EncodeLocalTimeZoneToNas("-04:45"),
			out: nasType.LocalTimeZone{
				Octet: uint8(0x99),
			},
		},
		{
			in: nasConvert.EncodeLocalTimeZoneToNas("+10:45"),
			out: nasType.LocalTimeZone{
				Octet: uint8(0x34),
			},
		},
		{
			in: nasConvert.EncodeLocalTimeZoneToNas("+01:00+1"), // CEST
			out: nasType.LocalTimeZone{
				Octet: uint8(0x80),
			},
		},
		{
			in: nasConvert.EncodeLocalTimeZoneToNas("-05:00+1"), // EDT
			out: nasType.LocalTimeZone{
				Octet: uint8(0x69),
			},
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.out, tc.in)
	}
}

func TestDecodeLocalTimeZone(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in: nasConvert.DecodeLocalTimeZone(nasType.LocalTimeZone{
				Octet: uint8(0x23),
			}),
			out: "+08:00",
		},
		{
			in: nasConvert.DecodeLocalTimeZone(nasType.LocalTimeZone{
				Octet: uint8(0x99),
			}),
			out: "-04:45",
		},
		{
			in: nasConvert.DecodeLocalTimeZone(nasType.LocalTimeZone{
				Octet: uint8(0x80),
			}),
			out: "+02:00",
		},
		{
			in: nasConvert.DecodeLocalTimeZone(nasType.LocalTimeZone{
				Octet: uint8(0x69),
			}),
			out: "-04:00",
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.out, tc.in)
	}
}

func TestDaylightSavingTimeToNas(t *testing.T) {
	nasConvertNetworkDaylightSavingTimeTable := []struct {
		in  nasType.NetworkDaylightSavingTime
		out nasType.NetworkDaylightSavingTime
	}{
		{
			in: nasConvert.EncodeDaylightSavingTimeToNas("-05:00+1"), // EST to EDT
			out: nasType.NetworkDaylightSavingTime{
				Len:   uint8(0x01),
				Octet: uint8(0x01),
			},
		},
		{
			in: nasConvert.EncodeDaylightSavingTimeToNas("+08:00+2"),
			out: nasType.NetworkDaylightSavingTime{
				Len:   uint8(0x01),
				Octet: uint8(0x02),
			},
		},
		{
			in: nasConvert.EncodeDaylightSavingTimeToNas("-03:00"),
			out: nasType.NetworkDaylightSavingTime{
				Len:   uint8(0x01),
				Octet: uint8(0x00),
			},
		},
	}

	for _, tc := range nasConvertNetworkDaylightSavingTimeTable {
		require.Equal(t, tc.out, tc.in)
	}
}

func TestDecodeDaylightSavingTime(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  nasConvert.DecodeDaylightSavingTime(nasConvert.EncodeDaylightSavingTimeToNas("-05:00+1")),
			out: "+1",
		},
		{
			in:  nasConvert.DecodeDaylightSavingTime(nasConvert.EncodeDaylightSavingTimeToNas("+08:00+2")),
			out: "+2",
		},
		{
			in:  nasConvert.DecodeDaylightSavingTime(nasConvert.EncodeDaylightSavingTimeToNas("03:00")),
			out: "",
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.in, tc.out)
	}
}
