package cpu

import (
	"fmt"

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

	mem *memory.Memory
}

const (
	NOOP        = 0x00
	DI          = 0xf3
	LD_NN_A     = 0xea
	LD_FF00_n_A = 0xe0
	LD_A_H      = 0x7c
	LD_A_L      = 0x7d
	LD_A_IMM    = 0x3e
	LD_SP_NN    = 0x31
	LD_HL_NN    = 0x21
	CALL_NN     = 0xcd
	AND_A_E     = 0xa3
	INC_A       = 0x3c
	JP_NN       = 0xc3
	JR_N        = 0x18
)

// ExecuteOpCode ...
// TODO: Maybe an array of func would be better?
// TODO: There's probably a better way to handle CPU cycles count
func (cpu *CPU) ExecuteOpCode(opcode uint8) (uint, error) {
	mem := cpu.mem

	switch opcode {
	case NOOP:
		return 4, nil

	// Interrupt instructions
	case DI:
		// TODO: Implement DI
		return 4, nil

	// Load instructions
	case LD_NN_A:
		mem.Write(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()), cpu.AF.Hi())
		return 16, nil
	case LD_FF00_n_A:
		mem.Write(0xFF00+uint16(cpu.mem.Read(cpu.AdvancePC())), cpu.AF.Hi())
		return 12, nil
	case LD_A_L:
		cpu.AF.SetHi(cpu.HL.Lo())
		return 4, nil
	case LD_A_H:
		cpu.AF.SetHi(cpu.HL.Hi())
		return 4, nil
	case LD_A_IMM:
		cpu.AF.SetHi(mem.Read(cpu.AdvancePC()))
		return 8, nil
	case LD_SP_NN:
		cpu.SP.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil
	case LD_HL_NN:
		cpu.HL.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil

	// Call instructions
	case CALL_NN:
		cpu.call(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil

	// Jump instructions
	case JP_NN:
		cpu.jump(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil
	case JR_N:
		cpu.jump(cpu.PC + uint16(mem.Read(cpu.AdvancePC())))
		return 8, nil

	// ALU instructions
	case AND_A_E:
		cpu.and(cpu.AF.SetHi, cpu.DE.Lo(), cpu.AF.Hi())
		return 4, nil
	case INC_A:
		cpu.inc(cpu.AF.SetHi, cpu.AF.Hi())
		return 4, nil

	default:
		return 0, fmt.Errorf("Opcode %X not implemented\n", opcode)
	}
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
func NewCPU(mem *memory.Memory) *CPU {
	return &CPU{
		PC: 0x0100,
		AF: NewRegister(0x01b0),
		BC: NewRegister(0x01b0),
		DE: NewRegister(0x01b0),
		HL: NewRegister(0x01b0),
		SP: NewRegister(0xfffE),
	}
}
