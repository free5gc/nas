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
