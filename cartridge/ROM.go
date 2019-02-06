package cartridge

// ROM ...
type ROM struct {
	rom []uint8
}

func (c *ROM) Read(address uint16) uint8 {
	return c.rom[address]
}

func (c *ROM) Write(address uint16, value uint8) {}

func NewROM(data []uint8) *ROM {
	return &ROM{
		rom: data,
	}
}
