package nasType

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type QoSRuleOperationCode uint8

const (
	OperationCodeCreateNewQoSRule                                   QoSRuleOperationCode = 1
	OperationCodeDeleteExistingQoSRule                              QoSRuleOperationCode = 2
	OperationCodeModifyExistingQoSRuleAndAddPacketFilters           QoSRuleOperationCode = 3
	OperationCodeModifyExistingQoSRuleAndReplaceAllPacketFilters    QoSRuleOperationCode = 4
	OperationCodeModifyExistingQoSRuleAndDeletePacketFilters        QoSRuleOperationCode = 5
	OperationCodeModifyExistingQoSRuleWithoutModifyingPacketFilters QoSRuleOperationCode = 6
)

type QoSRules []QoSRule

type QoSRule struct {
	Identifier       uint8
	Operation        QoSRuleOperationCode
	DQR              bool
	PacketFilterList PacketFilterList
	Precedence       uint8
	Segregation      bool
	QFI              uint8
}

func bool2bit(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func bit2bool(n uint8) bool {
	return n != 0
}

func (q *QoSRules) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for _, rule := range *q {
		ruleBuf := bytes.NewBuffer(nil)
		if err := binary.Write(ruleBuf, binary.BigEndian, rule.Identifier); err != nil {
			return nil, err
		}

		ruleContentBuf := bytes.NewBuffer(nil)

		ruleHeader := uint8(rule.Operation)<<5 |
			bool2bit(rule.DQR)<<4 |
			uint8(len(rule.PacketFilterList))

		if err := binary.Write(ruleContentBuf, binary.BigEndian, ruleHeader); err != nil {
			return nil, err
		}

		var packetFilterError error
		var packetFilterBytes []byte
		if rule.Operation == OperationCodeModifyExistingQoSRuleAndDeletePacketFilters {
			packetFilterBytes, packetFilterError = buildPacketFilterDeleteList(rule.PacketFilterList)
		} else {
			packetFilterBytes, packetFilterError = buildPacketFilterList(rule.PacketFilterList)
		}

		if packetFilterError != nil {
			return nil, packetFilterError
		}

		if err := binary.Write(ruleContentBuf, binary.BigEndian, packetFilterBytes); err != nil {
			return nil, err
		}

		if err := binary.Write(ruleContentBuf, binary.BigEndian, rule.Precedence); err != nil {
			return nil, err
		}

		if err := binary.Write(ruleContentBuf, binary.BigEndian,
			bool2bit(rule.Segregation)<<6|rule.QFI); err != nil {
			return nil, err
		}

		if err := binary.Write(ruleBuf, binary.BigEndian, uint16(ruleContentBuf.Len())); err != nil {
			return nil, err
		}

		if _, err := ruleBuf.ReadFrom(ruleContentBuf); err != nil {
			return nil, err
		}

		if _, err := buf.ReadFrom(ruleBuf); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (q *QoSRules) UnmarshalBinary(b []byte) error {
	*q = make(QoSRules, 0)

	buf := bytes.NewBuffer(b)

	for {
		rule := QoSRule{}
		if err := binary.Read(buf, binary.BigEndian, &rule.Identifier); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		var ruleLen uint16
		if err := binary.Read(buf, binary.BigEndian, &ruleLen); err != nil {
			return err
		}

		var ruleHeader uint8
		if err := binary.Read(buf, binary.BigEndian, &ruleHeader); err != nil {
			return err
		}
		rule.Operation = QoSRuleOperationCode(ruleHeader >> 5)
		rule.DQR = bit2bool(ruleHeader & (1 << 4))
		pfLen := int(ruleHeader & 0x0F)

		var pfList PacketFilterList
		var pfListError error

		if rule.Operation == OperationCodeModifyExistingQoSRuleAndDeletePacketFilters {
			pfList, pfListError = parsePacketFilterDeleteList(buf, pfLen)
		} else {
			pfList, pfListError = parsePacketFilterList(buf, pfLen)
		}

		if pfListError != nil {
			return pfListError
		}
		rule.PacketFilterList = pfList

		if err := binary.Read(buf, binary.BigEndian, &rule.Precedence); err != nil {
			return err
		}

		var QFIByte uint8
		if err := binary.Read(buf, binary.BigEndian, &QFIByte); err != nil {
			return err
		}

		rule.Segregation = bit2bool(QFIByte >> 6)
		rule.QFI = QFIByte & (1<<6 - 1)

		*q = append(*q, rule)
	}
	return nil
}

type PacketFilterDirection uint8

const (
	PacketFilterDirectionDownlink      PacketFilterDirection = 1
	PacketFilterDirectionUplink        PacketFilterDirection = 2
	PacketFilterDirectionBidirectional PacketFilterDirection = 3
)

type PacketFilterComponentType uint8

// TS 24.501 Table 9.11.4.13.1
const (
	PacketFilterComponentTypeMatchAll                       PacketFilterComponentType = 0x01
	PacketFilterComponentTypeIPv4RemoteAddress              PacketFilterComponentType = 0x10
	PacketFilterComponentTypeIPv4LocalAddress               PacketFilterComponentType = 0x11
	PacketFilterComponentTypeIPv6RemoteAddress              PacketFilterComponentType = 0x21
	PacketFilterComponentTypeIPv6LocalAddress               PacketFilterComponentType = 0x23
	PacketFilterComponentTypeProtocolIdentifierOrNextHeader PacketFilterComponentType = 0x30
	PacketFilterComponentTypeSingleLocalPort                PacketFilterComponentType = 0x40
	PacketFilterComponentTypeLocalPortRange                 PacketFilterComponentType = 0x41
	PacketFilterComponentTypeSingleRemotePort               PacketFilterComponentType = 0x50
	PacketFilterComponentTypeRemotePortRange                PacketFilterComponentType = 0x51
	PacketFilterComponentTypeSecurityParameterIndex         PacketFilterComponentType = 0x60
	PacketFilterComponentTypeTypeOfServiceOrTrafficClass    PacketFilterComponentType = 0x70
	PacketFilterComponentTypeFlowLabel                      PacketFilterComponentType = 0x80
	PacketFilterComponentTypeDestinationMACAddress          PacketFilterComponentType = 0x81
	PacketFilterComponentTypeSourceMACAddress               PacketFilterComponentType = 0x82
	PacketFilterComponentType8021Q_CTAG_VID                 PacketFilterComponentType = 0x83
	PacketFilterComponentType8021Q_STAG_VID                 PacketFilterComponentType = 0x84
	PacketFilterComponentType8021Q_CTAG_PCPOrDEI            PacketFilterComponentType = 0x85
	PacketFilterComponentType8021Q_STAG_PCPOrDEI            PacketFilterComponentType = 0x86
	PacketFilterComponentTypeEthertype                      PacketFilterComponentType = 0x87
)

type PacketFilterList []PacketFilter

func buildPacketFilterList(pfList PacketFilterList) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for _, pf := range pfList {
		pfHeader := uint8(pf.Direction)<<4 | pf.Identifier
		if err := binary.Write(buf, binary.BigEndian, pfHeader); err != nil {
			return nil, err
		}

		if pfBuf, err := pf.Components.MarshalBinary(); err != nil {
			return nil, err
		} else {
			if err := binary.Write(buf, binary.BigEndian, uint8(len(pfBuf))); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.BigEndian, pfBuf); err != nil {
				return nil, err
			}
		}

	}

	return buf.Bytes(), nil
}

func buildPacketFilterDeleteList(pfList PacketFilterList) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for _, pf := range pfList {
		if err := binary.Write(buf, binary.BigEndian, pf.Identifier); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func parsePacketFilterList(buf *bytes.Buffer, n int) (PacketFilterList, error) {
	pfList := make(PacketFilterList, 0, n)

	for i := 0; i < n; i++ {
		var pfHeader uint8
		if err := binary.Read(buf, binary.BigEndian, &pfHeader); err != nil {
			return nil, err
		}

		var pfLen uint8
		if err := binary.Read(buf, binary.BigEndian, &pfLen); err != nil {
			return nil, err
		}

		dir := PacketFilterDirection((pfHeader & 0xF0) >> 4)
		id := pfHeader & 0x0F

		pf := PacketFilter{
			Identifier: id,
			Direction:  dir,
		}

		if err := pf.Components.UnmarshalBinary(buf.Next(int(pfLen))); err != nil {
			return nil, err
		}

		pfList = append(pfList, pf)

	}

	return pfList, nil
}

func parsePacketFilterDeleteList(buf *bytes.Buffer, n int) (PacketFilterList, error) {
	pfList := make(PacketFilterList, 0, n)

	for i := 0; i < n; i++ {
		var pfHeader uint8
		if err := binary.Read(buf, binary.BigEndian, &pfHeader); err != nil {
			return nil, err
		}

		pf := PacketFilter{
			Identifier: pfHeader & 0x0F,
		}

		pfList = append(pfList, pf)
	}

	return pfList, nil
}

func newPacketFilterComponent(id PacketFilterComponentType) PacketFilterComponent {
	switch id {
	case PacketFilterComponentTypeMatchAll:
		return &PacketFilterMatchAll{}
	case PacketFilterComponentTypeIPv4RemoteAddress:
		return &PacketFilterIPv4RemoteAddress{}
	case PacketFilterComponentTypeIPv4LocalAddress:
		return &PacketFilterIPv4LocalAddress{}
	case PacketFilterComponentTypeProtocolIdentifierOrNextHeader:
		return &PacketFilterProtocolIdentifier{}
	case PacketFilterComponentTypeSingleLocalPort:
		return &PacketFilterSingleLocalPort{}
	case PacketFilterComponentTypeLocalPortRange:
		return &PacketFilterLocalPortRange{}
	case PacketFilterComponentTypeSingleRemotePort:
		return &PacketFilterSingleRemotePort{}
	case PacketFilterComponentTypeRemotePortRange:
		return &PacketFilterRemotePortRange{}
	case PacketFilterComponentTypeSecurityParameterIndex:
		return &PacketFilterSecurityParameterIndex{}
	case PacketFilterComponentTypeTypeOfServiceOrTrafficClass:
		return &PacketFilterServiceClass{}
	case PacketFilterComponentTypeFlowLabel:
		return &PacketFilterFlowLabel{}
	case PacketFilterComponentTypeDestinationMACAddress:
		return &PacketFilterDestinationMACAddress{}
	case PacketFilterComponentTypeSourceMACAddress:
		return &PacketFilterSourceMACAddress{}
	case PacketFilterComponentType8021Q_CTAG_VID:
		return &PacketFilterCTagVID{}
	case PacketFilterComponentType8021Q_STAG_VID:
		return &PacketFilterSTagVID{}
	case PacketFilterComponentType8021Q_CTAG_PCPOrDEI:
		return &PacketFilterCTagPCPDEI{}
	case PacketFilterComponentType8021Q_STAG_PCPOrDEI:
		return &PacketFilterSTagPCPDEI{}
	case PacketFilterComponentTypeEthertype:
		return &PacketFilterEtherType{}
	}

	return nil
}

type PacketFilter struct {
	Identifier uint8
	Direction  PacketFilterDirection
	Components PacketFilterComponentList
}

type PacketFilterComponentList []PacketFilterComponent

func (p *PacketFilterComponentList) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for _, component := range *p {
		if err := binary.Write(buf, binary.BigEndian, component.Type()); err != nil {
			return nil, err
		}

		if componentBytes, err := component.MarshalBinary(); err != nil {
			return nil, err
		} else {
			if err = binary.Write(buf, binary.BigEndian, componentBytes); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

func (p *PacketFilterComponentList) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)

	for {
		var componentType uint8
		if err := binary.Read(buf, binary.BigEndian, &componentType); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		component := newPacketFilterComponent(PacketFilterComponentType(componentType))
		if component == nil {
			return errors.New("packet filter component type unknown")
		}

		componentBytes := buf.Next(component.Length())

		if err := component.UnmarshalBinary(componentBytes); err != nil {
			return fmt.Errorf("parse component[%d] fail: %s", componentType, err)
		}

		*p = append(*p, component)

	}

	return nil
}

type PacketFilterComponent interface {
	Type() PacketFilterComponentType
	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
	Length() int
}

type PacketFilterMatchAll struct{}

func (p *PacketFilterMatchAll) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeMatchAll
}

func (p *PacketFilterMatchAll) MarshalBinary() ([]byte, error) {
	return []byte{}, nil
}

func (p *PacketFilterMatchAll) UnmarshalBinary(b []byte) error {
	if len(b) != 0 {
		return errors.New("\"match all\" should be empty")
	}
	return nil
}

func (p *PacketFilterMatchAll) Length() int {
	return 0
}

type pfIPv4Address struct {
	Address net.IP
	Mask    net.IPMask
}

func (p *pfIPv4Address) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if len(p.Address) != net.IPv4len {
		return nil, fmt.Errorf("ip[%s] is not a IPv4 address", p.Address)
	}

	if len(p.Mask) != net.IPv4len {
		return nil, fmt.Errorf("mask[%s] is not a IPv4 mask", p.Mask)
	}

	if err := binary.Write(buf, binary.BigEndian, p.Address); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, p.Mask); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *pfIPv4Address) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"ipv4 address\" should be %d", p.Length())
	}

	p.Address = b[0:4]
	p.Mask = b[4:8]

	return nil
}

func (p *pfIPv4Address) Length() int {
	return 8
}

type PacketFilterIPv4RemoteAddress pfIPv4Address

func (p *PacketFilterIPv4RemoteAddress) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeIPv4RemoteAddress
}

func (p *PacketFilterIPv4RemoteAddress) MarshalBinary() ([]byte, error) {
	return (*pfIPv4Address)(p).MarshalBinary()
}

func (p *PacketFilterIPv4RemoteAddress) UnmarshalBinary(b []byte) error {
	return (*pfIPv4Address)(p).UnmarshalBinary(b)
}

func (p *PacketFilterIPv4RemoteAddress) Length() int {
	return (*pfIPv4Address)(p).Length()
}

type PacketFilterIPv4LocalAddress pfIPv4Address

func (p *PacketFilterIPv4LocalAddress) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeIPv4LocalAddress
}

func (p *PacketFilterIPv4LocalAddress) MarshalBinary() ([]byte, error) {
	return (*pfIPv4Address)(p).MarshalBinary()
}

func (p *PacketFilterIPv4LocalAddress) UnmarshalBinary(b []byte) error {
	return (*pfIPv4Address)(p).UnmarshalBinary(b)
}

func (p *PacketFilterIPv4LocalAddress) Length() int {
	return (*pfIPv4Address)(p).Length()
}

type PacketFilterProtocolIdentifier struct {
	Value uint8
}

func (p *PacketFilterProtocolIdentifier) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeProtocolIdentifierOrNextHeader
}

func (p *PacketFilterProtocolIdentifier) MarshalBinary() ([]byte, error) {
	return []byte{p.Value}, nil
}

func (p *PacketFilterProtocolIdentifier) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"protocol number\" should be %d", p.Length())
	}
	p.Value = b[0]

	return nil
}

func (p *PacketFilterProtocolIdentifier) Length() int {
	return 1
}

type pfPort struct {
	Value uint16
}

func (p *pfPort) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, p.Value)
	return b, nil
}

func (p *pfPort) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"port number\" should be %d", p.Length())
	}

	p.Value = binary.BigEndian.Uint16(b)

	return nil
}

func (p *pfPort) Length() int {
	return 2
}

type PacketFilterSingleLocalPort pfPort

func (p *PacketFilterSingleLocalPort) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeSingleLocalPort
}

func (p *PacketFilterSingleLocalPort) MarshalBinary() ([]byte, error) {
	return (*pfPort)(p).MarshalBinary()
}

func (p *PacketFilterSingleLocalPort) UnmarshalBinary(b []byte) error {
	return (*pfPort)(p).UnmarshalBinary(b)
}

func (p *PacketFilterSingleLocalPort) Length() int {
	return (*pfPort)(p).Length()
}

type PacketFilterSingleRemotePort pfPort

func (p *PacketFilterSingleRemotePort) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeSingleRemotePort
}

func (p *PacketFilterSingleRemotePort) MarshalBinary() ([]byte, error) {
	return (*pfPort)(p).MarshalBinary()
}

func (p *PacketFilterSingleRemotePort) UnmarshalBinary(b []byte) error {
	return (*pfPort)(p).UnmarshalBinary(b)
}

func (p *PacketFilterSingleRemotePort) Length() int {
	return (*pfPort)(p).Length()
}

type pfPortRange struct {
	LowLimit  uint16
	HighLimit uint16
}

func (p *pfPortRange) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if err := binary.Write(buf, binary.BigEndian, p.LowLimit); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, p.HighLimit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *pfPortRange) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"port range\" should be %d", p.Length())
	}

	buf := bytes.NewBuffer(b)

	if err := binary.Read(buf, binary.BigEndian, &p.LowLimit); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &p.HighLimit); err != nil {
		return err
	}

	return nil
}

func (p *pfPortRange) Length() int {
	return 4
}

type PacketFilterLocalPortRange pfPortRange

func (p *PacketFilterLocalPortRange) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeLocalPortRange
}

func (p *PacketFilterLocalPortRange) MarshalBinary() ([]byte, error) {
	return (*pfPortRange)(p).MarshalBinary()
}

func (p *PacketFilterLocalPortRange) UnmarshalBinary(b []byte) error {
	return (*pfPortRange)(p).UnmarshalBinary(b)
}

func (p *PacketFilterLocalPortRange) Length() int {
	return (*pfPortRange)(p).Length()
}

type PacketFilterRemotePortRange pfPortRange

func (p *PacketFilterRemotePortRange) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeRemotePortRange
}

func (p *PacketFilterRemotePortRange) MarshalBinary() ([]byte, error) {
	return (*pfPortRange)(p).MarshalBinary()
}

func (p *PacketFilterRemotePortRange) UnmarshalBinary(b []byte) error {
	return (*pfPortRange)(p).UnmarshalBinary(b)
}

func (p *PacketFilterRemotePortRange) Length() int {
	return (*pfPortRange)(p).Length()
}

type PacketFilterSecurityParameterIndex struct {
	Index uint32
}

func (p *PacketFilterSecurityParameterIndex) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeSecurityParameterIndex
}

func (p *PacketFilterSecurityParameterIndex) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, p.Index)
	return b, nil
}

func (p *PacketFilterSecurityParameterIndex) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"security parameter index\" should be %d", p.Length())
	}

	p.Index = binary.BigEndian.Uint32(b)

	return nil
}

func (p *PacketFilterSecurityParameterIndex) Length() int {
	return 4
}

type PacketFilterServiceClass struct {
	Class uint8
	Mask  uint8
}

func (p *PacketFilterServiceClass) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeTypeOfServiceOrTrafficClass
}

func (p *PacketFilterServiceClass) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(p.Class)<<8|uint16(p.Mask))
	return b, nil
}

func (p *PacketFilterServiceClass) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"service class\" should be %d", p.Length())
	}

	p.Class = uint8((binary.BigEndian.Uint16(b) & 0xFF00) >> 8)
	p.Mask = uint8(binary.BigEndian.Uint16(b) & 0x00FF)

	return nil
}

func (p *PacketFilterServiceClass) Length() int {
	return 2
}

type PacketFilterFlowLabel struct {
	Label uint32
}

func (p *PacketFilterFlowLabel) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeFlowLabel
}

func (p *PacketFilterFlowLabel) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)

	if p.Label >= (1 << 19) {
		return nil, errors.New("value of \"flow label\" should be less then 524288")
	}

	binary.BigEndian.PutUint32(b, p.Label)
	return b[1:], nil
}

func (p *PacketFilterFlowLabel) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"flow label\" should be %d", p.Length())
	}

	p.Label = uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])

	return nil
}

func (p *PacketFilterFlowLabel) Length() int {
	return 3
}

type pfMACAddress struct {
	MAC net.HardwareAddr
}

func (p *pfMACAddress) MarshalBinary() ([]byte, error) {
	return p.MAC, nil
}

func (p *pfMACAddress) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"mac address\" should be %d", p.Length())
	}

	p.MAC = b

	return nil
}

func (p *pfMACAddress) Length() int {
	return 6
}

type PacketFilterDestinationMACAddress pfMACAddress

func (p *PacketFilterDestinationMACAddress) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeDestinationMACAddress
}

func (p *PacketFilterDestinationMACAddress) MarshalBinary() ([]byte, error) {
	return (*pfMACAddress)(p).MarshalBinary()
}

func (p *PacketFilterDestinationMACAddress) UnmarshalBinary(b []byte) error {
	return (*pfMACAddress)(p).UnmarshalBinary(b)
}

func (p *PacketFilterDestinationMACAddress) Length() int {
	return (*pfMACAddress)(p).Length()
}

type PacketFilterSourceMACAddress pfMACAddress

func (p *PacketFilterSourceMACAddress) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeSourceMACAddress
}

func (p *PacketFilterSourceMACAddress) MarshalBinary() ([]byte, error) {
	return (*pfMACAddress)(p).MarshalBinary()
}

func (p *PacketFilterSourceMACAddress) UnmarshalBinary(b []byte) error {
	return (*pfMACAddress)(p).UnmarshalBinary(b)
}

func (p *PacketFilterSourceMACAddress) Length() int {
	return (*pfMACAddress)(p).Length()
}

type pfVID struct {
	VID uint16
}

func (p *pfVID) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)

	binary.BigEndian.PutUint16(b, p.VID)

	return b, nil
}

func (p *pfVID) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"vid\" should be %d", p.Length())
	}

	p.VID = binary.BigEndian.Uint16(b)

	return nil
}

func (p *pfVID) Length() int {
	return 2
}

type PacketFilterCTagVID pfVID

func (p *PacketFilterCTagVID) Type() PacketFilterComponentType {
	return PacketFilterComponentType8021Q_CTAG_VID
}

func (p *PacketFilterCTagVID) MarshalBinary() ([]byte, error) {
	return (*pfVID)(p).MarshalBinary()
}

func (p *PacketFilterCTagVID) UnmarshalBinary(b []byte) error {
	return (*pfVID)(p).UnmarshalBinary(b)
}

func (p *PacketFilterCTagVID) Length() int {
	return (*pfVID)(p).Length()
}

type PacketFilterSTagVID pfVID

func (p *PacketFilterSTagVID) Type() PacketFilterComponentType {
	return PacketFilterComponentType8021Q_STAG_VID
}

func (p *PacketFilterSTagVID) MarshalBinary() ([]byte, error) {
	return (*pfVID)(p).MarshalBinary()
}

func (p *PacketFilterSTagVID) UnmarshalBinary(b []byte) error {
	return (*pfVID)(p).UnmarshalBinary(b)
}

func (p *PacketFilterSTagVID) Length() int {
	return (*pfVID)(p).Length()
}

type pfPCPDEI struct {
	Value uint8
}

func (p *pfPCPDEI) MarshalBinary() ([]byte, error) {
	return []byte{p.Value}, nil
}

func (p *pfPCPDEI) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"PCP or DIE\" should be %d", p.Length())
	}

	p.Value = b[0]

	return nil
}

func (p *pfPCPDEI) Length() int {
	return 1
}

type PacketFilterSTagPCPDEI pfPCPDEI

func (s *PacketFilterSTagPCPDEI) Type() PacketFilterComponentType {
	return PacketFilterComponentType8021Q_STAG_PCPOrDEI
}

func (p *PacketFilterSTagPCPDEI) MarshalBinary() ([]byte, error) {
	return (*pfPCPDEI)(p).MarshalBinary()
}

func (p *PacketFilterSTagPCPDEI) UnmarshalBinary(b []byte) error {
	return (*pfPCPDEI)(p).UnmarshalBinary(b)
}

func (p *PacketFilterSTagPCPDEI) Length() int {
	return (*pfPCPDEI)(p).Length()
}

type PacketFilterCTagPCPDEI pfPCPDEI

func (s *PacketFilterCTagPCPDEI) Type() PacketFilterComponentType {
	return PacketFilterComponentType8021Q_CTAG_PCPOrDEI
}

func (p *PacketFilterCTagPCPDEI) MarshalBinary() ([]byte, error) {
	return (*pfPCPDEI)(p).MarshalBinary()
}

func (p *PacketFilterCTagPCPDEI) UnmarshalBinary(b []byte) error {
	return (*pfPCPDEI)(p).UnmarshalBinary(b)
}

func (p *PacketFilterCTagPCPDEI) Length() int {
	return (*pfPCPDEI)(p).Length()
}

type PacketFilterEtherType struct {
	EtherType uint16
}

func (s *PacketFilterEtherType) Type() PacketFilterComponentType {
	return PacketFilterComponentTypeEthertype
}

func (p *PacketFilterEtherType) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)

	binary.BigEndian.PutUint16(b, p.EtherType)

	return b, nil
}

func (p *PacketFilterEtherType) UnmarshalBinary(b []byte) error {
	if len(b) != p.Length() {
		return fmt.Errorf("length of \"ether type\" should be %d", p.Length())
	}

	p.EtherType = binary.BigEndian.Uint16(b)

	return nil
}

func (p *PacketFilterEtherType) Length() int {
	return 2
}
