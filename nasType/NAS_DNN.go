package nasType

import (
	"bytes"
	"errors"
	"strings"
)

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
type DNN struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewDNN(iei uint8) (dNN *DNN) {
	dNN = &DNN{}
	dNN.SetIei(iei)
	return dNN
}

// DNN 9.11.2.1A
// Iei Row, sBit, len = [], 8, 8
func (a *DNN) GetIei() (iei uint8) {
	return a.Iei
}

// DNN 9.11.2.1A
// Iei Row, sBit, len = [], 8, 8
func (a *DNN) SetIei(iei uint8) {
	a.Iei = iei
}

// DNN 9.11.2.1A
// Len Row, sBit, len = [], 8, 8
func (a *DNN) GetLen() (len uint8) {
	return a.Len
}

// DNN 9.11.2.1A
// Len Row, sBit, len = [], 8, 8
func (a *DNN) SetLen(len uint8) {
	a.Len = len
	a.Buffer = make([]uint8, a.Len)
}

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
func (a *DNN) GetDNN() string {
	return rfc1035tofqdn(a.Buffer)
}

// DNN 9.11.2.1A
// DNN Row, sBit, len = [0, 0], 8 , INF
func (a *DNN) SetDNN(dNN string) {
	if b, err := fqdnToRfc1035(dNN); err == nil {
		a.Buffer = b
		a.Len = uint8(len(a.Buffer))
	}
}

func fqdnToRfc1035(fqdn string) ([]byte, error) {
	var rfc1035RR []byte
	domainSegments := strings.Split(fqdn, ".")

	for _, segment := range domainSegments {
		if len(segment) > 63 {
			return nil, errors.New("fqdn limit the label to 63 octets or less")
		}
		rfc1035RR = append(rfc1035RR, uint8(len(segment)))
		rfc1035RR = append(rfc1035RR, segment...)
	}

	if len(rfc1035RR) > 255 {
		return nil, errors.New("fqdn should less then 255 octet")
	}
	return rfc1035RR, nil
}

func rfc1035tofqdn(rfc1035RR []byte) string {
	rfc1035Reader := bytes.NewBuffer(rfc1035RR)
	fqdn := ""

	for {
		// length of label
		if labelLen, err := rfc1035Reader.ReadByte(); err != nil {
			break
		} else {
			fqdn += string(rfc1035Reader.Next(int(labelLen))) + "."
		}
	}

	return fqdn[0 : len(fqdn)-1]
}
