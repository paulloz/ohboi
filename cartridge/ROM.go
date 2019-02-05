package cartridge

import (
	. "github.com/paulloz/ohboi/types"
)

type ROM struct {
	rom []Byte
}

func (c *ROM) Read(address Word) Byte {
	return c.rom[address]
}

func (c *ROM) Write(address Word, data Byte) {}
