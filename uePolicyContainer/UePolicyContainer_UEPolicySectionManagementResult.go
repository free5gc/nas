package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
)

// refer to TS24501 v17, section D.6.3
type UEPolicySectionManagementResult struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewUEPolicySectionManagementResult(iei uint8) (uEPolicySectionManagementResult *UEPolicySectionManagementResult) {
	uEPolicySectionManagementResult = &UEPolicySectionManagementResult{}
	uEPolicySectionManagementResult.SetIei(iei)
	return uEPolicySectionManagementResult
}

func (u *UEPolicySectionManagementResult) MarshalBinary() ([]uint8, error) {
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

func (u *UEPolicySectionManagementResult) UnmarshalBinary(buf *bytes.Buffer) error {
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

func (u *UEPolicySectionManagementResult) GetIei() uint8 {
	return u.Iei
}

func (u *UEPolicySectionManagementResult) SetIei(iei uint8) {
	u.Iei = iei
}

func (u *UEPolicySectionManagementResult) GetLen() uint16 {
	return u.Len
}

func (u *UEPolicySectionManagementResult) SetLen(length uint16) {
	u.Len = length
}

// UEPolicySectionManagementResult element D.6.3.1
func (a *UEPolicySectionManagementResult) GetUEPolicySectionManagementResultContent() (ctnByte []uint8) {
	ctnByte = make([]uint8, len(a.Buffer))
	copy(ctnByte, a.Buffer)
	return
}

func (a *UEPolicySectionManagementResult) SetUEPolicySectionManagementResultContent(ctnByte []uint8) {
	a.Buffer = make([]uint8, len(ctnByte))
	copy(a.Buffer, ctnByte)
}
