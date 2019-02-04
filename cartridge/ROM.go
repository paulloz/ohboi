package cartridge

import (
	. "github.com/paulloz/gb-emulator/types"
)

type ROM struct {
	rom []Byte
}

func (c *ROM) Read(address Word) Byte {
	return c.rom[address]
}
