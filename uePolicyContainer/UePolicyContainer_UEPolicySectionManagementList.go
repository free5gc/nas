package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
)

// refer to TS24501 v17, section D.6.2
type UEPolicySectionManagementList struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewUEPolicySectionManagementList(iei uint8) (uEPolicySectionManagementList *UEPolicySectionManagementList) {
	uEPolicySectionManagementList = &UEPolicySectionManagementList{}
	uEPolicySectionManagementList.SetIei(iei)
	return uEPolicySectionManagementList
}

func (u *UEPolicySectionManagementList) MarshalBinary() ([]uint8, error) {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.BigEndian, u.Iei); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, u.Len); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, u.Buffer); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (u *UEPolicySectionManagementList) UnmarshalBinary(buf *bytes.Buffer) error {
	if err := binary.Read(buf, binary.BigEndian, &u.Iei); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &u.Len); err != nil {
		return err
	}
	u.Buffer = make([]uint8, u.Len)
	if err := binary.Read(buf, binary.BigEndian, u.Buffer[:u.Len]); err != nil {
		return err
	}
	return nil
}

// UEPolicySectionManagementList element D.6.2.1
// Iei Row, sBit, len = [], 8, 8
func (a *UEPolicySectionManagementList) GetIei() uint8 {
	return a.Iei
}

// UEPolicySectionManagementList element D.6.2.1
// Iei Row, sBit, len = [], 8, 8
func (a *UEPolicySectionManagementList) SetIei(iei uint8) {
	a.Iei = iei
}

// UEPolicySectionManagementList element D.6.2.1
// Len Row, sBit, len = [], 8, 16
func (a *UEPolicySectionManagementList) GetLen() uint16 {
	return a.Len
}

// UEPolicySectionManagementList element D.6.2.1
// Len Row, sBit, len = [], 8, 16
func (a *UEPolicySectionManagementList) SetLen(length uint16) {
	a.Len = length
}

// UEPolicySectionManagementList element D.6.2.1
// QoSFlowDescriptions Row, sBit, len = [0, 0], 8 , INF
func (a *UEPolicySectionManagementList) GetUEPolicySectionManagementListContent() (ctnByte []uint8) {
	ctnByte = make([]uint8, len(a.Buffer))
	copy(ctnByte, a.Buffer)
	return ctnByte
}

// set a byte list(consit of one or many sublist) to list content
func (a *UEPolicySectionManagementList) SetUEPolicySectionManagementListContent(ctnByte []uint8) {
	a.Buffer = make([]uint8, len(ctnByte))
	copy(a.Buffer, ctnByte)
}
