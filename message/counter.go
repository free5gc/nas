package message

// TS 33.501 6.4.3.1
// COUNT (32 bits) := 0x00 || NAS COUNT (24 bits)
// NAS COUNT (24 bits) := NAS OVERFLOW (16 bits) || NAS SQN (8 bits)
type Count struct {
	Count uint32
	// accepted is set on a receive COUNT once a message with this COUNT passed the
	// integrity check, so Parse requires a higher COUNT next (TS 24.501 4.4.3.1).
	accepted bool
}

func (c *Count) maskTo24Bits() {
	c.Count &= 0x00ffffff
}

// Set sets the COUNT, e.g. to the start value of a new NAS security context.
// A receive COUNT set this way is not yet accepted; see MarkAccepted.
func (c *Count) Set(overflow uint16, sqn uint8) {
	c.SetOverflow(overflow)
	c.SetSQN(sqn)
	c.accepted = false
}

// MarkAccepted records that a message with the current receive COUNT was
// already accepted, e.g. for an uplink COUNT taken over from another AMF, so
// Parse does not accept that COUNT again (TS 33.501 6.4.3.2).
func (c *Count) MarkAccepted() {
	c.accepted = true
}

func (c *Count) Get() uint32 {
	c.maskTo24Bits()
	return c.Count
}

func (c *Count) AddOne() {
	c.Count++
	c.maskTo24Bits()
}

func (c *Count) SQN() uint8 {
	return uint8(c.Count & 0x000000ff)
}

func (c *Count) SetSQN(sqn uint8) {
	c.Count = (c.Count & 0xffffff00) | uint32(sqn)
}

func (c *Count) Overflow() uint16 {
	return uint16((c.Count & 0x00ffff00) >> 8)
}

func (c *Count) SetOverflow(overflow uint16) {
	c.Count = (c.Count & 0xff0000ff) | (uint32(overflow) << 8)
}
