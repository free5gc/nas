//go:build go1.18
// +build go1.18

package nas_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/free5gc/nas"
)

func FuzzNAS(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		msg := new(nas.Message)
		err := msg.PlainNasDecode(&d)
		if err == nil {
			buf, err := msg.PlainNasEncode()
			if err != nil {
				panic(fmt.Sprintf("Re-encoding failed: %s", err.Error()))
			}
			msg2 := new(nas.Message)
			err = msg2.PlainNasDecode(&buf)
			if err != nil {
				panic(fmt.Sprintf("Re-decoding failed: %s", err.Error()))
			}
		}
	})
}

func FuzzGmmMessageDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		msg := new(nas.Message)
		err := msg.GmmMessageDecode(&d)
		if err == nil {
			buf := new(bytes.Buffer)
			err := msg.GmmMessageEncode(buf)
			if err != nil {
				panic(fmt.Sprintf("Re-encoding failed: %s", err.Error()))
			}
			msg2 := new(nas.Message)
			buf2 := buf.Bytes()
			err = msg2.GmmMessageDecode(&buf2)
			if err != nil {
				panic(fmt.Sprintf("Re-decoding failed: %s", err.Error()))
			}
		}
	})
}

func FuzzGsmMessageDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		msg := new(nas.Message)
		err := msg.GsmMessageDecode(&d)
		if err == nil {
			buf := new(bytes.Buffer)
			err := msg.GsmMessageEncode(buf)
			if err != nil {
				panic(fmt.Sprintf("Re-encoding failed: %s", err.Error()))
			}
			msg2 := new(nas.Message)
			buf2 := buf.Bytes()
			err = msg2.GsmMessageDecode(&buf2)
			if err != nil {
				panic(fmt.Sprintf("Re-decoding failed: %s", err.Error()))
			}
		}
	})
}
