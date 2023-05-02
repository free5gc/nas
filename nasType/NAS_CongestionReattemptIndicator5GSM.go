package nasType

type CongestionReattemptIndicator5GSM struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewCongestionReattemptIndicator5GSM(iei uint8) (congestionReattemptIndicator5GSM *CongestionReattemptIndicator5GSM) {
	congestionReattemptIndicator5GSM = &CongestionReattemptIndicator5GSM{}
	congestionReattemptIndicator5GSM.SetIei(iei)
	return congestionReattemptIndicator5GSM
}

func (a *CongestionReattemptIndicator5GSM) GetIei() (iei uint8) {
	return a.Iei
}

func (a *CongestionReattemptIndicator5GSM) SetIei(iei uint8) {
	a.Iei = iei
}

func (a *CongestionReattemptIndicator5GSM) GetLen() (len uint8) {
	return a.Len
}

func (a *CongestionReattemptIndicator5GSM) SetLen(len uint8) {
	a.Len = len
}
