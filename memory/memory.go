package memory

import (
	"fmt"

	"github.com/paulloz/ohboi/cartridge"
)

const (
	VRAMAddr            = 0x8000
	SwitchableRAMAddr   = 0xa000
	InternalRAMAddr     = 0xc000
	EchoInternalRAMAddr = 0xe000
	OAMAddr             = 0xfe00
	IOPortsAddr         = 0xff00
	InternalRAM2Addr    = 0xff80
)

// Memory ...
type Memory struct {
	cartridge *cartridge.Cartridge

	inBootMode bool

	vRAM [0x2000]uint8

	hRAM [0x80]uint8

	wRAM [0x2000]uint8 // 2 4KB banks
}

// Read ...
func (mem *Memory) Read(address uint16) uint8 {
	if address <= 0x7FFF { // Cartridge ROM
		if mem.inBootMode && address < uint16(len(bootRom)) {
			return bootRom[address]
		}
		return mem.cartridge.Read(address)
	} else if address <= 0x9FFF { // VRAM
		return mem.vRAM[address-VRAMAddr]
		// TODO: Implement VRAM
	} else if address <= 0xBFFF { // Cartridge RAM
		return mem.cartridge.Read(address)
	} else if address <= 0xDFFF { // WRAM
		return mem.wRAM[address-InternalRAMAddr]
	} else if address <= 0xFDFF { // ECHO RAM
		// TODO: Implement ECHO RAM
	} else if address <= 0xFE9F { // OAM
		// TODO: Implement OAM
	} else if address <= 0xFEFF { // Not Usable
		return 0xFF
	} else if address <= 0xFF7F { // I/O Ports
		// TODO: Implement I/O Ports
	} else if address <= 0xFFFE { // HRAM
		return mem.hRAM[address-0xFF80]
	} else if address == 0xFFFF { // Interrupt Enable Register
		// TODO: Implement Interrupt Enable Register
	}

	fmt.Printf("Memory Read not implemented at address %X\n", address)
	return 0xFF
}

// ReadWord ...
func (mem *Memory) ReadWord(addr uint16) uint16 {
	lo := mem.Read(addr)
	hi := mem.Read(addr + 1)
	return (uint16(hi) << 8) | uint16(lo)
}

func (mem *Memory) WriteWord(addr, value uint16) {
	mem.Write(addr, uint8(value&0xff))
	mem.Write(addr+1, uint8(value>>8))
}

// Write ...
func (mem *Memory) Write(address uint16, value uint8) {
	if address <= 0x7FFF { // Cartridge ROM
		mem.cartridge.Write(address, value)
		return
	} else if address <= 0x9FFF { // VRAM
		mem.vRAM[address-VRAMAddr] = value
		return
	} else if address <= 0xBFFF { // Cartridge RAM
		mem.cartridge.Write(address, value)
		return
	} else if address <= 0xDFFF { // WRAM
		mem.wRAM[address-InternalRAMAddr] = value
		return
	} else if address <= 0xFDFF { // ECHO RAM
		// TODO: Implement ECHO RAM
	} else if address <= 0xFE9F { // OAM
		// TODO: Implement OAM
	} else if address <= 0xFEFF { // Not Usable
		return
	} else if address <= 0xFF7F { // I/O Ports
		// TODO: Implement I/O Ports
	} else if address <= 0xFFFE { // HRAM
		mem.hRAM[address-0xFF80] = value
		return
	} else if address == 0xFFFF { // Interrupt Enable Register
		// TODO: Implement Interrupt Enable Register
	}

	fmt.Printf("Memory Write not implemented at address %X\n", address)
}

// LoadCartridge ...
func (mem *Memory) LoadCartridge(cartridge *cartridge.Cartridge) {
	mem.cartridge = cartridge
}

// LoadCartridgeFromFile ...
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

// NewMemory ...
func NewMemory() *Memory {
	mem := &Memory{inBootMode: true}

	/*
		// Disable for now as it generates a lot of errors

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
	*/

	return mem
}
