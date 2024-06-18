package uePolicyContainer

import "errors"

type UEPolicyNetworkClassmark struct {
	Iei   uint8
	Len   uint8
	NSSUI uint8 // Non-subscribed SNPN signalled URSP handling indication
	Spare uint8 // always set to zero
}

func NewUEPolicyNetworkClassmark() (uEPolicyNetworkClassmark *UEPolicyNetworkClassmark) {
	uEPolicyNetworkClassmark = &UEPolicyNetworkClassmark{}
	uEPolicyNetworkClassmark.SetLen(2)
	return uEPolicyNetworkClassmark
}

// D.6.7 UE policy network classmark
// Iei Row, sBit, len = [], 8, 8
func (a *UEPolicyNetworkClassmark) GetIei() (iei uint8) {
	return a.Iei
}

// D.6.7 UE policy network classmark
// Iei Row, sBit, len = [], 8, 8
func (a *UEPolicyNetworkClassmark) SetIei(iei uint8) {
	a.Iei = iei
}

// D.6.7 UE policy network classmark
// Len Row, sBit, len = [], 8, 8
func (a *UEPolicyNetworkClassmark) GetLen() (len uint8) {
	return a.Len
}

// D.6.7 UE policy network classmark
// Len Row, sBit, len = [], 8, 8
func (a *UEPolicyNetworkClassmark) SetLen(len uint8) {
	a.Len = len
}

// D.6.7 UE policy network classmark
// Len Row, sBit, len = 1, 8, 8
// Bits of NSSUI
// 0 UE is allowed to accept URSP signalled by non-subscribed SNPNs
// 1 UE is not allowed to accept URSP signalled by non-subscribed SNPNs
func (a *UEPolicyNetworkClassmark) SetNSSUI(allowAcceptURSP uint8) error {
	if allowAcceptURSP == 0 {
		a.NSSUI = 0
	} else if allowAcceptURSP == 1 {
		a.NSSUI = 1
	} else {
		return errors.New("only allowed to set NSSUI as 0 or 1")
	}
	return nil
}

func (a *UEPolicyNetworkClassmark) GetNSSUI() uint8 {
	return a.NSSUI
}

func (a *UEPolicyNetworkClassmark) GetSpare() uint8 {
	return a.Spare
}
