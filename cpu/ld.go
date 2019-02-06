package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newLoadRegister(dst Setter, src Getter) Instruction {
	cycles := uint(4)
	if src == AddressHL {
		cycles = 8
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			dst.Set(cpu, src.Get(cpu))
			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterIntructions(map[uint8]Instruction{
		op.LD_B_N: newLoadRegister(RegisterB, Immediate),
		op.LD_C_N: newLoadRegister(RegisterC, Immediate),
		op.LD_D_N: newLoadRegister(RegisterD, Immediate),
		op.LD_E_N: newLoadRegister(RegisterE, Immediate),
		op.LD_H_N: newLoadRegister(RegisterH, Immediate),
		op.LD_L_N: newLoadRegister(RegisterL, Immediate),

		op.LD_A_A:  NoopInstruction,
		op.LD_A_B:  newLoadRegister(RegisterA, RegisterB),
		op.LD_A_C:  newLoadRegister(RegisterA, RegisterC),
		op.LD_A_D:  newLoadRegister(RegisterA, RegisterD),
		op.LD_A_E:  newLoadRegister(RegisterA, RegisterE),
		op.LD_A_H:  newLoadRegister(RegisterA, RegisterH),
		op.LD_A_L:  newLoadRegister(RegisterA, RegisterL),
		op.LD_A_HL: newLoadRegister(RegisterA, AddressHL),

		op.LD_B_B:  NoopInstruction,
		op.LD_B_C:  newLoadRegister(RegisterB, RegisterC),
		op.LD_B_D:  newLoadRegister(RegisterB, RegisterD),
		op.LD_B_E:  newLoadRegister(RegisterB, RegisterE),
		op.LD_B_H:  newLoadRegister(RegisterB, RegisterH),
		op.LD_B_L:  newLoadRegister(RegisterB, RegisterL),
		op.LD_B_HL: newLoadRegister(RegisterB, AddressHL),

		op.LD_C_B:  newLoadRegister(RegisterC, RegisterB),
		op.LD_C_C:  NoopInstruction,
		op.LD_C_D:  newLoadRegister(RegisterC, RegisterD),
		op.LD_C_E:  newLoadRegister(RegisterC, RegisterE),
		op.LD_C_H:  newLoadRegister(RegisterC, RegisterH),
		op.LD_C_L:  newLoadRegister(RegisterC, RegisterL),
		op.LD_C_HL: newLoadRegister(RegisterC, AddressHL),

		op.LD_D_B:  newLoadRegister(RegisterD, RegisterB),
		op.LD_D_C:  newLoadRegister(RegisterD, RegisterC),
		op.LD_D_D:  NoopInstruction,
		op.LD_D_E:  newLoadRegister(RegisterD, RegisterE),
		op.LD_D_H:  newLoadRegister(RegisterD, RegisterH),
		op.LD_D_L:  newLoadRegister(RegisterD, RegisterL),
		op.LD_D_HL: newLoadRegister(RegisterD, AddressHL),

		op.LD_E_B:  newLoadRegister(RegisterE, RegisterB),
		op.LD_E_C:  newLoadRegister(RegisterE, RegisterC),
		op.LD_E_D:  newLoadRegister(RegisterE, RegisterD),
		op.LD_E_E:  NoopInstruction,
		op.LD_E_H:  newLoadRegister(RegisterE, RegisterH),
		op.LD_E_L:  newLoadRegister(RegisterE, RegisterL),
		op.LD_E_HL: newLoadRegister(RegisterE, AddressHL),

		op.LD_H_B:  newLoadRegister(RegisterH, RegisterB),
		op.LD_H_C:  newLoadRegister(RegisterH, RegisterC),
		op.LD_H_D:  newLoadRegister(RegisterH, RegisterD),
		op.LD_H_E:  newLoadRegister(RegisterH, RegisterE),
		op.LD_H_H:  NoopInstruction,
		op.LD_H_L:  newLoadRegister(RegisterH, RegisterL),
		op.LD_H_HL: newLoadRegister(RegisterH, AddressHL),

		op.LD_L_B:  newLoadRegister(RegisterL, RegisterB),
		op.LD_L_C:  newLoadRegister(RegisterL, RegisterC),
		op.LD_L_D:  newLoadRegister(RegisterL, RegisterD),
		op.LD_L_E:  newLoadRegister(RegisterL, RegisterE),
		op.LD_L_H:  newLoadRegister(RegisterL, RegisterH),
		op.LD_L_L:  NoopInstruction,
		op.LD_L_HL: newLoadRegister(RegisterL, AddressHL),

		op.LD_HL_B: newLoadRegister(AddressHL, RegisterB),
		op.LD_HL_C: newLoadRegister(AddressHL, RegisterC),
		op.LD_HL_D: newLoadRegister(AddressHL, RegisterD),
		op.LD_HL_E: newLoadRegister(AddressHL, RegisterE),
		op.LD_HL_H: newLoadRegister(AddressHL, RegisterH),
		op.LD_HL_L: newLoadRegister(AddressHL, RegisterL),
		op.LD_HL_N: newLoadRegister(AddressHL, Immediate),

		op.LD_A_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.AF.SetHi(mem.Read(cpu.AdvancePC()))
				return nil
			},
			Cycles: 8,
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
	})
}
