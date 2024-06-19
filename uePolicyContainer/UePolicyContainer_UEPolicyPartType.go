package uePolicyContainer

type UEPolicyPartType struct {
	Octet uint8
}

const (
	UEPolicyPartType_URSP   uint8 = 0x01
	UEPolicyPartType_ANDSP  uint8 = 0x02
	UEPolicyPartType_V2XP   uint8 = 0x03
	UEPolicyPartType_ProSeP uint8 = 0x04
)

func (u *UEPolicyPartType) SetPartType(partType byte) {
	u.Octet = partType
}

func (u *UEPolicyPartType) GetPartType() byte {
	return u.Octet
}
