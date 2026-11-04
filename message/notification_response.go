package message

import (
	"bytes"
	"io"

	"github.com/pkg/errors"

	"github.com/free5gc/nas/ie"
)

var _ Message = &NotifRsp{}

// NotifRsp is detailed in 8.2.24 Notification response, 24.501
type NotifRsp struct {
	PDUSessStatus *ie.PDUSessStatus //   TLV,    4-34B, 9.11.3.44
}

const (
	NotifRspIEIPDUSessStatus uint8 = 0x50
)

func (m *NotifRsp) ExtendedProtocolDiscriminator() Epd {
	return Epd5GSMobilityMgmtMsg
}

func (m *NotifRsp) MsgType() MsgType {
	return MsgTypeNotifRsp
}

func (m *NotifRsp) UnmarshalBinary(b []byte) error {
	if len(b) < int(GmmHdrLen) {
		return errors.Errorf("NotifRsp len(b)=%d, < GmmHdrLen(%d)",
			len(b), GmmHdrLen)
	}

	// skip 3 bytes header (b should not be used after this call)
	reader := bytes.NewBuffer(b[3:])
	m.PDUSessStatus = nil
	var seenPDUSessStatus bool

	// This message contains 0 Mandatory IE
	// Optional / Conditional IE
	for reader.Len() > 0 {
		iei := reader.Next(1)[0]
		switch iei {
		case NotifRspIEIPDUSessStatus: // TLV, 4-34B
			ieLen, err := getIeLen(reader, IELen8Bits)
			if err != nil {
				// Ignore the incomplete optional IE (TS 24.501 section 7.7.1).
				// Its boundary is unavailable, so stop parsing the tail.
				return nil
			}
			value := reader.Next(int(ieLen))
			if seenPDUSessStatus {
				continue // Only the first occurrence counts (section 7.6.3).
			}
			seenPDUSessStatus = true
			if len(value) < 2 {
				continue // Syntactically incorrect optional IE: treat as absent.
			}
			status := new(ie.PDUSessStatus)
			if err = status.UnmarshalBinary(value); err != nil {
				return errors.Wrap(err, "NotifRsp.PDUSessStatus.UnmarshalBinary")
			}
			m.PDUSessStatus = status
		default:
			if err := skipUnknownOptionalIE(reader, iei); err != nil {
				if errors.Is(err, io.ErrUnexpectedEOF) {
					// Ignore an incomplete optional tail, keeping decoded IEs.
					return nil
				}
				return errors.Wrap(err, "NotifRsp")
			}
		}
	}
	return nil
}

func (m *NotifRsp) MarshalBinary() ([]byte, error) {
	writer := bytes.NewBuffer([]byte{
		byte(Epd5GSMobilityMgmtMsg),
		byte(SecHdrTypePlainNas),
		byte(MsgTypeNotifRsp),
	})

	// m.PDUSessStatus TLV, 4-34B, IEI=0x50
	if m.PDUSessStatus != nil {
		out, err := m.PDUSessStatus.MarshalBinary()
		if err != nil {
			return nil, errors.Wrap(err, "NotifRsp.PDUSessStatus.MarshalBinary()")
		}
		writer.WriteByte(NotifRspIEIPDUSessStatus)
		writer.WriteByte(byte(len(out)))
		writer.Write(out)
	}
	return writer.Bytes(), nil
}
