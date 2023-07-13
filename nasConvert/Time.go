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
	nasUniversalTimeAndLocalTimeZoneToNas nasType.UniversalTimeAndLocalTimeZone,
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
	} else {
		timezone += "+"
	}
	timezone += fmt.Sprintf("%2d:%2d", offset/3600, (offset%3600)/60)

	nasUniversalTimeAndLocalTimeZoneToNas.SetYear(uint8(year))
	nasUniversalTimeAndLocalTimeZoneToNas.SetMonth(uint8(month))
	nasUniversalTimeAndLocalTimeZoneToNas.SetDay(uint8(day))
	nasUniversalTimeAndLocalTimeZoneToNas.SetHour(uint8(hour))
	nasUniversalTimeAndLocalTimeZoneToNas.SetMinute(uint8(minute))
	nasUniversalTimeAndLocalTimeZoneToNas.SetSecond(uint8(second))
	nasUniversalTimeAndLocalTimeZoneToNas.SetTimeZone(uint8(parseTimeZone(timezone)))

	return
}

func LocalTimeZoneToNas(timezone string) (nasTimezone nasType.LocalTimeZone) {
	time := parseTimeZone(timezone)

	nasTimezone.SetTimeZone(uint8(time))
	return
}

func DaylightSavingTimeToNas(timezone string) (nasDaylightSavingTimeToNas nasType.NetworkDaylightSavingTime) {
	value := 0

	if strings.Contains(timezone, "+1") {
		value = 1
	}

	if strings.Contains(timezone, "+2") {
		value = 2
	}

	nasDaylightSavingTimeToNas.SetLen(1)
	nasDaylightSavingTimeToNas.Setvalue(uint8(value))
	return
}
