package ie

import "github.com/pkg/errors"

// PagingRestriction is detailed in 9.11.3.77 Paging restriction, 24.501
// Psi holds the PSI of the types 3 and 4 only. Spare bits, and the octets after the
// last field of the type (24.007 11.4.2), are ignored on decode and not emitted.
type PagingRestriction struct {
	PagingRestrictionType uint8
	Psi
}

func (i *PagingRestriction) UnmarshalBinary(b []byte) error {
	if len(b) < 1 {
		return errors.Errorf("PagingRestriction IE length %d, expected at least 1", len(b))
	}
	i.PagingRestrictionType = Get4Bits41(b[0])
	i.Psi = Psi{}
	if i.withPSI() {
		return i.Psi.UnmarshalBinary(b[1:])
	}
	return nil
}

func (i *PagingRestriction) MarshalBinary() ([]byte, error) {
	if i.PagingRestrictionType > 15 {
		return nil, errors.Errorf("invalid paging restriction type %d", i.PagingRestrictionType)
	}
	b := []byte{i.PagingRestrictionType}
	if !i.withPSI() {
		return b, nil
	}
	psi, err := i.Psi.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(b, psi...), nil
}

// withPSI: "All paging is restricted except for specified PDU session(s)", with or without voice service.
func (i *PagingRestriction) withPSI() bool {
	return i.PagingRestrictionType == 3 || i.PagingRestrictionType == 4
}
