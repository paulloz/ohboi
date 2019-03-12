package cartridge

import (
	"io/ioutil"
	"os"
)

type MBC1 struct {
	rom []uint8

	romBank            uint8
	romBanks           uint8
	activeRomBankStart uint32

	ram []uint8

	isRAMBanking       bool
	isRAMEnabled       bool
	ramBank            uint8
	ramBankMask        uint8
	activeRAMBankStart uint16
}

func NewMBC1(rom []uint8, romBanks uint8, ramSize uint16, ramBankMask uint8) *MBC1 {
	mbc1 := &MBC1{
		rom:      rom,
		romBank:  1,
		romBanks: romBanks,

		ram:         make([]uint8, ramSize),
		ramBankMask: ramBankMask,

		isRAMBanking: false,
		isRAMEnabled: false,
	}

	mbc1.bankROM(mbc1.romBank)
	mbc1.bankRAM(mbc1.ramBank)

	return mbc1
}

func (c *MBC1) GetMode() bool {
	return c.isRAMBanking
}

func (c *MBC1) GetRAMState() (bool, int, uint8, uint16) {
	return c.isRAMEnabled, len(c.ram), c.ramBank, c.activeRAMBankStart
}

func (c *MBC1) GetROMState() (uint8, uint32) {
	return c.romBank, c.activeRomBankStart
}

func (c *MBC1) Read(address uint16) uint8 {
	if address <= 0x3FFF {
		return c.rom[address] // First bank is always there
	} else if address <= 0x7FFF {
		if addr := uint32(address) - 0x4000 + c.activeRomBankStart; addr < uint32(len(c.rom)) {
			return c.rom[addr]
		}
		return 0xff
	}

	if c.isRAMEnabled {
		if addr := address - 0xA000 + c.activeRAMBankStart; addr < uint16(len(c.ram)) {
			return c.ram[addr]
		}
	}

	return 0xff
}

func (c *MBC1) Write(address uint16, value uint8) {
	if address <= 0x1FFF { // 0x0A on lower 4 bits enable RAM, other values disable RAM
		c.isRAMEnabled = (value & 0xF) == 0xA
	} else if address <= 0x3FFF { // Lower 5bits of romBank
		c.bankROM((c.romBank & 0x60) | (value & 0x1F))
	} else if address <= 0x5FFF {
		if c.isRAMBanking { // 2 bits of ramBank
			c.bankRAM((c.ramBank & 0xFC) | (value & 0x3))
		} else { // Bits 5 and 6 of romBank
			// c.bankROM((c.romBank & 0x1F) | (value & 0x60))
			c.bankROM((c.romBank & 0x1F) | ((value & 0x3) << 5))
		}
	} else if address <= 0x7FFF {
		c.isRAMBanking = (value & 1) != 0
		if !c.isRAMBanking {
			c.bankRAM(0)
		}
	} else if address >= 0xA000 && address <= 0xBFFF { // Writing to RAM
		if c.isRAMEnabled {
			if addr := address - 0xA000 + c.activeRAMBankStart; addr < uint16(len(c.ram)) {
				c.ram[address-0xA000+c.activeRAMBankStart] = value
			}
		}
	}
}

func (c *MBC1) bankROM(newBank uint8) {
	if newBank == 0x00 || newBank == 0x20 || newBank == 0x40 || newBank == 0x60 {
		newBank++
	}
	if newBank < c.romBanks {
		c.romBank = newBank
		c.activeRomBankStart = uint32(c.romBank) * 0x4000
	}
}

func (c *MBC1) bankRAM(newBank uint8) {
	c.ramBank = newBank & c.ramBankMask
	c.activeRAMBankStart = uint16(c.ramBank) * uint16(0x2000)
}

func (c *MBC1) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if _, err := f.Write(c.ram); err != nil {
		return err
	}

	return f.Close()
}

func (c *MBC1) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	copy(c.ram, content)
	return nil
}
