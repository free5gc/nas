package nasType_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
)

func TestNasTypeNewMobileIdentity5GS(t *testing.T) {
	a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
	assert.NotNil(t, a)
}

var nasTypeMobileIdentity5GSRegistrationRequestAdditionalGUTITable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestAdditionalGUTIType, nasMessage.RegistrationRequestAdditionalGUTIType},
}

func TestNasTypeMobileIdentity5GSGetSetIei(t *testing.T) {
	a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
	for _, table := range nasTypeMobileIdentity5GSRegistrationRequestAdditionalGUTITable {
		a.SetIei(table.in)
		assert.Equal(t, table.out, a.GetIei())
	}
}

var nasTypeMobileIdentity5GSLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeMobileIdentity5GSGetSetLen(t *testing.T) {
	a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
	for _, table := range nasTypeMobileIdentity5GSLenTable {
		a.SetLen(table.in)
		assert.Equal(t, table.out, a.GetLen())
	}
}

type nasTypeMobileIdentity5GSMobileIdentity5GSContentsData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeMobileIdentity5GSMobileIdentity5GSContentsTable = []nasTypeMobileIdentity5GSMobileIdentity5GSContentsData{
	{2, []uint8{0xff, 0xff}, []uint8{0xff, 0xff}},
}

func TestNasTypeMobileIdentity5GSGetSetMobileIdentity5GSContents(t *testing.T) {
	a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
	for _, table := range nasTypeMobileIdentity5GSMobileIdentity5GSContentsTable {
		a.SetLen(table.inLen)
		a.SetMobileIdentity5GSContents(table.in)
		assert.Equal(t, table.out, a.GetMobileIdentity5GSContents())
	}
}

type testMobileIdentity5GSDataTemplate struct {
	inIei                        uint8
	inLen                        uint16
	inMobileIdentity5GSContents  []uint8
	outIei                       uint8
	outLen                       uint16
	outMobileIdentity5GSContents []uint8
}

var testMobileIdentity5GSTestTable = []testMobileIdentity5GSDataTemplate{
	{
		nasMessage.RegistrationRequestAdditionalGUTIType, 2,
		[]uint8{0xff, 0xff},
		nasMessage.RegistrationRequestAdditionalGUTIType, 2,
		[]uint8{0xff, 0xff},
	},
}

func TestNasTypeMobileIdentity5GS(t *testing.T) {
	for i, table := range testMobileIdentity5GSTestTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)

		a.SetIei(table.inIei)
		a.SetLen(table.inLen)
		a.SetMobileIdentity5GSContents(table.inMobileIdentity5GSContents)

		assert.Equalf(t, table.outIei, a.Iei, "in(%v): out %v, actual %x", table.inIei, table.outIei, a.Iei)
		assert.Equalf(t, table.outLen, a.Len, "in(%v): out %v, actual %x", table.inLen, table.outLen, a.Len)
		assert.Equalf(t, table.outMobileIdentity5GSContents, a.GetMobileIdentity5GSContents(), "in(%v): out %v, actual %x", table.inMobileIdentity5GSContents, table.outMobileIdentity5GSContents, a.GetMobileIdentity5GSContents())
	}
}

type GetTypeOfIdentityTemplate struct {
	inBuffer []uint8
	outType  string
	outError error
}

var GetTypeOfIdentityTable = []GetTypeOfIdentityTemplate{
	{[]uint8{0x00, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "", errors.New("no identity")},
	{[]uint8{0x01, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "SUCI", nil},
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "5G-GUTI", nil},
	{[]uint8{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xf0}, "IMEI", nil},
	{[]uint8{0xf4, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "5G-S-TMSI", nil},
	{[]uint8{0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xf0}, "IMEISV", nil},
	{[]uint8{0x09, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "SUCI", nil},
}

func TestNasTypeGetTypeOfIdentity(t *testing.T) {
	for i, table := range GetTypeOfIdentityTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		idtype, err := a.GetTypeOfIdentity()

		assert.Equalf(t, table.outType, idtype, "out %v, actual %x", table.outType, idtype)
		assert.Equalf(t, table.outError, err, "out %v, actual %x", table.outError, err)
	}
}

type GetSUCITemplate struct {
	inBuffer []uint8
	out      string
}

var GetSUCITable = []GetSUCITemplate{
	{
		[]uint8{0x01, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0},
		"suci-0-310-310--0-0-140000120",
	},
}

func TestNasTypeGetSUCI(t *testing.T) {
	for i, table := range GetSUCITable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetSUCI(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetSUCI())
	}
}

type GetPlmnIDTemplate struct {
	inBuffer []uint8
	out      string
}

var GetPlmnIDTable = []GetPlmnIDTemplate{
	{[]uint8{0x01, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "310310"},
	{[]uint8{0x01, 0x13, 0xf0, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "31031"},
}

func TestNasTypeGetPlmnID(t *testing.T) {
	for i, table := range GetPlmnIDTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetPlmnID(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetPlmnID())
	}
}

type GetMCCTemplate struct {
	inBuffer []uint8
	out      string
}

var GetMCCTable = []GetMCCTemplate{
	{[]uint8{0x01, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "310"},
}

func TestNasTypeGetMCC(t *testing.T) {
	for i, table := range GetMCCTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetMCC(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetMCC())
	}
}

type GetMNCTemplate struct {
	inBuffer []uint8
	out      string
}

var GetMNCTable = []GetMNCTemplate{
	{[]uint8{0x01, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "310"},
}

func TestNasTypeGetMNC(t *testing.T) {
	for i, table := range GetMNCTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetMNC(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetMNC())
	}
}

type Get5GGUTITemplate struct {
	inBuffer []uint8
	out      string
}

var Get5GGUTITable = []Get5GGUTITemplate{
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "310310cafe0000000001"},
}

func TestNasTypeGet5GGUTI(t *testing.T) {
	for i, table := range Get5GGUTITable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.Get5GGUTI(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.Get5GGUTI())
	}
}

type GetAmfIDTemplate struct {
	inBuffer []uint8
	out      string
}

var GetAmfIDTable = []GetAmfIDTemplate{
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "cafe00"},
}

func TestNasTypeGetAmfID(t *testing.T) {
	for i, table := range GetAmfIDTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetAmfID(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetAmfID())
	}
}

type GetAmfRegionIDTemplate struct {
	inBuffer []uint8
	out      string
}

var GetAmfRegionIDTable = []GetAmfRegionIDTemplate{
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "ca"},
}

func TestNasTypeGetAmfRegionID(t *testing.T) {
	for i, table := range GetAmfRegionIDTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetAmfRegionID(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetAmfRegionID())
	}
}

type GetAmfSetIDTemplate struct {
	inBuffer []uint8
	out      string
}

var GetAmfSetIDTable = []GetAmfSetIDTemplate{
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "1016"},
	{[]uint8{0xf4, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "1016"},
}

func TestNasTypeGetAmfSetID(t *testing.T) {
	for i, table := range GetAmfSetIDTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetAmfSetID(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetAmfSetID())
	}
}

type GetAmfPointerTemplate struct {
	inBuffer []uint8
	out      string
}

var GetAmfPointerTable = []GetAmfPointerTemplate{
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "0"},
	{[]uint8{0xf4, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "0"},
}

func TestNasTypeAmfPointerID(t *testing.T) {
	for i, table := range GetAmfPointerTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetAmfPointer(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetAmfPointer())
	}
}

type Get5GTMSITemplate struct {
	inBuffer []uint8
	out      string
}

var Get5GTMSITable = []Get5GTMSITemplate{
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "00000001"},
	{[]uint8{0xf4, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "00000001"},
}

func TestNasTypeGet5GTMSI(t *testing.T) {
	for i, table := range Get5GTMSITable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.Get5GTMSI(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.Get5GTMSI())
	}
}

type GetIMEITemplate struct {
	inBuffer []uint8
	out      string
}

var GetIMEITable = []GetIMEITemplate{
	{[]uint8{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xf0}, "imei-0000000000000200"},
}

func TestNasTypeGetIMEI(t *testing.T) {
	for i, table := range GetIMEITable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetIMEI(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetIMEI())
	}
}

type GetIMEISVTemplate struct {
	inBuffer []uint8
	out      string
}

var GetIMEISVTable = []GetIMEISVTemplate{
	{[]uint8{0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xf0}, "imeisv-0000000000000200"},
}

func TestNasTypeGetIMEISV(t *testing.T) {
	for i, table := range GetIMEISVTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		assert.Equalf(t, table.out, a.GetIMEISV(), "in(%v): out %v, actual %x",
			table.inBuffer, table.out, a.GetIMEISV())
	}
}

type Get5GSTMSITemplate struct {
	inBuffer []uint8
	out      string
	outType  string
	outerr   error
}

var Get5GSTMSITable = []Get5GSTMSITemplate{
	{[]uint8{0xf4, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "fe0000000001", "5G-S-TMSI", nil},
}

func TestNasTypeGet5GSTMSI(t *testing.T) {
	for i, table := range Get5GSTMSITable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		TMSI5GS, idtype, err := a.Get5GSTMSI()

		assert.Equalf(t, table.out, TMSI5GS, "in(%v): out %v, actual %x",
			table.inBuffer, table.out, TMSI5GS)
		assert.Equalf(t, table.outType, idtype, "in(%v): out %v, actual %x",
			table.inBuffer, table.outType, idtype)
		assert.Equalf(t, table.outerr, err, "in(%v): out %v, actual %x",
			table.inBuffer, table.outerr, err)
	}
}

type GetMobileIdentityTemplate struct {
	inBuffer []uint8
	outID    string
	outType  string
	outErr   string // Expected error message substring, empty if no error expected
}

var GetMobileIdentityTable = []GetMobileIdentityTemplate{
	// SUCI Cases
	{[]uint8{0x01, 0x13, 0x00, 0x13, 0x0f, 0xff, 0x00, 0x00, 0x41, 0x00, 0x00, 0x21, 0xf0}, "suci-0-310-310--0-0-140000120", "SUCI", ""},
	{[]uint8{0x01, 0x02, 0x03}, "", "SUCI", "invalid SUCI length"}, // Malformed SUCI

	// 5G-GUTI Cases
	{[]uint8{0xf2, 0x13, 0x00, 0x13, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "310310cafe0000000001", "5G-GUTI", ""},
	{[]uint8{0xf2, 0x01, 0x02}, "", "5G-GUTI", "invalid 5G-GUTI length"}, // Malformed GUTI

	// 5G-S-TMSI Cases
	{[]uint8{0xf4, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01}, "00000001", "5G-S-TMSI", ""},

	// IMEI Cases
	{[]uint8{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xf0}, "imei-0000000000000200", "IMEI", ""},

	// IMEISV Cases
	{[]uint8{0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xf0}, "imeisv-0000000000000200", "IMEISV", ""},

	// Error Cases
	{[]uint8{}, "", "", "empty buffer"},
	{[]uint8{0x00}, "", "", "no identity"},
}

func TestNasTypeGetMobileIdentity(t *testing.T) {
	for i, table := range GetMobileIdentityTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer

		id, idType, err := a.GetMobileIdentity()

		if table.outErr != "" {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), table.outErr)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, table.outID, id)
			assert.Equal(t, table.outType, idType)
		}
	}
}

// --- Malicious Packet Tests for Individual Getters ---

type MaliciousGetterTemplate struct {
	inBuffer []uint8
	out      string
}

var MaliciousGetterTable = []MaliciousGetterTemplate{
	{[]uint8{}, ""},             // Empty
	{[]uint8{0x00}, ""},         // Too short
	{[]uint8{0x00, 0x00}, ""},   // Too short
}

func TestNasTypeGetMCC_Malicious(t *testing.T) {
	for i, table := range MaliciousGetterTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		assert.NotPanics(t, func() {
			res := a.GetMCC()
			assert.Equal(t, "", res)
		})
	}
}

func TestNasTypeGetMNC_Malicious(t *testing.T) {
	for i, table := range MaliciousGetterTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		assert.NotPanics(t, func() {
			res := a.GetMNC()
			assert.Equal(t, "", res)
		})
	}
}

func TestNasTypeGetAmfID_Malicious(t *testing.T) {
	// GetAmfID requires at least 7 bytes
	maliciousTable := []MaliciousGetterTemplate{
		{[]uint8{}, ""},
		{[]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}, ""}, // 6 bytes (too short)
	}
	for i, table := range maliciousTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		assert.NotPanics(t, func() {
			res := a.GetAmfID()
			assert.Equal(t, "", res)
		})
	}
}

func TestNasTypeGetAmfRegionID_Malicious(t *testing.T) {
	// GetAmfRegionID requires at least 5 bytes
	maliciousTable := []MaliciousGetterTemplate{
		{[]uint8{}, ""},
		{[]uint8{0x00, 0x01, 0x02, 0x03}, ""}, // 4 bytes (too short)
	}
	for i, table := range maliciousTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		assert.NotPanics(t, func() {
			res := a.GetAmfRegionID()
			assert.Equal(t, "", res)
		})
	}
}

func TestNasTypeGetAmfSetID_Malicious(t *testing.T) {
	// GetAmfSetID logic depends on Identity Type
	// 5G-GUTI (0xF2): Start at index 5, need 2 bytes -> total 7 bytes
	// 5G-S-TMSI (0xF4): Start at index 1, need 2 bytes -> total 3 bytes
	
	maliciousTable := []struct {
		inBuffer []uint8
		desc     string
	}{
		{[]uint8{}, "Empty"},
		{[]uint8{0xf2, 0x00, 0x00, 0x00, 0x00}, "GUTI too short (5 bytes)"},
		{[]uint8{0xf4, 0x00}, "S-TMSI too short (2 bytes)"},
		{[]uint8{0x00, 0x00}, "Unknown Type too short"},
	}

	for i, table := range maliciousTable {
		t.Logf("Test Cnt:%d (%s)", i, table.desc)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		assert.NotPanics(t, func() {
			res := a.GetAmfSetID()
			assert.Equal(t, "", res)
		})
	}
}

func TestNasTypeGetAmfPointer_Malicious(t *testing.T) {
	// GetAmfPointer logic depends on Identity Type
	// 5G-GUTI (0xF2): Start at index 6
	// 5G-S-TMSI (0xF4): Start at index 2
	
	maliciousTable := []struct {
		inBuffer []uint8
		desc     string
	}{
		{[]uint8{}, "Empty"},
		{[]uint8{0xf2, 0x00, 0x00, 0x00, 0x00, 0x00}, "GUTI too short (6 bytes, need index 6)"},
		{[]uint8{0xf4, 0x00}, "S-TMSI too short (2 bytes, need index 2)"},
	}

	for i, table := range maliciousTable {
		t.Logf("Test Cnt:%d (%s)", i, table.desc)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		assert.NotPanics(t, func() {
			res := a.GetAmfPointer()
			assert.Equal(t, "", res)
		})
	}
}

func TestNasTypeGet5GTMSI_Malicious(t *testing.T) {
	// Get5GTMSI logic:
	// GUTI: need 7+ bytes (Buffer[7:])
	// S-TMSI: need 7 bytes (Buffer[3:7])
	
	maliciousTable := []struct {
		inBuffer []uint8
		desc     string
	}{
		{[]uint8{}, "Empty"},
		{[]uint8{0xf2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "GUTI too short (7 bytes, need >7)"},
		{[]uint8{0xf4, 0x00, 0x00, 0x00, 0x00, 0x00}, "S-TMSI too short (6 bytes, need 7)"},
	}

	for i, table := range maliciousTable {
		t.Logf("Test Cnt:%d (%s)", i, table.desc)
		a := nasType.NewMobileIdentity5GS(nasMessage.RegistrationRequestAdditionalGUTIType)
		a.Buffer = table.inBuffer
		assert.NotPanics(t, func() {
			res := a.Get5GTMSI()
			assert.Equal(t, "", res)
		})
	}
}