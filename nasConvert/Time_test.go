package nasConvert_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/stretchr/testify/assert"
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

type nasConvertUniversalTimeAndLocalTimeZone struct {
	in  nasType.UniversalTimeAndLocalTimeZone
	out nasType.UniversalTimeAndLocalTimeZone
}

func TestUniversalTimeAndLocalTimeZoneToNas(t *testing.T) {
	nasConvertUniversalTimeAndLocalTimeZoneTable := []nasConvertUniversalTimeAndLocalTimeZone{
		{
			in: nasConvert.UniversalTimeAndLocalTimeZoneToNas(time.Date(2023, time.July, 13, 12, 27, 39, 0, CST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Iei: nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType,
				Octet: [7]uint8{
					uint8(0x32), uint8(0x70), uint8(0x31), uint8(0x21), uint8(0x72), uint8(0x93), uint8(0x23),
				},
			},
		},
		{
			in: nasConvert.UniversalTimeAndLocalTimeZoneToNas(time.Date(2019, time.December, 15, 16, 55, 46, 0, IST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Iei: nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType,
				Octet: [7]uint8{
					uint8(0x91), uint8(0x21), uint8(0x51), uint8(0x61), uint8(0x55), uint8(0x64), uint8(0x22),
				},
			},
		},
		{
			in: nasConvert.UniversalTimeAndLocalTimeZoneToNas(time.Date(2001, time.February, 2, 9, 3, 6, 0, EST)),
			out: nasType.UniversalTimeAndLocalTimeZone{
				Iei: nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType,
				Octet: [7]uint8{
					uint8(0x10), uint8(0x20), uint8(0x20), uint8(0x90), uint8(0x30), uint8(0x60), uint8(0x0A),
				},
			},
		},
	}

	for _, testData := range nasConvertUniversalTimeAndLocalTimeZoneTable {
		assert.Equal(t, testData.out.GetYear(), testData.in.GetYear())
		assert.Equal(t, testData.out.GetMonth(), testData.in.GetMonth())
		assert.Equal(t, testData.out.GetDay(), testData.in.GetDay())
		assert.Equal(t, testData.out.GetHour(), testData.in.GetHour())
		assert.Equal(t, testData.out.GetMinute(), testData.in.GetMinute())
		assert.Equal(t, testData.out.GetSecond(), testData.in.GetSecond())
		assert.Equal(t, testData.out.GetTimeZone(), testData.in.GetTimeZone())
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

type testDecodedUniversalTimeAndLocalTimeZone struct {
	in  time.Time
	out time.Time
}

func TestDecodeUniversalTimeAndLocalTimeZone(t *testing.T) {
	testDecodedUniversalTimeAndLocalTimeZoneTable := []testDecodedUniversalTimeAndLocalTimeZone{
		{
			in: decodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Iei: nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType,
				Octet: [7]uint8{
					uint8(0x32), uint8(0x70), uint8(0x31), uint8(0x21), uint8(0x72), uint8(0x93), uint8(0x23),
				},
			}),
			out: time.Date(2023, time.July, 13, 12, 27, 39, 0, CST),
		},
		{
			in: decodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Iei: nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType,
				Octet: [7]uint8{
					uint8(0x91), uint8(0x21), uint8(0x51), uint8(0x61), uint8(0x55), uint8(0x64), uint8(0x22),
				},
			}),
			out: time.Date(2019, time.December, 15, 16, 55, 46, 0, IST),
		},
		{
			in: decodeUniversalTimeAndLocalTimeZone(nasType.UniversalTimeAndLocalTimeZone{
				Iei: nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType,
				Octet: [7]uint8{
					uint8(0x10), uint8(0x20), uint8(0x20), uint8(0x90), uint8(0x30), uint8(0x60), uint8(0x0A),
				},
			}),
			out: time.Date(2001, time.February, 2, 9, 3, 6, 0, EST),
		},
	}

	for _, testData := range testDecodedUniversalTimeAndLocalTimeZoneTable {
		assert.Equal(t, testData.out.Format(time.RFC822Z), testData.in.Format(time.RFC822Z))
	}
}

type nasConvertLocalTimeZone struct {
	in  nasType.LocalTimeZone
	out nasType.LocalTimeZone
}

func TestLocalTimeZoneToNas(t *testing.T) {
	nasConvertLocalTimeZoneTable := []nasConvertLocalTimeZone{
		{
			in: nasConvert.LocalTimeZoneToNas("+08:30"),
			out: nasType.LocalTimeZone{
				Iei:   nasMessage.ConfigurationUpdateCommandLocalTimeZoneType,
				Octet: uint8(0x43),
			},
		},
		{
			in: nasConvert.LocalTimeZoneToNas("-04:45"),
			out: nasType.LocalTimeZone{
				Iei:   nasMessage.ConfigurationUpdateCommandLocalTimeZoneType,
				Octet: uint8(0x99),
			},
		},
	}

	for _, testData := range nasConvertLocalTimeZoneTable {
		assert.Equal(t, testData.out.GetTimeZone(), testData.in.GetTimeZone())
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

type testDecodedLocalTimeZone struct {
	in  string
	out string
}

func TestDecodeLocalTimeZone(t *testing.T) {
	testDecodedLocalTimeZoneTable := []testDecodedLocalTimeZone{
		{
			in: decodeLocalTimeZone(nasType.LocalTimeZone{
				Iei:   nasMessage.ConfigurationUpdateCommandLocalTimeZoneType,
				Octet: uint8(0x23),
			}),
			out: "+08:00",
		},
		{
			in: decodeLocalTimeZone(nasType.LocalTimeZone{
				Iei:   nasMessage.ConfigurationUpdateCommandLocalTimeZoneType,
				Octet: uint8(0x99),
			}),
			out: "-04:45",
		},
	}

	for _, testData := range testDecodedLocalTimeZoneTable {
		assert.Equal(t, testData.out, testData.in)
	}
}

type nasConvertNetworkDaylightSavingTime struct {
	in  nasType.NetworkDaylightSavingTime
	out nasType.NetworkDaylightSavingTime
}

func TestDaylightSavingTimeToNas(t *testing.T) {
	nasConvertNetworkDaylightSavingTimeTable := []nasConvertNetworkDaylightSavingTime{
		{
			in: nasConvert.DaylightSavingTimeToNas("-05:00+1"), // EST to EDT
			out: nasType.NetworkDaylightSavingTime{
				Iei:   nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType,
				Len:   uint8(0x01),
				Octet: uint8(0x01),
			},
		},
	}

	for _, testData := range nasConvertNetworkDaylightSavingTimeTable {
		assert.Equal(t, testData.out.GetLen(), testData.in.GetLen())
		assert.Equal(t, testData.out.Getvalue(), testData.in.Getvalue())
	}
}

func TestDecodeDaylightSavingTime(t *testing.T) {
	daylightSavingTimeTable := []nasType.NetworkDaylightSavingTime{
		{
			Iei:   nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType,
			Len:   uint8(0x01),
			Octet: uint8(0x00),
		},
		{
			Iei:   nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType,
			Len:   uint8(0x01),
			Octet: uint8(0x01),
		},
		{
			Iei:   nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType,
			Len:   uint8(0x01),
			Octet: uint8(0x02),
		},
	}

	for _, testData := range daylightSavingTimeTable {
		assert.LessOrEqual(t, testData.Getvalue(), uint8(0x02))
	}
}
