package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/nas/nasType"
)

func TestNasTypeNewDNN(t *testing.T) {
	a := nasType.NewDNN(0)
	assert.NotNil(t, a)
}

var NNIeiTable = []NasTypeIeiData{
	{0, 0},
}

func TestNasTypDNNGetSetIei(t *testing.T) {
	a := nasType.NewDNN(0)
	for _, table := range NNIeiTable {
		a.SetIei(table.in)
		assert.Equal(t, table.out, a.GetIei())
	}
}

var NNLenTable = []NasTypeLenuint8Data{
	{1, 1},
}

func TestNasTypeDNNGetSetLen(t *testing.T) {
	a := nasType.NewDNN(0)
	for _, table := range NNLenTable {
		a.SetLen(table.in)
		assert.Equal(t, table.out, a.GetLen())
	}
}

type DNNData struct {
	in  string
	out []uint8
}

func TestNasTypeDNNGetSetDNNValue(t *testing.T) {
	NNTable := []DNNData{
		{
			"internet",
			[]uint8{0x8, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74},
		},
		{
			"www.example.com",
			[]uint8{
				0x03, 0x77, 0x77, 0x77,
				0x07, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6c, 0x65,
				0x03, 0x63, 0x6f, 0x6d,
			},
		},
		{
			// length of label = 62, length of encoded buffer = 100
			"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" +
				".ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghij",
			[]uint8{
				0x3E,
				0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
				0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52,
				0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
				0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69,
				0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72,
				0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A,
				0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
				0x24,
				0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
				0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52,
				0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A,
				0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A,
			},
		},
	}

	a := nasType.NewDNN(0)
	for _, table := range NNTable {
		a.SetDNN(table.in)
		assert.Equalf(t, table.out, a.Buffer, "in(%v): out %v, actual %x", table.in, table.out, a.Buffer)
		assert.Equalf(t, uint8(len(table.out)), a.Len, "in(%v): outlen %d, actual %d", table.in, len(table.out), a.Len)
		assert.Equalf(t, table.in, a.GetDNN(), "in(%v): GetDNN %x", table.in, a.GetDNN())
	}
}

func TestNasTypeDNNSetDNNInvalidValue(t *testing.T) {
	invalidDnns := []string{
		// length of label > 62
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789A",
		// length of encoded buffer > 100
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijk",
	}

	var emptyDNN []uint8
	a := nasType.NewDNN(0)
	for _, dnn := range invalidDnns {
		a.SetDNN(dnn)
		assert.Equalf(t, emptyDNN, a.Buffer, "in(%v): out <>, actual %x", dnn, a.Buffer)
		assert.Equalf(t, uint8(0), a.Len, "in(%v): outlen 0, actual %d", dnn, a.Len)
	}
}
