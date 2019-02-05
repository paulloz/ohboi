package cartridge

import (
	. "github.com/paulloz/ohboi/types"
)

type MBC1 struct {
	rom []Byte

	romBank            Byte
	activeRomBankStart Word

	ram []Byte

	isRAMBanking       bool
	isRAMEnabled       bool
	ramBank            Byte
	activeRamBankStart Word
}

func NewMBC1(rom []Byte, ramSize Word) *MBC1 {
	return &MBC1{
		rom:     rom,
		romBank: 1,

		ram: make([]Byte, ramSize),

		isRAMBanking: false,
		isRAMEnabled: false,
	}
}

func (c *MBC1) Read(address Word) Byte {
	if address <= 0x3FFF {
		return c.rom[address] // First bank is always there
	} else if address <= 0x7FFF {
		return c.rom[address-0x4000+c.activeRomBankStart]
	}

	return c.ram[address-0xA000+c.activeRamBankStart]
}

func (c *MBC1) Write(address Word, data Byte) {
	if address <= 0x1FFF { // 0x0A on lower 4 bits enable RAM, other values disable RAM
		c.isRAMEnabled = (data & 0xF) == 0xA
	} else if address <= 0x3FFF { // Lower 5bits of romBank
		c.bankROM((c.romBank & 0x60) | (data & 0x1F))
	} else if address <= 0x5FFF {
		if c.isRAMBanking { // 2 bits of ramBank
			c.bankRAM((c.ramBank & 0xFC) | (data & 0x3))
		} else { // Bits 5 and 6 of romBank
			c.bankROM((c.romBank & 0x1F) | (data & 0x60))
		}
	} else if address <= 0x7FFF {
		c.isRAMBanking = (data & 1) != 0
		if !c.isRAMBanking {
			c.bankRAM(0)
		}
	} else if address >= 0xA000 && address <= 0xBFFF { // Writing to RAM
		if c.isRAMEnabled {
			c.ram[address-0xA000+c.activeRamBankStart] = data
		}
	}
}

func (c *MBC1) bankROM(newBank Byte) {
	if newBank == 0x00 || newBank == 0x20 || newBank == 0x40 || newBank == 0x60 {
		newBank += 1
	}
	c.romBank = newBank
	c.activeRomBankStart = Word(c.romBank) * Word(0x4000)
}

func (c *MBC1) bankRAM(newBank Byte) {
	c.ramBank = newBank
	c.activeRamBankStart = Word(c.ramBank) * Word(0x2000)
}
