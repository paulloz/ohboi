package apu

type channel2 struct {
	basechannel

	nr21 uint8
	nr23 uint8
	nr24 uint8
}

func (c *channel2) ReadNR13() uint8 { return c.nr21 }
func (c *channel2) ReadNR23() uint8 { return c.nr23 }
func (c *channel2) ReadNR24() uint8 { return c.nr24 }

func (c *channel2) WriteNR21(val uint8) {
	c.nr21 = val
	switch (c.nr21 >> 6) & 0x03 {
	case 1:
		c.duty = 0.75
	case 2:
		c.duty = 0.5
	case 3:
		c.duty = 0
	case 4:
		c.duty = -0.5
	}
}

func (c *channel2) WriteNR23(val uint8) {
	c.nr23 = val
	c.updateFrequency()
}

func (c *channel2) WriteNR24(val uint8) {
	c.nr24 = val
	c.updateFrequency()
}

func (c *channel2) updateFrequency() {
	frequency := (uint16((c.nr24 & 0x07)) << 8) | uint16(c.nr23)
	c.frequency = 131072 / (2048 - float64(frequency))
}

func newChannel2() *channel2 {
	c := &channel2{}
	c.active = true
	return c
}
