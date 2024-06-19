package uePolicyContainer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type UEPolicySectionManagementListContent []UEPolicySectionManagementSubList

type UEPolicySectionManagementSubList struct {
	Len                                      uint16
	PlmnDigit1                               uint8 // ref D.6.2.3, PLMN1=MCC digit 2 + MCC digit 1
	PlmnDigit2                               uint8 // PLMN2=MNC digit 3 + MCC digit 3
	PlmnDigit3                               uint8 // PLMN3=MNC digit 2 + MNC digit 1
	Mcc                                      *int  // not specific in specm just recode the origin MCC before encoding
	Mnc                                      *int  // not specific in specm just recode the origin MNC before encoding
	UEPolicySectionManagementSubListContents UEPolicySectionManagementSubListContents
	// a UE policy section code (UPSC) containing a unique value within the PLMN or SNPN selected by the PCF.
	UpscGenerator IDGenerator
}

func (u *UEPolicySectionManagementListContent) AppendSublist(sublist UEPolicySectionManagementSubList) {
	*u = append(*u, sublist)
}

// Marshal Strcuture into byte slice
func (u *UEPolicySectionManagementListContent) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for i, sublist := range *u {
		if sublistBuf, err := sublist.MarshalBinary(); err != nil {
			return nil, err
		} else {
			_, err = buf.Write(sublistBuf)
			if err != nil {
				return nil, err
			}
			(*u)[i] = sublist
		}
	}
	return buf.Bytes(), nil
}

// UnMarshal byte slice into Strctute
func (u *UEPolicySectionManagementListContent) UnmarshalBinary(b []byte) error {
	// initial an empty slice that length and capacity = 0
	// *u = make(UEPolicySectionManagementListContent, 0)
	buf := bytes.NewBuffer(b)
	for {
		if sublist, err := parseUEPlcSublist(buf); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		} else {
			*u = append(*u, *sublist)
		}
	}
}

func (u *UEPolicySectionManagementSubList) SetLen(length uint16) {
	u.Len = length
}

func (u *UEPolicySectionManagementSubList) GetLen() uint16 {
	return u.Len
}

func (u *UEPolicySectionManagementSubList) SetPlmnDigit(mcc, mnc int) error {
	u.Mcc = &mcc
	u.Mnc = &mnc
	if *u.Mcc < 99 || *u.Mcc > 999 {
		return fmt.Errorf("MCC must be positive 3-digit, mcc:%d", u.Mcc)
	}
	if *u.Mnc < 9 || *u.Mcc > 999 {
		return fmt.Errorf("MCC must be positive 2 or 3-digit, mnc:%d", u.Mnc)
	}
	// PlmnDigit1
	u.PlmnDigit1 = (uint8((*u.Mcc%100)/10) << 4) | (uint8(*u.Mcc % 10))

	// PlmnDigit2
	if *u.Mnc < 100 {
		u.PlmnDigit2 = (0xF0) | (uint8(*u.Mcc / 100))
	} else {
		u.PlmnDigit2 = (uint8(*u.Mnc/100) << 4) | (uint8(*u.Mcc / 100))
	}

	// PlmnDigit3
	u.PlmnDigit3 = (uint8((*u.Mnc%100)/10) << 4) | (uint8(*u.Mnc % 10))

	return nil
}

func (u *UEPolicySectionManagementSubList) GetPlmnDigit() (int, int) {
	return *u.Mcc, *u.Mnc
}

func (u *UEPolicySectionManagementSubList) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	// Preprocess the content to count for length, then encode length first and content last
	contentByte, err := u.UEPolicySectionManagementSubListContents.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// Set byte length of UEPolicySectionManagementSubList--plmn1、2、3+content length
	u.SetLen(uint16(1 + 1 + 1 + len(contentByte)))

	// len
	if err := binary.Write(buf, binary.BigEndian, u.Len); err != nil {
		return nil, err
	}
	// PlmnDigit1
	if err := binary.Write(buf, binary.BigEndian, u.PlmnDigit1); err != nil {
		return nil, err
	}
	// PlmnDigit2
	if err := binary.Write(buf, binary.BigEndian, u.PlmnDigit2); err != nil {
		return nil, err
	}
	// PlmnDigit3
	if err := binary.Write(buf, binary.BigEndian, u.PlmnDigit3); err != nil {
		return nil, err
	}
	// SubListContents
	if _, err := buf.Write(contentByte); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseUEPlcSublist(buf *bytes.Buffer) (*UEPolicySectionManagementSubList, error) {
	var u UEPolicySectionManagementSubList

	// Len
	if err := binary.Read(buf, binary.BigEndian, &u.Len); err != nil {
		return nil, err
	}

	// PlmnDigit1
	if err := binary.Read(buf, binary.BigEndian, &u.PlmnDigit1); err != nil {
		return nil, err
	}
	var mccDig1, mccDig2 uint8
	mccDig1 = (0x0F) & u.PlmnDigit1
	mccDig2 = ((0xF0) & u.PlmnDigit1) >> 4
	if mccDig1 > 9 {
		return nil, fmt.Errorf("MCC Digit1 larger than 9")
	}
	if mccDig2 > 9 {
		return nil, fmt.Errorf("MCC Digit2 larger than 9")
	}

	// PlmnDigit2
	if err := binary.Read(buf, binary.BigEndian, &u.PlmnDigit2); err != nil {
		return nil, err
	}
	var mncDig3, mccDig3 uint8
	mccDig3 = (0x0F) & u.PlmnDigit2
	mncDig3 = ((0xF0) & u.PlmnDigit2) >> 4
	if mccDig3 > 9 {
		return nil, fmt.Errorf("MCC Digit3 larger than 9")
	}
	if mncDig3 > 9 {
		if mncDig3 == 15 {
			// If a network operator decides to use only two digits in the MNC, MNC digit 3 shall be coded as "1111"
			mncDig3 = 0
		} else {
			return nil, fmt.Errorf("MNC Digit3 larger than 9")
		}
	}

	// PlmnDigit3
	if err := binary.Read(buf, binary.BigEndian, &u.PlmnDigit3); err != nil {
		return nil, err
	}
	var mncDig1, mncDig2 uint8
	mncDig1 = (0x0F) & u.PlmnDigit3
	mncDig2 = ((0xF0) & u.PlmnDigit3) >> 4
	if mncDig1 > 9 {
		return nil, fmt.Errorf("MCC Digit1 larger than 9")
	}
	if mncDig2 > 9 {
		return nil, fmt.Errorf("MNC Digit2 larger than 9")
	}
	u.Mcc = new(int)
	u.Mnc = new(int)
	*u.Mcc = int(mccDig1) + int(mccDig2)*10 + int(mccDig3)*100
	*u.Mnc = int(mncDig1) + int(mncDig2)*10 + int(mncDig3)*100

	// UEPolicySectionManagementSubListContents
	if int(u.Len-3) < 0 {
		return nil, fmt.Errorf("UEPolicySectionManagementSubList length should not less than 3")
	}
	err := u.UEPolicySectionManagementSubListContents.UnmarshalBinary(buf.Next(int(u.Len - 3)))
	if err != nil {
		return nil, err
	}

	return &u, nil
}
