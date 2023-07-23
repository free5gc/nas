package nasConvert_test

import (
	"fmt"
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
)

func TestMain(m *testing.M) {
	CST, _ = time.LoadLocation("Asia/Taipei")
	IST, _ = time.LoadLocation("Asia/Kolkata")
	EST, _ = time.LoadLocation("America/New_York") // If using DST, EST would be EDT.

	m.Run()
}

func TestUniversalTimeAndLocalTimeZoneToNas(t *testing.T) {
	tests := []struct {
		in  nasType.UniversalTimeAndLocalTimeZone
		out nasType.UniversalTimeAndLocalTimeZone
	}{
		{
			in: nasConvert.UniversalTimeAndLocalTimeZoneToNas(time.Date(2023, time.July, 13, 12, 27, 39, 0, CST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x32), uint8(0x70), uint8(0x31), uint8(0x21), uint8(0x72), uint8(0x93), uint8(0x23),
				},
			},
		},
		{
			in: nasConvert.UniversalTimeAndLocalTimeZoneToNas(time.Date(2019, time.December, 15, 16, 55, 46, 0, IST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x91), uint8(0x21), uint8(0x51), uint8(0x61), uint8(0x55), uint8(0x64), uint8(0x22),
				},
			},
		},
		{
			in: nasConvert.UniversalTimeAndLocalTimeZoneToNas(time.Date(2001, time.February, 2, 9, 3, 6, 0, EST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x10), uint8(0x20), uint8(0x20), uint8(0x90), uint8(0x30), uint8(0x60), uint8(0x0A),
				},
			},
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.out, tc.in)
	}
}

func getTimeZoneOffset(timezone uint8) int {
	octet := int((timezone >> 4) + (timezone&0x07)*10)
	offset := (octet / 4 * 60 * 60) + (octet % 4 * 15 * 60)
	if timezone&0x08 == 0x08 {
		// sign is "-"
		offset = 0 - offset
	}

	return offset
}

func decodeUniversalTimeAndLocalTimeZone(
	nasUniversalTimeAndLocalTimeZone nasType.UniversalTimeAndLocalTimeZone,
) time.Time {
	year := 2000 + int((nasUniversalTimeAndLocalTimeZone.GetYear()&0x0f)*10+
		((nasUniversalTimeAndLocalTimeZone.GetYear()&0xf0)>>4))

	month := int((nasUniversalTimeAndLocalTimeZone.GetMonth()&0x0f)*10 +
		((nasUniversalTimeAndLocalTimeZone.GetMonth() & 0xf0) >> 4))

	day := int((nasUniversalTimeAndLocalTimeZone.GetDay()&0x0f)*10 +
		((nasUniversalTimeAndLocalTimeZone.GetDay() & 0xf0) >> 4))

	hour := int((nasUniversalTimeAndLocalTimeZone.GetHour()&0x0f)*10 +
		((nasUniversalTimeAndLocalTimeZone.GetHour() & 0xf0) >> 4))

	minute := int((nasUniversalTimeAndLocalTimeZone.GetMinute()&0x0f)*10 +
		((nasUniversalTimeAndLocalTimeZone.GetMinute() & 0xf0) >> 4))

	second := int((nasUniversalTimeAndLocalTimeZone.GetSecond()&0x0f)*10 +
		((nasUniversalTimeAndLocalTimeZone.GetSecond() & 0xf0) >> 4))

	offset := getTimeZoneOffset(nasUniversalTimeAndLocalTimeZone.GetTimeZone())
	location := time.FixedZone("NameIsNotImportant", offset)

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, location)
}

func TestDecodeUniversalTimeAndLocalTimeZone(t *testing.T) {
	tests := []struct {
		in  time.Time
		out time.Time
	}{
		{
			in: decodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x32), uint8(0x70), uint8(0x31), uint8(0x21), uint8(0x72), uint8(0x93), uint8(0x23),
				},
			}),
			out: time.Date(2023, time.July, 13, 12, 27, 39, 0, CST),
		},
		{
			in: decodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x91), uint8(0x21), uint8(0x51), uint8(0x61), uint8(0x55), uint8(0x64), uint8(0x22),
				},
			}),
			out: time.Date(2019, time.December, 15, 16, 55, 46, 0, IST),
		},
		{
			in: decodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Octet: [7]uint8{
					uint8(0x10), uint8(0x20), uint8(0x20), uint8(0x90), uint8(0x30), uint8(0x60), uint8(0x0A),
				},
			}),
			out: time.Date(2001, time.February, 2, 9, 3, 6, 0, EST),
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
			in: nasConvert.LocalTimeZoneToNas("+08:30"),
			out: nasType.LocalTimeZone{
				Octet: uint8(0x43),
			},
		},
		{
			in: nasConvert.LocalTimeZoneToNas("-04:45"),
			out: nasType.LocalTimeZone{
				Octet: uint8(0x99),
			},
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.out, tc.in)
	}
}

func decodeLocalTimeZone(nasLocalTimeZone nasType.LocalTimeZone) string {
	offset := getTimeZoneOffset(nasLocalTimeZone.GetTimeZone())
	timezone := ""

	if offset < 0 {
		timezone += "-"
		offset = 0 - offset
	} else {
		timezone += "+"
	}

	timezone += fmt.Sprintf("%02d:%02d", offset/3600, (offset%3600)/60)
	return timezone
}

func TestDecodeLocalTimeZone(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in: decodeLocalTimeZone(nasType.LocalTimeZone{
				Octet: uint8(0x23),
			}),
			out: "+08:00",
		},
		{
			in: decodeLocalTimeZone(nasType.LocalTimeZone{
				Octet: uint8(0x99),
			}),
			out: "-04:45",
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
			in: nasConvert.DaylightSavingTimeToNas("-05:00+1"), // EST to EDT
			out: nasType.NetworkDaylightSavingTime{
				Len:   uint8(0x01),
				Octet: uint8(0x01),
			},
		},
		{
			in: nasConvert.DaylightSavingTimeToNas("+08:00+2"),
			out: nasType.NetworkDaylightSavingTime{
				Len:   uint8(0x01),
				Octet: uint8(0x02),
			},
		},
		{
			in: nasConvert.DaylightSavingTimeToNas("-03:00"),
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

func decodeDaylightSavingTime(nasDaylightSavingTime nasType.NetworkDaylightSavingTime) string {
	result := ""

	switch nasDaylightSavingTime.Getvalue() {
	case uint8(0x00):
		result = ""
	case uint8(0x01):
		result = "+1"
	case uint8(0x02):
		result = "+2"
	}

	return result
}

func TestDecodeDaylightSavingTime(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  decodeDaylightSavingTime(nasConvert.DaylightSavingTimeToNas("-05:00+1")),
			out: "+1",
		},
		{
			in:  decodeDaylightSavingTime(nasConvert.DaylightSavingTimeToNas("+08:00+2")),
			out: "+2",
		},
		{
			in:  decodeDaylightSavingTime(nasConvert.DaylightSavingTimeToNas("03:00")),
			out: "",
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.in, tc.out)
	}
}
