package nas_test

import (
	"testing"

	"github.com/free5gc/nas"
)

func FuzzNAS(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		msg := new(nas.Message)
		msg.PlainNasDecode(&d)
	})
}

func FuzzGmmMessageDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		msg := new(nas.Message)
		msg.GmmMessageDecode(&d)
	})
}

func FuzzGsmMessageDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		msg := new(nas.Message)
		msg.GsmMessageDecode(&d)
	})
}
