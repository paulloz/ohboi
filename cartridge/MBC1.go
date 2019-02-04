package cartridge

import (
	. "github.com/paulloz/gb-emulator/types"
)

type MBC1 struct {
	rom     []Byte
	romBank Byte

	isRAMBanking bool
	isRAMEnabled bool
}

func NewMBC1(rom []Byte) *MBC1 {
	return &MBC1{
		rom:     rom,
		romBank: 1,

		isRAMBanking: false,
		isRAMEnabled: false,
	}
}

func (c *MBC1) Read(address Word) Byte {
	if address > 0x3FFF {
		return 0
	}

	return c.rom[address] // First bank is always there
}

func (c *MBC1) Write(address Word, data Byte) {
	if address <= 0x1FFF { // 0x0A on lower 4 bits enable RAM, other values disable RAM
		c.isRAMEnabled = (data & 0xF) == 0xA
	} else if address <= 0x3FFF { // Lower 5bits of romBank
		c.romBank = (c.romBank & 0x60) | (data & 0x1F)
	} else if address <= 0x5FFF {
		if c.isRAMBanking {
		} else { // Bits 5 and 6 of romBank
			c.romBank = (c.romBank & 0x1F) | (data & 0x60)
		}
	} else if address <= 0x6FFF {
		c.isRAMBanking = (data & 1) != 0
	}
}
