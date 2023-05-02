package nasType

type EPSBearerContextStatus struct {
	Iei   uint8
	Len   uint8
	Octet [2]uint8
}

func NewEPSBearerContextStatus(iei uint8) (ePSBearerContextStatus *EPSBearerContextStatus) {
	ePSBearerContextStatus = &EPSBearerContextStatus{}
	ePSBearerContextStatus.SetIei(iei)
	return ePSBearerContextStatus
}

func (a *EPSBearerContextStatus) GetIei() (iei uint8) {
	return a.Iei
}

func (a *EPSBearerContextStatus) SetIei(iei uint8) {
	a.Iei = iei
}

func (a *EPSBearerContextStatus) GetLen() (len uint8) {
	return a.Len
}

func (a *EPSBearerContextStatus) SetLen(len uint8) {
	a.Len = len
}
