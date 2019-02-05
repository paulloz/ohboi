package memory

import (
	"fmt"

	"github.com/paulloz/ohboi/cartridge"
)

// Memory ...
type Memory struct {
	cartridge *cartridge.Cartridge

	wRAM [0x2000]uint8 // 2 4KB banks
}

// Read ...
func (mem *Memory) Read(address uint16) uint8 {
	if address <= 0x7FFF { // Cartridge ROM
		return mem.cartridge.Read(address)
	}
	return 0
}

// ReadWord ...
func (mem *Memory) ReadWord(loAddress uint16, hiAddress uint16) uint16 {
	hi := mem.Read(hiAddress)
	lo := mem.Read(loAddress)
	return (uint16(hi) << 8) | uint16(lo)
}

// Write ...
func (mem *Memory) Write(address uint16, value uint8) {
	if address <= 0x7FFF { // Cartridge ROM
		mem.cartridge.Write(address, value)
		return
	} else if address <= 0x9FFF { // VRAM
		// TODO: Implement VRAM
	} else if address <= 0xBFFF { // Cartridge RAM
		mem.cartridge.Write(address, value)
	} else if address <= 0xDFFF { // WRAM
		mem.wRAM[address-0xC000] = value
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
		// TODO: Implement HRAM
	} else if address == 0xFFFF { // Interrupt Enable Register
		// TODO: Implement Interrupt Enable Register
	}

	fmt.Printf("Memory Write not implemented at address %X", address)
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

// NewMemory ...
func NewMemory() *Memory {
	return &Memory{}
}
