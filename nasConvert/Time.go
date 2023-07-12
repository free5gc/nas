package nasConvert

import (
	"strings"

	"github.com/free5gc/nas/nasType"
)

func LocalTimeZoneToNas(timezone string) (nasTimezone nasType.LocalTimeZone) {
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
	time = ((time / 10) << 4) + (time % 10)

	// Add signed number
	if timezone[0] == '-' {
		time |= 0x80
	}

	// Swap the semi-octet
	time = ((time & 0x0F) << 4) | ((time & 0xF0) >> 4)

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
