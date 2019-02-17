package cartridge

type MBC1 struct {
	rom []uint8

	romBank            uint8
	activeRomBankStart uint16

	ram []uint8

	isRAMBanking       bool
	isRAMEnabled       bool
	ramBank            uint8
	activeRAMBankStart uint16
}

func NewMBC1(rom []uint8, ramSize uint16) *MBC1 {
	mbc1 := &MBC1{
		rom:     rom,
		romBank: 1,

		ram: make([]uint8, ramSize),

		isRAMBanking: false,
		isRAMEnabled: false,
	}

	mbc1.bankROM(mbc1.romBank)
	mbc1.bankRAM(mbc1.ramBank)

	return mbc1
}

func (c *MBC1) Read(address uint16) uint8 {
	if address <= 0x3FFF {
		return c.rom[address] // First bank is always there
	} else if address <= 0x7FFF {
		return c.rom[address-0x4000+c.activeRomBankStart]
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
			c.bankROM((c.romBank & 0x1F) | (value & 0x60))
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
	c.romBank = newBank
	c.activeRomBankStart = uint16(c.romBank) * uint16(0x4000)
}

func (c *MBC1) bankRAM(newBank uint8) {
	c.ramBank = newBank
	c.activeRAMBankStart = uint16(c.ramBank) * uint16(0x2000)
}
