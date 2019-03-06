package memory

import (
	"github.com/paulloz/ohboi/cartridge"
	"github.com/paulloz/ohboi/io"
)

const (
	VRAMAddr            uint16 = 0x8000
	SwitchableRAMAddr   uint16 = 0xa000
	InternalRAMAddr     uint16 = 0xc000
	EchoInternalRAMAddr uint16 = 0xe000
	OAMAddr             uint16 = 0xfe00
	IOPortsAddr         uint16 = 0xff00
	InternalRAM2Addr    uint16 = 0xff80
)

type Memory struct {
	cartridge *cartridge.Cartridge
	io        *io.IO

	oam [0x100]uint8

	vRAM [0x2000]uint8
	hRAM [0x80]uint8

	wRAM [0x2000]uint8 // 2 4KB banks

	lastWrite uint8
}

func (mem *Memory) Read(address uint16) uint8 {
	switch {
	case address >= 0xFFFF:
		// IE
		return mem.io.Read(io.IE)

	case address >= InternalRAM2Addr:
		// High RAM
		return mem.hRAM[address-InternalRAM2Addr]

	case address >= IOPortsAddr:
		return mem.io.Read(uint8(address & 0xff))

	case address >= 0xFEA0:
		// Not usable
		return 0xFF

	case address >= OAMAddr:
		// OAM
		return mem.oam[address-OAMAddr]

	case address >= EchoInternalRAMAddr:
		// Echo on WRAM
		return mem.Read(EchoInternalRAMAddr - InternalRAMAddr)

	case address >= InternalRAMAddr:
		// Work RAM
		return mem.wRAM[address-InternalRAMAddr]

	case address >= SwitchableRAMAddr:
		// Cartridge RAM
		return mem.cartridge.Read(address)

	case address >= VRAMAddr:
		// Video RAM
		return mem.vRAM[address-VRAMAddr]

	default:
		// Cartridge ROM
		if mem.io.Read(io.BOOTROM) == 0 && address < uint16(len(bootRom)) {
			return bootRom[address]
		}

		return mem.cartridge.Read(address)
	}
}

func (mem *Memory) ReadWord(addr uint16) uint16 {
	lo := mem.Read(addr)
	hi := mem.Read(addr + 1)
	return (uint16(hi) << 8) | uint16(lo)
}

func (mem *Memory) WriteWord(addr, value uint16) {
	mem.Write(addr, uint8(value&0xff))
	mem.Write(addr+1, uint8(value>>8))
}

func (mem *Memory) Write(address uint16, value uint8) {
	switch {
	case address >= 0xFFFF:
		mem.io.Write(io.IE, value)

	case address >= InternalRAM2Addr:
		// High RAM
		mem.hRAM[address-InternalRAM2Addr] = value

	case address >= IOPortsAddr:
		mem.io.Write(uint8(address&0xff), value)

	case address >= 0xFEA0:
		// Not usable

	case address >= OAMAddr:
		// OAM
		mem.oam[address-OAMAddr] = value

	case address >= EchoInternalRAMAddr:
		// Echo to WRAM
		mem.Write(address-(EchoInternalRAMAddr-InternalRAMAddr), value)

	case address >= InternalRAMAddr:
		// Work RAM
		mem.wRAM[address-InternalRAMAddr] = value

	case address >= SwitchableRAMAddr:
		// Cartridge RAM
		mem.cartridge.Write(address, value)

	case address >= VRAMAddr:
		// Video RAM
		mem.vRAM[address-VRAMAddr] = value

	default:
		// Cartridge ROM
		mem.cartridge.Write(address, value)

	}

	mem.lastWrite = value
}

func (mem *Memory) LastWrittenValue() uint8 {
	return mem.lastWrite
}

func (mem *Memory) LoadCartridge(cartridge *cartridge.Cartridge) {
	mem.cartridge = cartridge
}

func (mem *Memory) LoadCartridgeFromFile(filename string) {
	cartridge, err := cartridge.NewCartridge(filename)
	if err != nil {
		panic(err)
	}
	mem.LoadCartridge(cartridge)
}

func (mem *Memory) Cartridge() *cartridge.Cartridge {
	return mem.cartridge
}

func NewMemory(io *io.IO) *Memory {
	mem := &Memory{io: io}

	mem.Write(0xFF05, 0x00)
	mem.Write(0xFF06, 0x00)
	mem.Write(0xFF07, 0x00)
	mem.Write(0xFF10, 0x80)
	mem.Write(0xFF11, 0xBF)
	mem.Write(0xFF12, 0xF3)
	mem.Write(0xFF14, 0xBF)
	mem.Write(0xFF16, 0x3F)
	mem.Write(0xFF17, 0x00)
	mem.Write(0xFF19, 0xBF)
	mem.Write(0xFF1A, 0x7F)
	mem.Write(0xFF1B, 0xFF)
	mem.Write(0xFF1C, 0x9F)
	mem.Write(0xFF1E, 0xBF)
	mem.Write(0xFF20, 0xFF)
	mem.Write(0xFF21, 0x00)
	mem.Write(0xFF22, 0x00)
	mem.Write(0xFF23, 0xBF)
	mem.Write(0xFF24, 0x77)
	mem.Write(0xFF25, 0xF3)
	mem.Write(0xFF26, 0xF1)
	mem.Write(0xFF40, 0x91)
	mem.Write(0xFF42, 0x00)
	mem.Write(0xFF43, 0x00)
	mem.Write(0xFF45, 0x00)
	mem.Write(0xFF47, 0xFC)
	mem.Write(0xFF48, 0xFF)
	mem.Write(0xFF49, 0xFF)
	mem.Write(0xFF4A, 0x00)
	mem.Write(0xFF4B, 0x00)
	mem.Write(0xFFFF, 0x00)

	return mem
}
