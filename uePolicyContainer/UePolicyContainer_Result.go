package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
	"io"
)

type UEPolicySectionManagementSubResultContents []Result

// Refer to TS23501 v17.7.1,Table D.6.3.1
type Result struct {
	Upsc                 uint16 // a UE policy section code (UPSC) containing a value assigned by the PCF.
	FailInstructionOrder uint16 // This field contains the binary encoding of the order of the failed instruction in the UE policy section management sublist.
	Cause                uint8  // The receiving entity shall treat any other value as 0110 1111, "protocol error, unspecified".
}

func (u *UEPolicySectionManagementSubResultContents) AppendResult(ins Result) {
	*u = append(*u, ins)
}

// Marshal Strcuture into byte slice
func (u *UEPolicySectionManagementSubResultContents) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for i, result := range *u {
		if resultBuf, err := result.MarshalBinary(); err != nil {
			return nil, err
		} else {
			_, err = buf.Write(resultBuf)
			if err != nil {
				return nil, err
			}
			(*u)[i] = result
		}
	}
	return buf.Bytes(), nil
}

// UnMarshal byte slice into Strctute
func (u *UEPolicySectionManagementSubResultContents) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	for {
		if result, err := parseResult(buf); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		} else {
			*u = append(*u, *result)
		}
	}
}

func NewResult() Result {
	return Result{
		Cause: 0b01101111,
	}
}

func (r *Result) SetUpsc(upsc uint16) {
	r.Upsc = upsc
}

func (r *Result) GetUpsc() uint16 {
	return r.Upsc
}

func (r *Result) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// UPSC
	if err := binary.Write(buf, binary.BigEndian, r.Upsc); err != nil {
		return nil, err
	}
	// FailInstructionOrder
	if err := binary.Write(buf, binary.BigEndian, r.FailInstructionOrder); err != nil {
		return nil, err
	}
	// Cause
	if r.Cause != 0b01101111 {
		r.Cause = 0b01101111
	}
	if err := binary.Write(buf, binary.BigEndian, r.Cause); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parseResult(buf *bytes.Buffer) (*Result, error) {
	var result Result
	// UPSC
	if err := binary.Read(buf, binary.BigEndian, &result.Upsc); err != nil {
		return nil, err
	}
	// FailInstructionOrder
	if err := binary.Read(buf, binary.BigEndian, &result.FailInstructionOrder); err != nil {
		return nil, err
	}
	// Ue policy section contents
	if err := binary.Read(buf, binary.BigEndian, &result.Cause); err != nil {
		return nil, err
	}
	if result.Cause != 0b01101111 {
		result.Cause = 0b01101111
	}

	return &result, nil
}
