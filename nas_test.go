package nas

import (
	"bytes"
	"encoding/hex"
	"io"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var hexString = "7e00560102000021e440b883d63a9f9c56b3703217152eba2010068f241c77748000b2180e54a9760068"

func TestNasGmmMessage(t *testing.T) {
	data, _ := hex.DecodeString(hexString)
	m := NewMessage()
	err := m.GmmMessageDecode(&data)
	assert.Nil(t, err)

	buff := new(bytes.Buffer)
	err = m.GmmMessageEncode(buff)
	assert.Nil(t, err)
}

func TestNasGsmMessage(t *testing.T) {
	data, _ := hex.DecodeString(hexString)
	m := NewMessage()
	err := m.GsmMessageDecode(&data)
	assert.NotNil(t, err)

	buff := new(bytes.Buffer)
	err = m.GsmMessageEncode(buff)
	assert.NotNil(t, err)
}

func TestPlainNas(t *testing.T) {
	data, _ := hex.DecodeString(hexString)
	m := NewMessage()
	err := m.PlainNasDecode(&data)
	assert.Nil(t, err)
	buff, err1 := m.PlainNasEncode()
	assert.Nil(t, err1)
	if !reflect.DeepEqual(data, buff) {
		t.Errorf("Expect : 0x%0x\nOutput: 0x%0x", data, buff)
	}
}

func TestGmmMessage(t *testing.T) {
	for _, tt := range testsGmmMessage {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fData, err := os.Open("testdata/GmmMessage/" + tt.name)
			require.NoError(t, err)
			data, err := io.ReadAll(fData)
			require.NoError(t, err)
			err = fData.Close()
			require.NoError(t, err)
			var msg Message
			err = msg.GmmMessageDecode(&data)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, msg)
			}
			buf := new(bytes.Buffer)
			err = tt.want.GmmMessageEncode(buf)
			if assert.NoError(t, err) {
				assert.Equal(t, data, buf.Bytes())
			}
		})
	}
}

func TestGsmMessage(t *testing.T) {
	for _, tt := range testsGsmMessage {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fData, err := os.Open("testdata/GsmMessage/" + tt.name)
			require.NoError(t, err)
			data, err := io.ReadAll(fData)
			require.NoError(t, err)
			err = fData.Close()
			require.NoError(t, err)
			var msg Message
			err = msg.GsmMessageDecode(&data)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, msg)
			}
			buf := new(bytes.Buffer)
			err = tt.want.GsmMessageEncode(buf)
			if assert.NoError(t, err) {
				assert.Equal(t, data, buf.Bytes())
			}
		})
	}
}

var bufGenerateBufferSlice []uint8

func generateBufferSlice(l int) []uint8 {
	bufGenerateBufferSlice = make([]byte, 256)
	for i := 0; i < 256; i++ {
		bufGenerateBufferSlice[i] = byte(i)
	}
	for len(bufGenerateBufferSlice) < math.MaxUint16 {
		bufGenerateBufferSlice = append(bufGenerateBufferSlice, bufGenerateBufferSlice...)
	}

	return bufGenerateBufferSlice[:l]
}
