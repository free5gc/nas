package nasConvert

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
	"unicode"

	"github.com/free5gc/nas/logger"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
)

func GetTypeOfIdentity(buf byte) uint8 {
	return buf & 0x07
}

// TS 24.501 9.11.3.4
// suci(imsi) =
// "suci-0-${mcc}-${mnc}-${routingIndentifier}-${protectionScheme}-${homeNetworkPublicKeyIdentifier}-${schemeOutput}"
// suci(nai) = "nai-${naiString}"
func SuciToString(buf []byte) (suci string, plmnId string) {
	var err error
	suci, plmnId, err = SuciToStringWithError(buf)
	if err != nil {
		logger.ConvertLog.Warnf("SuciToString: %+v", err)
		return "", ""
	}
	return
}

func SuciToStringWithError(buf []byte) (suci string, plmnId string, err error) {
	var mcc, mnc, routingInd, protectionScheme, homeNetworkPublicKeyIdentifier, schemeOutput string

	if len(buf) < 1 {
		return "", "", errors.New("too short SUCI")
	}

	supiFormat := (buf[0] & 0xf0) >> 4
	if supiFormat == nasMessage.SupiFormatNai {
		suci, err = naiToString(buf)
		return suci, "", err
	}

	if len(buf) < 9 {
		return "", "", errors.New("too short SUCI")
	}

	// Encode buf to SUCI in supi format "IMSI"

	// Plmn(MCC + MNC)
	mccDigit3 := (buf[2] & 0x0f)
	tmpBytes := []byte{bits.RotateLeft8(buf[1], 4), (mccDigit3 << 4)}
	mcc = hex.EncodeToString(tmpBytes)
	mcc = mcc[:3] // remove rear 0

	mncDigit3 := (buf[2] & 0xf0) >> 4
	tmpBytes = []byte{bits.RotateLeft8(buf[3], 4), mncDigit3 << 4}
	mnc = hex.EncodeToString(tmpBytes)
	if mnc[2] == 'f' {
		mnc = mnc[:2] // mnc is 2 digit -> remove 'f'
	} else {
		mnc = mnc[:3] // mnc is 3 digit -> remove rear 0
	}
	plmnId = mcc + mnc

	// Routing Indicator
	var routingIndBytes []byte
	routingIndBytes = append(routingIndBytes, bits.RotateLeft8(buf[4], 4))
	routingIndBytes = append(routingIndBytes, bits.RotateLeft8(buf[5], 4))
	routingInd = hex.EncodeToString(routingIndBytes)

	if idx := strings.Index(routingInd, "f"); idx != -1 {
		routingInd = routingInd[0:idx]
	}

	// Protection Scheme
	protectionScheme = fmt.Sprintf("%x", buf[6]) // convert byte to hex string without leading 0s

	// Home Network Public Key Indentifier
	homeNetworkPublicKeyIdentifier = fmt.Sprintf("%d", buf[7])

	// Scheme output
	if protectionScheme == strconv.Itoa(nasMessage.ProtectionSchemeNullScheme) {
		// MSIN
		var msinBytes []byte
		for i := 8; i < len(buf); i++ {
			msinBytes = append(msinBytes, bits.RotateLeft8(buf[i], 4))
		}
		schemeOutput = hex.EncodeToString(msinBytes)
		if schemeOutput[len(schemeOutput)-1] == 'f' {
			schemeOutput = schemeOutput[:len(schemeOutput)-1]
		}
	} else {
		schemeOutput = hex.EncodeToString(buf[8:])
	}

	suci = strings.Join([]string{
		"suci", "0", mcc, mnc, routingInd, protectionScheme, homeNetworkPublicKeyIdentifier,
		schemeOutput,
	}, "-")
	return suci, plmnId, nil
}

func NaiToString(buf []byte) (nai string) {
	var err error
	nai, err = naiToString(buf)
	if err != nil {
		logger.ConvertLog.Warnf("NaiToString: %+v", err)
		return ""
	}
	return
}

func naiToString(buf []byte) (nai string, err error) {
	if len(buf) < 2 {
		return "", errors.New("too short NAI")
	}
	prefix := "nai"
	naiBytes := buf[1:]
	naiStr := hex.EncodeToString(naiBytes)
	nai = strings.Join([]string{prefix, "1", naiStr}, "-")
	return
}

// nasType: TS 24.501 9.11.3.4
func GutiToString(buf []byte) (guami models.Guami, guti string) {
	var err error
	guami, guti, err = GutiToStringWithError(buf)
	if err != nil {
		logger.ConvertLog.Warnf("GutiToString: %+v", err)
		return models.Guami{}, ""
	}
	return
}

func GutiToStringWithError(buf []byte) (guami models.Guami, guti string, err error) {
	if len(buf) != 11 {
		return models.Guami{}, "", errors.New("invalid GUTI length")
	}
	plmnID := PlmnIDToString(buf[1:4])
	amfID := hex.EncodeToString(buf[4:7])
	tmsi5G := hex.EncodeToString(buf[7:])

	guami.PlmnId = new(models.PlmnId)
	guami.PlmnId.Mcc = plmnID[:3]
	guami.PlmnId.Mnc = plmnID[3:]
	guami.AmfId = amfID
	guti = plmnID + amfID + tmsi5G
	return
}

func GutiToNas(guti string) nasType.GUTI5G {
	gutiNas, err := GutiToNasWithError(guti)
	if err != nil {
		logger.ConvertLog.Warnf("GutiToNas: %+v", err)
		return nasType.GUTI5G{Len: 11}
	}
	return gutiNas
}

func GutiToNasWithError(guti string) (nasType.GUTI5G, error) {
	var gutiNas nasType.GUTI5G

	if len(guti) != 19 && len(guti) != 20 {
		return nasType.GUTI5G{}, errors.New("invalid GUTI length")
	}

	gutiNas.SetLen(11)
	gutiNas.SetSpare(0)
	gutiNas.SetSpare2(15)
	gutiNas.SetTypeOfIdentity(nasMessage.MobileIdentity5GSType5gGuti)

	var mcc1, mcc2, mcc3, mnc1, mnc2, mnc3 int
	if mcc1Tmp, err := strconv.Atoi(string(guti[0])); err != nil {
		return nasType.GUTI5G{}, fmt.Errorf("atoi mcc1 error: %w", err)
	} else {
		mcc1 = mcc1Tmp
	}
	if mcc2Tmp, err := strconv.Atoi(string(guti[1])); err != nil {
		return nasType.GUTI5G{}, fmt.Errorf("atoi mcc2 error: %w", err)
	} else {
		mcc2 = mcc2Tmp
	}
	if mcc3Tmp, err := strconv.Atoi(string(guti[2])); err != nil {
		return nasType.GUTI5G{}, fmt.Errorf("atoi mcc3 error: %w", err)
	} else {
		mcc3 = mcc3Tmp
	}
	if mnc1Tmp, err := strconv.Atoi(string(guti[3])); err != nil {
		return nasType.GUTI5G{}, fmt.Errorf("atoi mnc1 error: %w", err)
	} else {
		mnc1 = mnc1Tmp
	}
	if mnc2Tmp, err := strconv.Atoi(string(guti[4])); err != nil {
		return nasType.GUTI5G{}, fmt.Errorf("atoi mnc2 error: %w", err)
	} else {
		mnc2 = mnc2Tmp
	}
	mnc3 = 0x0f
	amfId := ""
	tmsi := ""
	if len(guti) == 20 {
		if mnc3Tmp, err := strconv.Atoi(string(guti[5])); err != nil {
			return nasType.GUTI5G{}, fmt.Errorf("atoi guti error: %w", err)
		} else {
			mnc3 = mnc3Tmp
		}
		amfId = guti[6:12]
		tmsi = guti[12:]
	} else {
		amfId = guti[5:11]
		tmsi = guti[11:]
	}
	gutiNas.SetMCCDigit1(uint8(mcc1))
	gutiNas.SetMCCDigit2(uint8(mcc2))
	gutiNas.SetMCCDigit3(uint8(mcc3))
	gutiNas.SetMNCDigit1(uint8(mnc1))
	gutiNas.SetMNCDigit2(uint8(mnc2))
	gutiNas.SetMNCDigit3(uint8(mnc3))

	if amfRegionId, amfSetId, amfPointer, err := AmfIdToNasWithError(amfId); err != nil {
		return nasType.GUTI5G{}, fmt.Errorf("decode AMF ID failed: %w", err)
	} else {
		gutiNas.SetAMFRegionID(amfRegionId)
		gutiNas.SetAMFSetID(amfSetId)
		gutiNas.SetAMFPointer(amfPointer)
	}
	if tmsiBytes, err := hex.DecodeString(tmsi); err != nil {
		return nasType.GUTI5G{}, fmt.Errorf("decode TMSI failed: %w", err)
	} else {
		copy(gutiNas.Octet[7:11], tmsiBytes[:])
	}
	return gutiNas, nil
}

// PEI: ^(imei-[0-9]{15}|imeisv-[0-9]{16}|.+)$
func PeiToString(buf []byte) string {
	pei, err := PeiToStringWithError(buf)
	if err != nil {
		logger.ConvertLog.Warnf("PeiToString: %+v", err)
		return ""
	}
	return pei
}

func PeiToStringWithError(buf []byte) (string, error) {
	var prefix string

	if len(buf) < 1 {
		return "", errors.New("too short PEI")
	}

	typeOfIdentity := buf[0] & 0x07
	if typeOfIdentity == 0x03 {
		prefix = "imei-"
	} else {
		prefix = "imeisv-"
	}

	oddIndication := (buf[0] & 0x08) >> 3

	digit1 := (buf[0] & 0xf0)

	tmpBytes := []byte{digit1}

	for _, octet := range buf[1:] {
		digitP := octet & 0x0f
		digitP1 := octet & 0xf0

		tmpBytes[len(tmpBytes)-1] += digitP
		tmpBytes = append(tmpBytes, digitP1)
	}

	digitStr := hex.EncodeToString(tmpBytes)
	digitStr = digitStr[:len(digitStr)-1] // remove the last digit

	if oddIndication == 0 { // even digits
		digitStr = digitStr[:len(digitStr)-1] // remove the last digit
	}

	if prefix == "imei-" {
		// Validate IMEI before returning
		if len(digitStr) != 15 {
			return "", fmt.Errorf("invalid IMEI length: expected 15 digits, got %d", len(digitStr))
		}
		valid, err := ValidateIMEI(digitStr)
		if err != nil {
			return "", fmt.Errorf("IMEI validation error: %w", err)
		}
		if !valid {
			return "", fmt.Errorf("invalid IMEI checksum")
		}
	} else {
		if len(digitStr) != 16 {
			return "", fmt.Errorf("invalid IMEISV length: expected 16 digits, got %d", len(digitStr))
		}
	}

	return prefix + digitStr, nil
}

func validateIMEI(imei string) (bool, error) {
	// Remove any non-digit characters
	cleanIMEI := strings.ReplaceAll(imei, "-", "")
	cleanIMEI = strings.ReplaceAll(cleanIMEI, " ", "")

	// Check if all characters are digits
	for _, char := range cleanIMEI {
		if !unicode.IsDigit(char) {
			return false, fmt.Errorf("IMEI contains non-digit character: %c", char)
		}
	}

	// Luhn algorithm validation
	sum := 0
	for i := len(cleanIMEI) - 1; i >= 0; i-- {
		digit := int(cleanIMEI[i] - '0')

		if (len(cleanIMEI)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}
		sum += digit
	}

	return sum%10 == 0, nil
}
