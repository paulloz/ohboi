package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

type Instruction struct {
	Handler func(cpu *CPU, mem *memory.Memory) error
	Cycles  uint
}

var InstructionSet map[uint8]Instruction

func init() {
	InstructionSet = map[uint8]Instruction{
		op.NOOP: {
			Handler: func(cpu *CPU, mem *memory.Memory) error { return nil },
			Cycles:  4,
		},

		op.LD_B_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.BC.SetHi(cpu.FetchByte())
				return nil
			},
			Cycles: 8,
		},

		op.LD_C_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.BC.SetLo(cpu.FetchByte())
				return nil
			},
			Cycles: 8,
		},

		op.LD_D_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.DE.SetHi(cpu.FetchByte())
				return nil
			},
			Cycles: 8,
		},

		op.LD_E_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.DE.SetLo(cpu.FetchByte())
				return nil
			},
			Cycles: 8,
		},

		op.LD_H_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.HL.SetHi(cpu.FetchByte())
				return nil
			},
		},

		op.LD_L_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.HL.SetLo(cpu.FetchByte())
				return nil
			},
		},

		op.LD_NN_A: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				mem.Write(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()), cpu.AF.Hi())
				return nil
			},
			Cycles: 16,
		},

		op.LD_FF00_n_A: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				mem.Write(0xFF00+uint16(cpu.mem.Read(cpu.AdvancePC())), cpu.AF.Hi())
				return nil
			},
			Cycles: 12,
		},

		op.LD_A_L: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.AF.SetHi(cpu.HL.Lo())
				return nil
			},
			Cycles: 4,
		},

		op.LD_A_H: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.AF.SetHi(cpu.HL.Hi())
				return nil
			},
			Cycles: 4,
		},

		op.LD_A_IMM: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.AF.SetHi(mem.Read(cpu.AdvancePC()))
				return nil
			},
			Cycles: 8,
		},

		op.LD_SP_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.SP.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
				return nil
			},
			Cycles: 12,
		},

		op.LD_HL_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.HL.Set(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
				return nil
			},
			Cycles: 12,
		},

		op.CALL_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Call(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
				return nil
			},
			Cycles: 12,
		},

		// Jump instructions
		op.JP_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Jump(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
				return nil
			},
			Cycles: 12,
		},

		op.JR_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Jump(cpu.PC + uint16(mem.Read(cpu.AdvancePC())))
				return nil
			},
			Cycles: 8,
		},

		op.AND_A_E: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.And(cpu.AF.SetHi, cpu.DE.Lo(), cpu.AF.Hi())
				return nil
			},
			Cycles: 4,
		},

		op.INC_A: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Inc(cpu.AF.SetHi, cpu.AF.Hi())
				return nil
			},
			Cycles: 4,
		},
	}
}
