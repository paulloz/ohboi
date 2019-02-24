package cartridge

// ROM ...
type ROM struct {
	rom []uint8
}

func (c *ROM) Read(address uint16) uint8 {
	if address < uint16(len(c.rom)) {
		return c.rom[address]
	}
	return 0xff
}

func (c *ROM) Write(address uint16, value uint8) {}

func NewROM(data []uint8) *ROM {
	return &ROM{
		rom: data,
	}
}
