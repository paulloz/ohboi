package cpu

import (
	"fmt"

	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
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

// ExecuteOpCode ...
// TODO: Maybe an array of func would be better?
// TODO: There's probably a better way to handle CPU cycles count
func (cpu *CPU) ExecuteOpCode(opcode uint8) (uint, error) {
	mem := cpu.mem

	switch opcode {
	case op.NOOP:
		return 4, nil

	// Interrupt instructions
	case op.DI:
		// TODO: Implement DI
		return 4, nil

	// Load instructions
	case op.LD_NN_A:
		mem.Write(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()), cpu.AF.Hi())
		return 16, nil
	case op.LD_FF00_n_A:
		mem.Write(0xFF00+uint16(cpu.mem.Read(cpu.AdvancePC())), cpu.AF.Hi())
		return 12, nil
	case op.LD_A_L:
		cpu.AF.SetHi(cpu.HL.Lo())
		return 4, nil
	case op.LD_A_H:
		cpu.AF.SetHi(cpu.HL.Hi())
		return 4, nil
	case op.LD_A_IMM:
		cpu.AF.SetHi(mem.Read(cpu.AdvancePC()))
		return 8, nil
	case op.LD_SP_NN:
		cpu.SP.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil
	case op.LD_HL_NN:
		cpu.HL.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil

	// Call instructions
	case op.CALL_NN:
		cpu.Call(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil

	// Jump instructions
	case op.JP_NN:
		cpu.Jump(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
		return 12, nil
	case op.JR_N:
		cpu.Jump(cpu.PC + uint16(mem.Read(cpu.AdvancePC())))
		return 8, nil

	// ALU instructions
	case op.AND_A_E:
		cpu.And(cpu.AF.SetHi, cpu.DE.Lo(), cpu.AF.Hi())
		return 4, nil
	case op.INC_A:
		cpu.Inc(cpu.AF.SetHi, cpu.AF.Hi())
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

// NewCPU ...
func NewCPU(mem *memory.Memory) *CPU {
	return &CPU{
		PC: 0x0100,
		AF: NewRegister(0x01b0),
		BC: NewRegister(0x01b0),
		DE: NewRegister(0x01b0),
		HL: NewRegister(0x01b0),
		SP: NewRegister(0xfffE),

		mem: mem,
	}
}
