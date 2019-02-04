package cartridge

import (
	. "github.com/paulloz/gb-emulator/types"
)

type MBC1 struct {
	rom []Byte
}

func (c MBC1) Read(address Word) Byte {
	if address <= 0x3FFF {
		return c.rom[address]
	}

	return 0
}
