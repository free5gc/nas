package nasConvert

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
)

func ModelsToSessionAMBR(ambr *models.Ambr) (sessAmbr nasType.SessionAMBR) {
	var bitRate int64
	var bitRateBytes [2]byte

	fmt.Println(ambr)

	uplink := strings.Split(ambr.Uplink, " ")
	bitRate, _ = strconv.ParseInt(uplink[0], 10, 16)
	binary.LittleEndian.PutUint16(bitRateBytes[:], uint16(bitRate))
	sessAmbr.SetSessionAMBRForUplink(bitRateBytes)
	sessAmbr.SetUnitForSessionAMBRForUplink(strToAMBRUnit(uplink[1]))

	downlink := strings.Split(ambr.Downlink, " ")
	bitRate, _ = strconv.ParseInt(downlink[0], 10, 16)
	binary.LittleEndian.PutUint16(bitRateBytes[:], uint16(bitRate))
	sessAmbr.SetSessionAMBRForDownlink(bitRateBytes)
	sessAmbr.SetUnitForSessionAMBRForDownlink(strToAMBRUnit(downlink[1]))
	return
}

func strToAMBRUnit(unit string) uint8 {
	switch unit {
	case "bps":
		return nasMessage.SessionAMBRUnitNotUsed
	case "Kbps":
		return nasMessage.SessionAMBRUnit1Kbps
	case "Mbps":
		return nasMessage.SessionAMBRUnit1Mbps
	case "Gbps":
		return nasMessage.SessionAMBRUnit1Gbps
	case "Tbps":
		return nasMessage.SessionAMBRUnit1Tbps
	case "Pbps":
		return nasMessage.SessionAMBRUnit1Pbps
	}
	return nasMessage.SessionAMBRUnitNotUsed
}
