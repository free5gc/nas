package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
	"io"
)

type UEPolicySectionManagementSubListContents []Instruction

type Instruction struct {
	Len                     uint16
	Upsc                    uint16 // a UE policy section code (UPSC) containing a value assigned by the PCF.
	UEPolicySectionContents UEPolicySectionContents
}

func (u *UEPolicySectionManagementSubListContents) AppendInstruction(ins Instruction) {
	*u = append(*u, ins)
}

// Marshal Strcuture into byte slice
func (u *UEPolicySectionManagementSubListContents) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for i, instruct := range *u {
		if instrucBuf, err := instruct.MarshalBinary(); err != nil {
			return nil, err
		} else {
			_, err = buf.Write(instrucBuf)
			if err != nil {
				return nil, err
			}
			(*u)[i] = instruct
		}
	}
	return buf.Bytes(), nil
}

// UnMarshal byte slice into Strctute
func (u *UEPolicySectionManagementSubListContents) UnmarshalBinary(b []byte) error {
	// initial an empty slice that length and capacity = 0
	// *u = make(UEPolicySectionManagementSubListContents, 0)
	buf := bytes.NewBuffer(b)
	for {
		if instrcut, err := parseInstruction(buf); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		} else {
			*u = append(*u, *instrcut)
		}
	}
}

func (i *Instruction) SetLen(len uint16) {
	i.Len = len
}

func (i *Instruction) GetLen() uint16 {
	return i.Len
}

func (i *Instruction) SetUpsc(upsc uint16) {
	i.Upsc = upsc
}

func (i *Instruction) GetUpsc() uint16 {
	return i.Upsc
}

func (i *Instruction) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	// Preprocess the content into byte slice, but append to buffer at the end
	policySectionContentBuf, err := i.UEPolicySectionContents.MarshalBinary()
	if err != nil {
		return nil, err
	}
	i.SetLen(uint16(len(policySectionContentBuf) + 2))

	// len
	if err := binary.Write(buf, binary.BigEndian, i.Len); err != nil {
		return nil, err
	}
	// UPSC
	if err := binary.Write(buf, binary.BigEndian, i.Upsc); err != nil {
		return nil, err
	}
	// Policy Section Content
	if _, err := buf.Write(policySectionContentBuf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parseInstruction(buf *bytes.Buffer) (*Instruction, error) {
	var instruction Instruction
	// Len
	if err := binary.Read(buf, binary.BigEndian, &instruction.Len); err != nil {
		return nil, err
	}
	// UPSC
	if err := binary.Read(buf, binary.BigEndian, &instruction.Upsc); err != nil {
		return nil, err
	}
	// Ue policy section contents
	if err := instruction.UEPolicySectionContents.UnmarshalBinary(buf.Next(int(instruction.Len) - 2)); err != nil {
		return nil, err
	}

	return &instruction, nil
}
