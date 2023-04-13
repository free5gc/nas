package nasType

type Non3GppNwPolicies struct {
	Octet uint8
}

func NewNon3GppNwPolicies(iei uint8) (non3GppNwPolicies *Non3GppNwPolicies) {
	non3GppNwPolicies = &Non3GppNwPolicies{}
	non3GppNwPolicies.SetIei(iei)
	return non3GppNwPolicies
}

func (a *Non3GppNwPolicies) GetIei() (iei uint8) {
	return a.Octet & GetBitMask(8, 4) >> (4)
}

func (a *Non3GppNwPolicies) SetIei(iei uint8) {
	a.Octet = (a.Octet & 15) + ((iei & 15) << 4)
}
