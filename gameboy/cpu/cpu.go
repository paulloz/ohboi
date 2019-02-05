package cpu

import (
	"fmt"
	"os"

	"github.com/paulloz/ohboi/gameboy/memory"
)

// CPU ...
type CPU struct {
	AF Register
	BC Register
	DE Register
	HL Register

	SP Register
	PC uint16
}

// ExecuteOpCode ...
// TODO: Might want to move this to GameBoy so we don't have to pass *Memory as parameter. Also, it'd make it easier to implement stuff like DI.
func (cpu *CPU) ExecuteOpCode(opCode uint8, mem *memory.Memory) uint {
	switch opCode {
	case 0xF3: // DI
		// TODO
		return 4
	case 0xEA: // LD nn, A
		mem.Write(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()), cpu.AF.Hi())
		return 16
	case 0xE0: // LD 0xFF00+n, A
		mem.Write(0xFF00+uint16(mem.Read(cpu.AdvancePC())), cpu.AF.Hi())
		return 12
	case 0xC3: // JP nn
		cpu.jump(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12
	case 0x3E: // LD a, #
		cpu.AF.SetHi(mem.Read(cpu.AdvancePC()))
		return 8
	case 0x3C: // INC A
		cpu.inc(cpu.AF.SetHi, cpu.AF.Hi())
		return 4
	case 0x31: // LD SP, nn
		cpu.SP.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12
	case 0x00: // NOOP
		return 4
	default:
		fmt.Printf("Unknown OpCode: %X\n", opCode)
		os.Exit(0)
	}

	return 0
}

// AdvancePC return PC value and increments it
func (cpu *CPU) AdvancePC() uint16 {
	n := cpu.PC
	cpu.PC++
	return n
}

func (cpu *CPU) setZFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x40)
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xBF)
	}
}

func (cpu *CPU) setNFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x20)
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xDF)
	}
}

func (cpu *CPU) setHFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x10)
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xEF)
	}
}

func (cpu *CPU) setCFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x08)
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xF7)
	}
}

func (cpu *CPU) inc(out func(uint8), in uint8) {
	new := in + 1
	out(new)

	cpu.setZFlag(new == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(false) // TODO: Implement HalfCarry Flag
}

func (cpu *CPU) jump(nn uint16) {
	cpu.PC = nn
}

// NewCPU ...
func NewCPU() *CPU {
	cpu := &CPU{PC: 0x100}

	cpu.AF.Set(0x01B0)
	cpu.BC.Set(0x0013)
	cpu.DE.Set(0x00D8)
	cpu.HL.Set(0x014D)
	cpu.SP.Set(0xFFFE)

	return cpu
}
