package nasType

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

// MobileIdentity5GS 9.11.3.4
// MobileIdentity5GSContents Row, sBit, len = [0, 0], 8 , INF
type MobileIdentity5GS struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

// const
const (
	noIdentity uint8 = iota
	suci
	fiveGGUTI
	imei
	fiveGSTMSI
	imeisv
)

// const
const (
	high5BitMask uint8 = 0x07
	low4BitMask  uint8 = 0xf0
	high4BitMask uint8 = 0x0f
	Bit4         uint8 = 0x08
)

func NewMobileIdentity5GS(iei uint8) (mobileIdentity5GS *MobileIdentity5GS) {
	mobileIdentity5GS = &MobileIdentity5GS{}
	mobileIdentity5GS.SetIei(iei)
	return mobileIdentity5GS
}

// MobileIdentity5GS 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *MobileIdentity5GS) GetIei() (iei uint8) {
	return a.Iei
}

// MobileIdentity5GS 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *MobileIdentity5GS) SetIei(iei uint8) {
	a.Iei = iei
}

// MobileIdentity5GS 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *MobileIdentity5GS) GetLen() (len uint16) {
	return a.Len
}

// MobileIdentity5GS 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *MobileIdentity5GS) SetLen(len uint16) {
	a.Len = len
	a.Buffer = make([]uint8, a.Len)
}

// MobileIdentity5GS 9.11.3.4
// MobileIdentity5GSContents Row, sBit, len = [0, 0], 8 , INF
func (a *MobileIdentity5GS) GetMobileIdentity5GSContents() (mobileIdentity5GSContents []uint8) {
	mobileIdentity5GSContents = make([]uint8, len(a.Buffer))
	copy(mobileIdentity5GSContents, a.Buffer)
	return mobileIdentity5GSContents
}

// MobileIdentity5GS 9.11.3.4
// MobileIdentity5GSContents Row, sBit, len = [0, 0], 8 , INF
func (a *MobileIdentity5GS) SetMobileIdentity5GSContents(mobileIdentity5GSContents []uint8) {
	copy(a.Buffer, mobileIdentity5GSContents)
}

// TS 24.501 9.11.3.3 Table 9.11.3.3.1 identity type information element
// All other values are unused and shall be interpreted
// as "SUCI", if received by the UE
func (a *MobileIdentity5GS) GetTypeOfIdentity() (string, error) {
	if len(a.Buffer) == 0 {
		return "", errors.New("buffer is empty")
	}
	idType := a.Buffer[0] & high5BitMask
	switch idType {
	case noIdentity:
		return "", errors.New("no identity")
	case suci:
		return "SUCI", nil
	case fiveGGUTI:
		return "5G-GUTI", nil
	case imei:
		return "IMEI", nil
	case fiveGSTMSI:
		return "5G-S-TMSI", nil
	case imeisv:
		return "IMEISV", nil
	default:
		return "SUCI", nil
	}
}

// GetMobileIdentity
func (a *MobileIdentity5GS) GetMobileIdentity() (string, string, error) {
	if len(a.Buffer) == 0 {
		return "", "", fmt.Errorf("empty buffer")
	}

	idType, err := a.GetTypeOfIdentity()
	if err != nil {
		return "", idType, err
	}

	switch idType {
	case "SUCI":
		suci, err := a.GetSUCI()
		if err != nil {
			return "", idType, err
		}
		return suci, idType, nil
	case "5G-GUTI":
		guti, err := a.Get5GGUTI()
		if err != nil {
			return "", idType, err
		}
		return guti, idType, nil
	case "IMEI":
		imei, err := a.GetIMEI()
		if err != nil {
			return "", idType, err
		}
		return imei, idType, nil
	case "5G-S-TMSI":
		tmsi, err := a.Get5GTMSI()
		if err != nil {
			return "", idType, err
		}
		return tmsi, idType, nil
	case "IMEISV":
		imeisv, err := a.GetIMEISV()
		if err != nil {
			return "", idType, err
		}
		return imeisv, idType, nil
	default:
		// Default to SUCI
		suci, err := a.GetSUCI()
		if err != nil {
			return "", "SUCI", err
		}
		return suci, "SUCI", nil
	}
}

// GetSUCI
func (a *MobileIdentity5GS) GetSUCI() (string, error) {
	if len(a.Buffer) < 8 {
		return "", fmt.Errorf("invalid SUCI length")
	}
	idType, err := a.GetTypeOfIdentity()
	if err != nil {
		return "", err
	}

	if idType == "SUCI" {
		var schemeOutput string

		// Encode buf to SUCI in supi format "IMSI"
		supiFormat := (a.Buffer[0] & low4BitMask) >> 4
		if supiFormat == suci {
			return naiToString(a.Buffer), nil
		}

		mcc, err := a.GetMCC()
		if err != nil {
			return "", err
		}
		mnc, err := a.GetMNC()
		if err != nil {
			return "", err
		}

		var routingIndBytes []byte
		routingIndBytes = append(routingIndBytes, bits.RotateLeft8(a.Buffer[4], 4))
		routingIndBytes = append(routingIndBytes, bits.RotateLeft8(a.Buffer[5], 4))
		routingInd := hex.EncodeToString(routingIndBytes)

		if idx := strings.Index(routingInd, "f"); idx != -1 {
			routingInd = routingInd[0:idx]
		}

		// Protection Scheme
		protectionScheme := fmt.Sprintf("%x", a.Buffer[6]) // convert byte to hex string without leading 0s

		// Home Network Public Key Indentifier
		homeNetworkPublicKeyIdentifier := fmt.Sprintf("%d", a.Buffer[7])

		// Scheme output
		// TS 24.501 9.11.3.4
		if protectionScheme == "0" {
			if len(a.Buffer) < 9 {
				return "", fmt.Errorf("invalid SUCI length")
			}
			// MSIN
			var msinBytes []byte
			for i := 8; i < len(a.Buffer); i++ {
				msinBytes = append(msinBytes, bits.RotateLeft8(a.Buffer[i], 4))
			}
			schemeOutput = hex.EncodeToString(msinBytes)
			if schemeOutput[len(schemeOutput)-1] == 'f' {
				schemeOutput = schemeOutput[:len(schemeOutput)-1]
			}
		} else {
			schemeOutput = hex.EncodeToString(a.Buffer[8:])
		}

		// "suci-0-208-93-0-0-0-00007487"
		suci := strings.Join([]string{
			"suci", "0", mcc, mnc, routingInd, protectionScheme, homeNetworkPublicKeyIdentifier,
			schemeOutput,
		}, "-")

		return suci, nil
	}
	return "", fmt.Errorf("identity type is not SUCI")
}

// GetPlmnID
func (a *MobileIdentity5GS) GetPlmnID() (string, error) {
	mcc, err := a.GetMCC()
	if err != nil {
		return "", err
	}
	mnc, err := a.GetMNC()
	if err != nil {
		return "", err
	}
	plmnId := mcc + mnc
	return plmnId, nil
}

// GetMCC
func (a *MobileIdentity5GS) GetMCC() (string, error) {
	if len(a.Buffer) < 3 {
		return "", fmt.Errorf("invalid MCC length")
	}
	mccDigit3 := (a.Buffer[2] & high4BitMask)
	tmpBytes := []byte{bits.RotateLeft8(a.Buffer[1], 4), (mccDigit3 << 4)}
	mcc := hex.EncodeToString(tmpBytes)
	mcc = mcc[:3] // remove rear 0
	return mcc, nil
}

// GetMNC
func (a *MobileIdentity5GS) GetMNC() (string, error) {
	if len(a.Buffer) < 4 {
		return "", fmt.Errorf("invalid MNC length")
	}
	mncDigit3 := (a.Buffer[2] & low4BitMask) >> 4
	tmpBytes := []byte{bits.RotateLeft8(a.Buffer[3], 4), mncDigit3 << 4}
	mnc := hex.EncodeToString(tmpBytes)
	if mnc[2] == 'f' {
		mnc = mnc[:2] // mnc is 2 digit -> remove 'f'
	} else {
		mnc = mnc[:3] // mnc is 3 digit -> remove rear 0
	}
	return mnc, nil
}

// Get5GGUTI
func (a *MobileIdentity5GS) Get5GGUTI() (string, error) {
	mcc, err := a.GetMCC()
	if err != nil {
		return "", err
	}
	mnc, err := a.GetMNC()
	if err != nil {
		return "", err
	}
	amfID, err := a.GetAmfID()
	if err != nil {
		return "", err
	}
	tmsi, err := a.Get5GTMSI()
	if err != nil {
		return "", err
	}
	return mcc + mnc + amfID + tmsi, nil
}

// GetAmfID
func (a *MobileIdentity5GS) GetAmfID() (string, error) {
	if len(a.Buffer) < 7 {
		return "", fmt.Errorf("invalid AmfID length")
	}
	return hex.EncodeToString(a.Buffer[4:7]), nil
}

// GetAmfRegionID
func (a *MobileIdentity5GS) GetAmfRegionID() (string, error) {
	if len(a.Buffer) < 5 {
		return "", fmt.Errorf("invalid AmfRegionID length")
	}
	return hex.EncodeToString(a.Buffer[4:5]), nil
}

// GetAmfSetID
func (a *MobileIdentity5GS) GetAmfSetID() (string, error) {
	var amfSetStartPoint int
	idType, err := a.GetTypeOfIdentity()
	if err != nil {
		return "", err
	}

	switch idType {
	case "5G-GUTI":
		amfSetStartPoint = 5
	case "5G-S-TMSI":
		amfSetStartPoint = 1
	default:
		return "", fmt.Errorf("wrong identity type for AmfSetID")
	}

	if len(a.Buffer) < amfSetStartPoint+2 {
		return "", fmt.Errorf("invalid AmfSetID length")
	}

	amfSetID := (uint16(a.Buffer[amfSetStartPoint])<<2 + uint16((a.Buffer[amfSetStartPoint+1])&GetBitMask(8, 2))>>6)
	amfSetID_string := strconv.FormatUint(uint64(amfSetID), 10)
	return amfSetID_string, nil
}

// GetAmfPointer
func (a *MobileIdentity5GS) GetAmfPointer() (string, error) {
	var amfPointerStartPoint int
	idType, err := a.GetTypeOfIdentity()
	if err != nil {
		return "", err
	}

	switch idType {
	case "5G-GUTI":
		amfPointerStartPoint = 6
	case "5G-S-TMSI":
		amfPointerStartPoint = 2
	default:
		return "", fmt.Errorf("wrong identity type for AmfPointer")
	}

	if len(a.Buffer) <= amfPointerStartPoint {
		return "", fmt.Errorf("invalid AmfPointer length")
	}

	AMFPointer := (a.Buffer[amfPointerStartPoint]) & GetBitMask(6, 0)
	AMFPointer_string := strconv.FormatUint(uint64(AMFPointer), 10)
	return AMFPointer_string, nil
}

// Get5GTMSI
func (a *MobileIdentity5GS) Get5GTMSI() (string, error) {
	idType, err := a.GetTypeOfIdentity()
	if err != nil {
		return "", err
	}

	switch idType {
	case "5G-GUTI":
		if len(a.Buffer) <= 7 {
			return "", fmt.Errorf("invalid 5G-GUTI length for TMSI")
		}
		tmsi5G_string := hex.EncodeToString(a.Buffer[7:])
		return tmsi5G_string, nil
	case "5G-S-TMSI":
		if len(a.Buffer) < 7 {
			return "", fmt.Errorf("invalid 5G-S-TMSI length for TMSI")
		}
		tmsi5G := a.Buffer[3:7]
		tmsi5G_string := hex.EncodeToString(tmsi5G[0:])

		return tmsi5G_string, nil
	default:
		return "", fmt.Errorf("wrong identity type for 5G-TMSI")
	}
}

// GetIMEI
func (a *MobileIdentity5GS) GetIMEI() (string, error) {
	idType, err := a.GetTypeOfIdentity()
	if err != nil {
		return "", err
	}
	switch idType {
	case "IMEI":
		return "imei-" + peiToString(a.Buffer), nil
	default:
		return "", fmt.Errorf("identity type is not IMEI")
	}
}

// GetIMEISV
func (a *MobileIdentity5GS) GetIMEISV() (string, error) {
	idType, err := a.GetTypeOfIdentity()
	if err != nil {
		return "", err
	}
	switch idType {
	case "IMEISV":
		return "imeisv-" + peiToString(a.Buffer), nil
	default:
		return "", fmt.Errorf("identity type is not IMEISV")
	}
}

func (a *MobileIdentity5GS) Get5GSTMSI() (tMSI5GS string, mobileIdType string, err error) {
	if len(a.Buffer) < 7 {
		return "", "", fmt.Errorf("invalid 5G-S-TMSI length")
	}
	partOfAmfId := hex.EncodeToString(a.Buffer[1:3])
	tmsi5g, err := a.Get5GTMSI()
	if err != nil {
		return "", "", err
	}
	tMSI5GS = partOfAmfId + tmsi5g
	return tMSI5GS, "5G-S-TMSI", nil
}

func naiToString(buf []byte) string {
	prefix := "nai"
	naiBytes := buf[1:]
	naiStr := hex.EncodeToString(naiBytes)
	nai := strings.Join([]string{prefix, "1", naiStr}, "-")
	return nai
}

func peiToString(buf []byte) string {
	oddIndication := (buf[0] & Bit4) >> 3
	digit1 := (buf[0] & low4BitMask)
	tmpBytes := []byte{digit1}

	for _, octet := range buf[1:] {
		digitP := octet & high4BitMask
		digitP1 := octet & low4BitMask

		tmpBytes[len(tmpBytes)-1] += digitP
		tmpBytes = append(tmpBytes, digitP1)
	}

	digitStr := hex.EncodeToString(tmpBytes)
	digitStr = digitStr[:len(digitStr)-1] // remove the last digit

	if oddIndication == 0 { // even digits
		digitStr = digitStr[:len(digitStr)-1] // remove the last digit
	}
	return digitStr
}
