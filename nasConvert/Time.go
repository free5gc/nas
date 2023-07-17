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

func toSemiOctet(val int) int {
	// Refer to TS 23.040 - 9.1.2.3  Semi-octet representation
	return ((val & 0x0F) << 4) | ((val & 0xF0) >> 4)
}

func parseTimeZone(timezone string) int {
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

func UniversalTimeAndLocalTimeZoneToNas(universalTime time.Time) (
	nasUniversalTimeAndLocalTimeZone nasType.UniversalTimeAndLocalTimeZone,
) {
	year := toSemiOctet(toBinaryCodedDecimal(universalTime.Year() % 100))
	month := toSemiOctet(toBinaryCodedDecimal(int(universalTime.Month())))
	day := toSemiOctet(toBinaryCodedDecimal(universalTime.Day()))
	hour := toSemiOctet(toBinaryCodedDecimal(universalTime.Hour()))
	minute := toSemiOctet(toBinaryCodedDecimal(universalTime.Minute()))
	second := toSemiOctet(toBinaryCodedDecimal(universalTime.Second()))
	timezone := ""

	_, offset := universalTime.Zone()
	if offset < 0 {
		timezone += "-"
		offset = 0 - offset
	} else {
		timezone += "+"
	}
	timezone += fmt.Sprintf("%02d:%02d", offset/3600, (offset%3600)/60)
	if universalTime.IsDST() {
		timezone += "+1"
	}

	nasUniversalTimeAndLocalTimeZone.SetYear(uint8(year))
	nasUniversalTimeAndLocalTimeZone.SetMonth(uint8(month))
	nasUniversalTimeAndLocalTimeZone.SetDay(uint8(day))
	nasUniversalTimeAndLocalTimeZone.SetHour(uint8(hour))
	nasUniversalTimeAndLocalTimeZone.SetMinute(uint8(minute))
	nasUniversalTimeAndLocalTimeZone.SetSecond(uint8(second))
	nasUniversalTimeAndLocalTimeZone.SetTimeZone(uint8(parseTimeZone(timezone)))

	return
}

func LocalTimeZoneToNas(timezone string) (nasLocalTimeZone nasType.LocalTimeZone) {
	time := parseTimeZone(timezone)

	nasLocalTimeZone.SetTimeZone(uint8(time))
	return
}

func DaylightSavingTimeToNas(timezone string) (nasDaylightSavingTime nasType.NetworkDaylightSavingTime) {
	value := 0

	if strings.Contains(timezone, "+1") {
		value = 1
	}

	if strings.Contains(timezone, "+2") {
		value = 2
	}

	nasDaylightSavingTime.SetLen(1)
	nasDaylightSavingTime.Setvalue(uint8(value))
	return
}
