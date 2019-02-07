package memory

import (
	"fmt"

	"github.com/paulloz/ohboi/cartridge"
	op "github.com/paulloz/ohboi/cpu/opcodes"
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

	bootRom [0x100]uint8

	hRAM [0x80]uint8

	wRAM [0x2000]uint8 // 2 4KB banks
}

// Read ...
func (mem *Memory) Read(address uint16) uint8 {
	if address <= 0x7FFF { // Cartridge ROM
		return mem.cartridge.Read(address)
	} else if address <= 0x9FFF { // VRAM
		// TODO: Implement VRAM
	} else if address <= 0xBFFF { // Cartridge RAM
		return mem.cartridge.Read(address)
	} else if address <= 0xDFFF { // WRAM
		return mem.wRAM[address-0xC000]
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
	mem := &Memory{
		bootRom: [0x100]uint8{
			// Setup stack
			op.LD_SP_NN, 0xfe, 0xff,

			// Zero memory from 0x8000 to 0x9fff
			op.XOR_A,
			op.LD_HL_NN, 0xff, 0x9f,
			op.LD_HLD_A,
			op.CB,
			op.BIT_7_H,
			op.JR_NZ_N, 0xfb,

			// Setup audio
			op.LD_HL_NN, 0x26, 0xff,
			op.LD_C_N, 0x11,
			op.LD_A_N, 0x80,
			op.LD_HLD_A,
			op.LD_CADDR_A,
			op.INC_C,
			op.LD_A_N, 0xf3,
			op.LD_CADDR_A,
			op.LD_HLD_A,
			op.LD_A_N, 0x77,
			op.LD_HL_A,

			// Setup BG palette
			op.LD_A_N, 0xfc,
			op.LDH_FF00N_A, 0x47,

			// Convert and load logo data from cartridge into VRAM
			op.LD_DE_NN, 0x04, 0x01,
			op.LD_HL_NN, 0x10, 0x80,
			op.LD_A_DE,
			op.CALL_NN, 0x95, 0x00,
			op.CALL_NN, 0x96, 0x00,
			op.INC16_DE,
			op.LD_A_E,
			op.CP_N, 0x34,
			op.JR_NZ_N, 0xf3,

			// Load 8 additional bytes into VRAM
			op.LD_DE_NN, 0xd8, 0x00,
			op.LD_B_N, 0x08,
			op.LD_A_DE,
			op.INC16_DE,
			op.LD_HLI_A,
			op.INC16_HL,
			op.DEC_B,
			op.JR_NZ_N, 0xf9,

			// Setup background tilemap
			op.LD_A_N, 0x19,
			op.LD_NN_A, 0x10, 0x99,
			op.LD_HL_NN, 0x2f, 0x99,
			op.LD_C_N, 0x0c,
			op.DEC_A,
			op.JR_Z_N, 0x08,
			op.LD_HLD_A,
			op.DEC_C,
			op.JR_NZ_N, 0xf9,
			op.LD_L_N, 0x0f,
			op.JR_N, 0xf3,

			// Scroll logo on screen, and play logo sound
			// Initialize scroll count, H=0
			op.LD_H_A,
			op.LD_A_N, 0x64,
			// Set loop count
			op.LD_D_A,
			// Set vertical scroll register
			op.LDH_FF00N_A, 0x42,
			op.LD_A_N, 0x91,
			// Turn on LCD, showing Background
			op.LDH_FF00N_A, 0x40,
			// Set B=1
			op.INC_B,
			op.LD_E_N, 0x02,
			op.LD_C_N, 0x0c,
			// Wait for screen frame
			op.LDH_A_FF00N, 0x44,
			op.CP_N, 0x90,
			op.JR_NZ_N, 0xfa,
			op.DEC_C,
			op.JR_NZ_N, 0xf7,
			op.DEC_E,
			op.JR_NZ_N, 0xf2,
			op.LD_C_N, 0x13,
			// Increment scroll count
			op.INC_H,
			op.LD_A_H,
			op.LD_E_N, 0x83,
			// 0x62 counts in, play sound #1
			op.CP_N, 0x62,
			op.JR_Z_N, 0x06,
			op.LD_E_N, 0xc1,
			// 0x64 counts in, play sound #2
			op.CP_N, 0x64,
			op.JR_NZ_N, 0x06,
			// Play sound
			op.LD_A_E,
			op.LD_CADDR_A,
			op.INC_C,
			op.LD_A_N, 0x87,
			op.LD_CADDR_A,
			op.LDH_A_FF00N, 0x42,
			op.SUB_B,
			// Scroll logo up if B=1
			op.LDH_FF00N_A, 0x42,
			op.DEC_D,
			op.JR_NZ_N, 0xd2,
			// Set B=0 first time
			op.DEC_B,
			// ... next time, cause jump to "Nintendo Logo check"
			op.JR_NZ_N, 0x4f,
			// Use scrolling loop to pause
			op.LD_D_N, 0x20,
			op.JR_N, 0xcb,

			// Graphic routine
			op.LD_C_A,
			op.LD_B_N, 0x04,
			op.PUSH_BC,
			op.CB, op.RL_C,
			op.RLA,
			op.POP_BC,
			op.CB, op.RL_C,
			op.RLA,
			op.DEC_B,
			op.JR_NZ_N, 0xf5,
			op.LD_HLI_A,
			op.INC16_HL,
			op.LD_HLI_A,
			op.INC16_HL,
			op.RET,

			// Nintendo Logo
			0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
			0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
			0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,

			// More video data
			0x3C, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x3C,

			// Nintendo logo comparison routine
			// Point HL to Nintendo logo in cart
			op.LD_HL_NN, 0x04, 0x01,
			// Point DE to Nintendo logo in DMG rom
			op.LD_DE_NN, 0xa8, 0x00,
			op.LD_A_DE,
			op.INC16_DE,
			// Compare logo data in cart to DMG rom
			op.CP_HL,
			// If not a match, lock up here
			op.JR_NZ_N, 0xfe,
			op.INC16_HL,
			op.LD_A_L,
			// Do this for 0x30 bytes
			op.CP_N, 0x34,
			op.JR_NZ_N, 0xf5,
			op.LD_B_N, 0x19,
			op.LD_A_B,
			op.ADD_A_HL,
			op.INC16_HL,
			op.DEC_B,
			op.JR_NZ_N, 0xfb,
			op.ADD_A_HL,
			// If 0x19 + bytes from 0x0134-0x014D don't add to 0x00, lock up
			op.JR_NZ_N, 0xfe,
			op.LD_A_N, 0x01,
			op.LDH_FF00N_A, 0x50,
		},
	}

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
