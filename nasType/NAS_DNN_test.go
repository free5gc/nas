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
	}

	a := nasType.NewDNN(0)
	for _, table := range NNTable {
		a.SetDNN(table.in)
		assert.Equalf(t, table.out, a.Buffer, "in(%v): out %v, actual %x", table.in, table.out, a.Buffer)
		assert.Equalf(t, uint8(len(table.out)), a.Len, "outlen %d, actual %d", table.in, len(table.out), a.Len)
		assert.Equalf(t, table.in, a.GetDNN(), "in(%v): GetDNN %x", table.in, a.GetDNN())
	}
}
