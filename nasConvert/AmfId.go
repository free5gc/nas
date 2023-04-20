package nasConvert

import (
	"encoding/hex"
	"fmt"
	"log"
)

func AmfIdToNas(amfId string) (amfRegionId uint8, amfSetId uint16, amfPointer uint8) {
	var err error
	amfRegionId, amfSetId, amfPointer, err = AmfIdToNasWithError(amfId)
	if err != nil {
		log.Printf("AmfIdToNas: %+v", err)
		return 0, 0, 0
	}
	return
}

func AmfIdToNasWithError(amfId string) (amfRegionId uint8, amfSetId uint16, amfPointer uint8, err error) {
	amfIdBytes, err := hex.DecodeString(amfId)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("amfId decode failed: %w", err)
	}

	amfRegionId = amfIdBytes[0]
	amfSetId = uint16(amfIdBytes[1])<<2 + (uint16(amfIdBytes[2])&0x00c0)>>6
	amfPointer = amfIdBytes[2] & 0x3f
	return
}

func AmfIdToModels(amfRegionId uint8, amfSetId uint16, amfPointer uint8) (amfId string) {
	tmpBytes := []uint8{amfRegionId, uint8(amfSetId>>2) & 0xff, uint8(amfSetId&0x03) + amfPointer&0x3f}
	amfId = hex.EncodeToString(tmpBytes)
	return
}
