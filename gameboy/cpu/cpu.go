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
// TODO: Maybe an array of func would be better?
// TODO: There's probably a better way to handle CPU cycles count
func (cpu *CPU) ExecuteOpCode(opCode uint8, mem *memory.Memory) uint {
	switch opCode {
	case 0xF3: // DI
		// TODO: Implement DI
		return 4
	case 0xEA: // LD nn, A
		mem.Write(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()), cpu.AF.Hi())
		return 16
	case 0xE0: // LD 0xFF00+n, A
		mem.Write(0xFF00+uint16(mem.Read(cpu.AdvancePC())), cpu.AF.Hi())
		return 12
	case 0xCD: // CALL nn
		cpu.call(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12
	case 0xC3: // JP nn
		cpu.jump(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12
	case 0xA3: // AND A, E
		cpu.and(cpu.AF.SetHi, cpu.DE.Lo(), cpu.AF.Hi())
		return 4
	case 0x7D: // LD A, L
		cpu.AF.SetHi(cpu.HL.Lo())
		return 4
	case 0x7C: // LD A, H
		cpu.AF.SetHi(cpu.HL.Hi())
		return 4
	case 0x3E: // LD a, #
		cpu.AF.SetHi(mem.Read(cpu.AdvancePC()))
		return 8
	case 0x3C: // INC A
		cpu.inc(cpu.AF.SetHi, cpu.AF.Hi())
		return 4
	case 0x31: // LD SP, nn
		cpu.SP.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12
	case 0x21: // LD HL, nn
		cpu.HL.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12
	case 0x18: // JR n
		cpu.jump(cpu.PC + uint16(mem.Read(cpu.AdvancePC())))
		return 8
	case 0x00: // NOOP
		return 4
	default:
		fmt.Printf("OpCode not implemented: %X\n", opCode)
		os.Exit(0)
	}

	return 0
}

// AdvancePC returns PC value and increments it
func (cpu *CPU) AdvancePC() uint16 {
	n := cpu.PC
	cpu.PC++
	return n
}

func (cpu *CPU) setZFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x80) // 0x80 -> 1000 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0x7F) // 0x7F -> 0111 1111
	}
}

func (cpu *CPU) setNFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x40) // 0x40 -> 0100 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xBF) // 0xBF -> 1011 1111
	}
}

func (cpu *CPU) setHFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x20) // 0x20 -> 0010 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xDF) // 0xDF -> 1101 1111
	}
}

func (cpu *CPU) setCFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x10) // 0x10 -> 0001 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xEF) // 0xEF -> 1110 1111
	}
}

func (cpu *CPU) and(out func(uint8), a uint8, b uint8) {
	result := b & a
	out(result)

	cpu.setZFlag(result == 0)
	cpu.setNFlag(false)
	cpu.setHFlag(true)
	cpu.setCFlag(false)
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

func (cpu *CPU) call(nn uint16) {
	// TODO: push PC on stack
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
