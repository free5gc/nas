package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
	"io"
)

type UEPolicySectionContents []UEPolicyPart

type UEPolicyPart struct {
	Len                  uint16
	UEPolicyPartType     UEPolicyPartType
	UEPolicyPartContents []byte
	// the encoding way of content is specifiled in other specs according to 'UEPolicyPartType',
	// like "URSP encoding" sepific in TS 24.526 R17
}

func (u *UEPolicySectionContents) AppendUEPolicyPart(policyPart *UEPolicyPart) {
	*u = append(*u, *policyPart)
}

func (u *UEPolicySectionContents) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for i, uePolicyPart := range *u {
		if uePolicyPartBuf, err := uePolicyPart.MarshalBinary(); err != nil {
			return nil, err
		} else {
			_, err = buf.Write(uePolicyPartBuf)
			if err != nil {
				return nil, err
			}
			(*u)[i] = uePolicyPart
		}
	}
	return buf.Bytes(), nil
}

func (u *UEPolicySectionContents) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	for {
		if policyPart, err := parseUEPolicyPart(buf); err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else {
			*u = append(*u, *policyPart)
		}
	}
	return nil
}

func (u *UEPolicyPart) SetLen(len uint16) {
	u.Len = len
}

func (u *UEPolicyPart) SetLen_byContent() uint16 {
	var ctnLen uint16
	// UEPolicyPartType
	ctnLen = uint16(1)
	// UEPolicyPartContents
	ctnLen += uint16(len(u.UEPolicyPartContents))
	u.Len = ctnLen
	return ctnLen
}

func (u *UEPolicyPart) GetLen() uint16 {
	return u.Len
}

func (u *UEPolicyPart) SetPartContent(encoding_content []uint8) {
	u.UEPolicyPartContents = encoding_content
}

func (u *UEPolicyPart) GetPartContent() []uint8 {
	return u.UEPolicyPartContents
}

func (u *UEPolicyPart) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	// len
	if u.Len == 0 {
		_ = u.SetLen_byContent()
	}
	if err := binary.Write(buf, binary.BigEndian, u.Len); err != nil {
		return nil, err
	}
	// UE policy part type
	if err := binary.Write(buf, binary.BigEndian, u.UEPolicyPartType); err != nil {
		return nil, err
	}
	// UE Policy Part Contents
	if err := binary.Write(buf, binary.BigEndian, u.UEPolicyPartContents); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parseUEPolicyPart(buf *bytes.Buffer) (*UEPolicyPart, error) {
	var uEPolicyPart UEPolicyPart
	// len
	if err := binary.Read(buf, binary.BigEndian, &uEPolicyPart.Len); err != nil {
		return nil, err
	}

	// UEPolicyPartType
	if err := binary.Read(buf, binary.BigEndian, &uEPolicyPart.UEPolicyPartType); err != nil {
		return nil, err
	}

	// UEPolicyPartContents
	uEPolicyPart.UEPolicyPartContents = make([]byte, uEPolicyPart.GetLen()-1)
	if err := binary.Read(buf, binary.BigEndian, uEPolicyPart.UEPolicyPartContents[:uEPolicyPart.GetLen()-1]); err != nil {
		return nil, err
	}
	return &uEPolicyPart, nil
}
