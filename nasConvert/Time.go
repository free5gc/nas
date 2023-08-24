package nasConvert

import (
	"fmt"
	"strings"
	"time"

	"github.com/free5gc/nas/nasType"
)

func toBinaryCodedDecimal(val int) int {
	return ((val / 10) << 4) + (val % 10)
}

// Refer to TS 23.040 - 9.1.2.3  Semi-octet representation
func toSemiOctet(val int) int {
	return ((val & 0x0F) << 4) | ((val & 0xF0) >> 4)
}

func parseTimeZoneToNas(timezone string) int {
	time := 0 // expressed in quarters of an hour

	// Parse hour
	if timezone[1] == '1' {
		time += (10 * 4)
	}
	for i := 0; i < 10; i++ {
		if int(timezone[2]) == (i + 0x30) {
			time += i * 4
		}
	}
	if timezone[len(timezone)-2:] == "+1" || timezone[len(timezone)-2:] == "+2" {
		idx := strings.LastIndex(timezone, "+")
		if idx != -1 {
			if timezone[0] == '-' {
				time -= (int(timezone[idx+1]) - 0x30) * 4
			} else {
				time += (int(timezone[idx+1]) - 0x30) * 4
			}
		}
	}

	// Parse minute
	switch timezone[4:6] {
	case "15":
		time += 1
	case "30":
		time += 2
	case "45":
		time += 3
	default:
		time += 0
	}

	// Convert decimal to binary-coded decimal
	time = toBinaryCodedDecimal(time)

	// Add signed number
	if timezone[0] == '-' {
		time |= 0x80
	}

	time = toSemiOctet(time)
	return time
}

// Get time zone string from time.Time structure
func GetTimeZone(now time.Time) string {
	timezone := ""
	_, offset := now.Zone()
	if now.IsDST() {
		// Adjust one hour to get the orignal time
		offset -= 3600
	}
	if offset < 0 {
		timezone += "-"
		offset = 0 - offset
	} else {
		timezone += "+"
	}
	timezone += fmt.Sprintf("%02d:%02d", offset/3600, (offset%3600)/60)
	if now.IsDST() {
		timezone += "+1"
	}

	return timezone
}

func EncodeUniversalTimeAndLocalTimeZoneToNas(
	universalTime time.Time,
) nasType.UniversalTimeAndLocalTimeZone {
	var nasUniversalTimeAndLocalTimeZone nasType.UniversalTimeAndLocalTimeZone

	year := toSemiOctet(toBinaryCodedDecimal(universalTime.Year() % 100))
	month := toSemiOctet(toBinaryCodedDecimal(int(universalTime.Month())))
	day := toSemiOctet(toBinaryCodedDecimal(universalTime.Day()))
	hour := toSemiOctet(toBinaryCodedDecimal(universalTime.Hour()))
	minute := toSemiOctet(toBinaryCodedDecimal(universalTime.Minute()))
	second := toSemiOctet(toBinaryCodedDecimal(universalTime.Second()))
	timezone := GetTimeZone(universalTime)

	nasUniversalTimeAndLocalTimeZone.SetYear(uint8(year))
	nasUniversalTimeAndLocalTimeZone.SetMonth(uint8(month))
	nasUniversalTimeAndLocalTimeZone.SetDay(uint8(day))
	nasUniversalTimeAndLocalTimeZone.SetHour(uint8(hour))
	nasUniversalTimeAndLocalTimeZone.SetMinute(uint8(minute))
	nasUniversalTimeAndLocalTimeZone.SetSecond(uint8(second))
	nasUniversalTimeAndLocalTimeZone.SetTimeZone(uint8(parseTimeZoneToNas(timezone)))
	return nasUniversalTimeAndLocalTimeZone
}

func EncodeLocalTimeZoneToNas(
	timezone string,
) nasType.LocalTimeZone {
	var nasLocalTimeZone nasType.LocalTimeZone

	nasLocalTimeZone.SetTimeZone(uint8(parseTimeZoneToNas(timezone)))
	return nasLocalTimeZone
}

func EncodeDaylightSavingTimeToNas(
	timezone string,
) nasType.NetworkDaylightSavingTime {
	var nasDaylightSavingTime nasType.NetworkDaylightSavingTime

	value := 0
	if strings.Contains(timezone, "+1") {
		value = 1
	}
	if strings.Contains(timezone, "+2") {
		value = 2
	}

	nasDaylightSavingTime.SetLen(1)
	nasDaylightSavingTime.Setvalue(uint8(value))
	return nasDaylightSavingTime
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

func DecodeUniversalTimeAndLocalTimeZone(
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

func DecodeLocalTimeZone(nasLocalTimeZone nasType.LocalTimeZone) string {
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

func DecodeDaylightSavingTime(nasDaylightSavingTime nasType.NetworkDaylightSavingTime) string {
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
