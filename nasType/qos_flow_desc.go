package nasType

import (
	"bytes"
	"encoding/binary"
	"io"
)

type QoSFlowDescs []QoSFlowDesc

func (q *QoSFlowDescs) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, desc := range *q {
		if descBuf, err := desc.MarshalBinary(); err != nil {
			return nil, err
		} else {
			_, err = buf.Write(descBuf)
			if err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil
}

func (q *QoSFlowDescs) UnmarshalBinary(b []byte) error {
	*q = make(QoSFlowDescs, 0)
	buf := bytes.NewBuffer(b)
	for {
		if desc, err := parseQoSFlowDesc(buf); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		} else {
			*q = append(*q, *desc)
		}
	}
}

type QoSFlowOperationCode uint8

const (
	OperationCodeCreateNewQoSFlowDescription QoSFlowOperationCode = iota + 1
	OperationCodeDeleteExistingQoSFlowDescription
	OperationCodeModifyExistingQoSFlowDescription
)

// QoSFlowDesc represnets the QoS Flow Description in NAS
type QoSFlowDesc struct {
	QFI           uint8
	OperationCode QoSFlowOperationCode
	Parameters    QoSFlowParameterList
}

func (q *QoSFlowDesc) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	// QFI
	if err := binary.Write(buf, binary.BigEndian, q.QFI); err != nil {
		return nil, err
	}

	// Operation Code
	if err := binary.Write(buf, binary.BigEndian, q.OperationCode<<5); err != nil {
		return nil, err
	}

	var E, parameterNum uint8
	if parameterNum = uint8(len(q.Parameters)); parameterNum != 0 {
		E = 1
	}
	if err := binary.Write(buf, binary.BigEndian, E<<6|parameterNum); err != nil {
		return nil, err
	}

	if E == 1 {
		if paraListBuf, err := q.Parameters.MarshalBinary(); err != nil {
			return nil, err
		} else {
			if _, err = buf.Write(paraListBuf); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

func parseQoSFlowDesc(buf *bytes.Buffer) (*QoSFlowDesc, error) {
	var q QoSFlowDesc
	// QFI
	if err := binary.Read(buf, binary.BigEndian, &q.QFI); err != nil {
		return nil, err
	}

	// Operation Code
	var OperationCodeOctet uint8
	if err := binary.Read(buf, binary.BigEndian, &OperationCodeOctet); err != nil {
		return nil, err
	}
	q.OperationCode = QoSFlowOperationCode(OperationCodeOctet >> 5)

	// Parameter List
	var ParameterNumOctet uint8
	if err := binary.Read(buf, binary.BigEndian, &ParameterNumOctet); err != nil {
		return nil, err
	}
	if ParameterNumOctet != 0 {
		if paraList, err := parseQoSFlowParameterList(buf, ParameterNumOctet&(1<<6-1)); err != nil {
			return nil, err
		} else {
			q.Parameters = paraList
		}
	}

	return &q, nil
}

type QoSFlowParameterIdentifier uint8

const (
	ParameterIdentifier5QI               QoSFlowParameterIdentifier = 0x01
	ParameterIdentifierGFBRUplink        QoSFlowParameterIdentifier = 0x02
	ParameterIdentifierGFBRDownlink      QoSFlowParameterIdentifier = 0x03
	ParameterIdentifierMFBRUplink        QoSFlowParameterIdentifier = 0x04
	ParameterIdentifierMFBRDownlink      QoSFlowParameterIdentifier = 0x05
	ParameterIdentifierAveragingWindow   QoSFlowParameterIdentifier = 0x06
	ParameterIdentifierEPSBearerIdentity QoSFlowParameterIdentifier = 0x07
)

type QoSFlowParameterList []QoSFlowParameter

func (l *QoSFlowParameterList) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for _, parameter := range *l {
		if err := binary.Write(buf, binary.BigEndian, parameter.Identifier()); err != nil {
			return nil, err
		}
		if parameterBuf, err := parameter.MarshalBinary(); err != nil {
			return nil, err
		} else {
			if err = binary.Write(buf, binary.BigEndian, uint8(len(parameterBuf))); err != nil {
				return nil, err
			}
			if _, err = buf.Write(parameterBuf); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

func parseQoSFlowParameterList(buf *bytes.Buffer, number uint8) (QoSFlowParameterList, error) {
	paraList := make(QoSFlowParameterList, 0)
	var parameterID QoSFlowParameterIdentifier
	var parameterLen uint8

	for i := 0; i < int(number); i++ {
		if err := binary.Read(buf, binary.BigEndian, &parameterID); err != nil {
			return nil, err
		}
		if err := binary.Read(buf, binary.BigEndian, &parameterLen); err != nil {
			return nil, err
		}

		parameter := newQoSFlowParameters(parameterID)

		if err := parameter.UnmarshalBinary(buf.Next(int(parameterLen))); err != nil {
			return nil, err
		}

		paraList = append(paraList, parameter)
	}

	return paraList, nil
}

type QoSFlowParameter interface {
	Identifier() QoSFlowParameterIdentifier
	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
}

func newQoSFlowParameters(id QoSFlowParameterIdentifier) QoSFlowParameter {
	var parameter QoSFlowParameter
	switch id {
	case ParameterIdentifier5QI:
		parameter = new(QoSFlow5QI)
	case ParameterIdentifierGFBRUplink:
		parameter = &QoSFlowGFBRUplink{}
	case ParameterIdentifierGFBRDownlink:
		parameter = &QoSFlowGFBRDownlink{}
	case ParameterIdentifierMFBRUplink:
		parameter = &QoSFlowMFBRUplink{}
	case ParameterIdentifierMFBRDownlink:
		parameter = &QoSFlowMFBRDownlink{}
	case ParameterIdentifierAveragingWindow:
		parameter = &QoSFlowAveragingWindow{}
	case ParameterIdentifierEPSBearerIdentity:
		parameter = &QoSFlowEBI{}
	default:
		parameter = nil
	}
	return parameter
}

type QoSFlow5QI struct {
	FiveQI uint8
}

func (q *QoSFlow5QI) Identifier() QoSFlowParameterIdentifier {
	return ParameterIdentifier5QI
}

func (q *QoSFlow5QI) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// 5QI value
	if err := binary.Write(buf, binary.BigEndian, q.FiveQI); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (q *QoSFlow5QI) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)

	// 5QI value
	if err := binary.Read(buf, binary.BigEndian, &q.FiveQI); err != nil {
		return err
	}

	return nil
}

type QoSFlowBitRateUnit uint8

const (
	QoSFlowBitRateUnit1Kbps QoSFlowBitRateUnit = iota + 1
	QoSFlowBitRateUnit4Kbps
	QoSFlowBitRateUnit16Kbps
	QoSFlowBitRateUnit64Kbps
	QoSFlowBitRateUnit256Kbps
	QoSFlowBitRateUnit1Mbps
	QoSFlowBitRateUnit4Mbps
	QoSFlowBitRateUnit16Mbps
	QoSFlowBitRateUnit64Mbps
	QoSFlowBitRateUnit256Mbps
	QoSFlowBitRateUnit1Gbps
	QoSFlowBitRateUnit4Gbps
	QoSFlowBitRateUnit16Gbps
	QoSFlowBitRateUnit64Gbps
	QoSFlowBitRateUnit256Gbps
	QoSFlowBitRateUnit1Tbps
	QoSFlowBitRateUnit4Tbps
	QoSFlowBitRateUnit16Tbps
	QoSFlowBitRateUnit64Tbps
	QoSFlowBitRateUnit256Tbps
	QoSFlowBitRateUnit1Pbps
	QoSFlowBitRateUnit4Pbps
	QoSFlowBitRateUnit16Pbps
	QoSFlowBitRateUnit64Pbps
	QoSFlowBitRateUnit256Pbps
)

type qoSFlowBitRate struct {
	Unit  QoSFlowBitRateUnit
	Value uint16
}

func (q *qoSFlowBitRate) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// BitRate Unit
	if err := binary.Write(buf, binary.BigEndian, uint8(q.Unit)); err != nil {
		return nil, err
	}

	// value
	if err := binary.Write(buf, binary.BigEndian, q.Value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (q *qoSFlowBitRate) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)

	// BitRate Unit
	if err := binary.Read(buf, binary.BigEndian, &q.Unit); err != nil {
		return err
	}

	// value
	if err := binary.Read(buf, binary.BigEndian, &q.Value); err != nil {
		return err
	}

	return nil
}

type QoSFlowGFBRUplink qoSFlowBitRate

func (q *QoSFlowGFBRUplink) Identifier() QoSFlowParameterIdentifier {
	return ParameterIdentifierGFBRUplink
}

func (q *QoSFlowGFBRUplink) MarshalBinary() ([]byte, error) {
	return (*qoSFlowBitRate)(q).MarshalBinary()
}

func (q *QoSFlowGFBRUplink) UnmarshalBinary(b []byte) error {
	return (*qoSFlowBitRate)(q).UnmarshalBinary(b)
}

type QoSFlowGFBRDownlink qoSFlowBitRate

func (q *QoSFlowGFBRDownlink) Identifier() QoSFlowParameterIdentifier {
	return ParameterIdentifierGFBRDownlink
}

func (q *QoSFlowGFBRDownlink) MarshalBinary() ([]byte, error) {
	return (*qoSFlowBitRate)(q).MarshalBinary()
}

func (q *QoSFlowGFBRDownlink) UnmarshalBinary(b []byte) error {
	return (*qoSFlowBitRate)(q).UnmarshalBinary(b)
}

type QoSFlowMFBRUplink qoSFlowBitRate

func (q *QoSFlowMFBRUplink) Identifier() QoSFlowParameterIdentifier {
	return ParameterIdentifierMFBRUplink
}

func (q *QoSFlowMFBRUplink) MarshalBinary() ([]byte, error) {
	return (*qoSFlowBitRate)(q).MarshalBinary()
}

func (q *QoSFlowMFBRUplink) UnmarshalBinary(b []byte) error {
	return (*qoSFlowBitRate)(q).UnmarshalBinary(b)
}

type QoSFlowMFBRDownlink qoSFlowBitRate

func (q *QoSFlowMFBRDownlink) Identifier() QoSFlowParameterIdentifier {
	return ParameterIdentifierMFBRDownlink
}

func (q *QoSFlowMFBRDownlink) MarshalBinary() ([]byte, error) {
	return (*qoSFlowBitRate)(q).MarshalBinary()
}

func (q *QoSFlowMFBRDownlink) UnmarshalBinary(b []byte) error {
	return (*qoSFlowBitRate)(q).UnmarshalBinary(b)
}

type QoSFlowAveragingWindow struct {
	AverageWindow uint16
}

func (q *QoSFlowAveragingWindow) Identifier() QoSFlowParameterIdentifier {
	return ParameterIdentifierAveragingWindow
}

func (q *QoSFlowAveragingWindow) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// AverageWindow value in milliseconds
	if err := binary.Write(buf, binary.BigEndian, q.AverageWindow); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (q *QoSFlowAveragingWindow) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)

	// AverageWindow value in milliseconds
	if err := binary.Read(buf, binary.BigEndian, &q.AverageWindow); err != nil {
		return err
	}

	return nil
}

type QoSFlowEBI struct {
	EBI uint8
}

func (q *QoSFlowEBI) Identifier() QoSFlowParameterIdentifier {
	return ParameterIdentifierEPSBearerIdentity
}

func (q *QoSFlowEBI) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// EBI
	if err := binary.Write(buf, binary.BigEndian, q.EBI); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (q *QoSFlowEBI) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)

	// EBI
	if err := binary.Read(buf, binary.BigEndian, &q.EBI); err != nil {
		return err
	}

	return nil
}
