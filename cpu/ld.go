package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newLoadRegister(dst Setter, src Getter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			dst.Set(cpu, src.Get(cpu))
			return nil
		},
		Cycles: cycles,
	}
}

func newLoadRegister16(dst Setter16, src Getter16, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			dst.Set(cpu, src.Get(cpu))
			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.LD_A_N: newLoadRegister(RegisterA, Immediate, 8),
		op.LD_B_N: newLoadRegister(RegisterB, Immediate, 8),
		op.LD_C_N: newLoadRegister(RegisterC, Immediate, 8),
		op.LD_D_N: newLoadRegister(RegisterD, Immediate, 8),
		op.LD_E_N: newLoadRegister(RegisterE, Immediate, 8),
		op.LD_H_N: newLoadRegister(RegisterH, Immediate, 8),
		op.LD_L_N: newLoadRegister(RegisterL, Immediate, 8),

		op.LD_A_A:  NoopInstruction,
		op.LD_A_B:  newLoadRegister(RegisterA, RegisterB, 4),
		op.LD_A_C:  newLoadRegister(RegisterA, RegisterC, 4),
		op.LD_A_D:  newLoadRegister(RegisterA, RegisterD, 4),
		op.LD_A_E:  newLoadRegister(RegisterA, RegisterE, 4),
		op.LD_A_H:  newLoadRegister(RegisterA, RegisterH, 4),
		op.LD_A_L:  newLoadRegister(RegisterA, RegisterL, 4),
		op.LD_A_BC: newLoadRegister(RegisterL, AddressBC, 8),
		op.LD_A_DE: newLoadRegister(RegisterL, AddressDE, 8),
		op.LD_A_HL: newLoadRegister(RegisterA, AddressHL, 8),
		op.LD_A_NN: newLoadRegister(RegisterA, AddressImmediate, 16),

		op.LD_B_B:  NoopInstruction,
		op.LD_B_C:  newLoadRegister(RegisterB, RegisterC, 4),
		op.LD_B_D:  newLoadRegister(RegisterB, RegisterD, 4),
		op.LD_B_E:  newLoadRegister(RegisterB, RegisterE, 4),
		op.LD_B_H:  newLoadRegister(RegisterB, RegisterH, 4),
		op.LD_B_L:  newLoadRegister(RegisterB, RegisterL, 4),
		op.LD_B_HL: newLoadRegister(RegisterB, AddressHL, 8),

		op.LD_C_B:  newLoadRegister(RegisterC, RegisterB, 4),
		op.LD_C_C:  NoopInstruction,
		op.LD_C_D:  newLoadRegister(RegisterC, RegisterD, 4),
		op.LD_C_E:  newLoadRegister(RegisterC, RegisterE, 4),
		op.LD_C_H:  newLoadRegister(RegisterC, RegisterH, 4),
		op.LD_C_L:  newLoadRegister(RegisterC, RegisterL, 4),
		op.LD_C_HL: newLoadRegister(RegisterC, AddressHL, 8),

		op.LD_D_B:  newLoadRegister(RegisterD, RegisterB, 4),
		op.LD_D_C:  newLoadRegister(RegisterD, RegisterC, 4),
		op.LD_D_D:  NoopInstruction,
		op.LD_D_E:  newLoadRegister(RegisterD, RegisterE, 4),
		op.LD_D_H:  newLoadRegister(RegisterD, RegisterH, 4),
		op.LD_D_L:  newLoadRegister(RegisterD, RegisterL, 4),
		op.LD_D_HL: newLoadRegister(RegisterD, AddressHL, 8),

		op.LD_E_B:  newLoadRegister(RegisterE, RegisterB, 4),
		op.LD_E_C:  newLoadRegister(RegisterE, RegisterC, 4),
		op.LD_E_D:  newLoadRegister(RegisterE, RegisterD, 4),
		op.LD_E_E:  NoopInstruction,
		op.LD_E_H:  newLoadRegister(RegisterE, RegisterH, 4),
		op.LD_E_L:  newLoadRegister(RegisterE, RegisterL, 4),
		op.LD_E_HL: newLoadRegister(RegisterE, AddressHL, 8),

		op.LD_H_B:  newLoadRegister(RegisterH, RegisterB, 4),
		op.LD_H_C:  newLoadRegister(RegisterH, RegisterC, 4),
		op.LD_H_D:  newLoadRegister(RegisterH, RegisterD, 4),
		op.LD_H_E:  newLoadRegister(RegisterH, RegisterE, 4),
		op.LD_H_H:  NoopInstruction,
		op.LD_H_L:  newLoadRegister(RegisterH, RegisterL, 4),
		op.LD_H_HL: newLoadRegister(RegisterH, AddressHL, 8),

		op.LD_L_B:  newLoadRegister(RegisterL, RegisterB, 4),
		op.LD_L_C:  newLoadRegister(RegisterL, RegisterC, 4),
		op.LD_L_D:  newLoadRegister(RegisterL, RegisterD, 4),
		op.LD_L_E:  newLoadRegister(RegisterL, RegisterE, 4),
		op.LD_L_H:  newLoadRegister(RegisterL, RegisterH, 4),
		op.LD_L_L:  NoopInstruction,
		op.LD_L_HL: newLoadRegister(RegisterL, AddressHL, 8),

		op.LD_HL_B: newLoadRegister(AddressHL, RegisterB, 8),
		op.LD_HL_C: newLoadRegister(AddressHL, RegisterC, 8),
		op.LD_HL_D: newLoadRegister(AddressHL, RegisterD, 8),
		op.LD_HL_E: newLoadRegister(AddressHL, RegisterE, 8),
		op.LD_HL_H: newLoadRegister(AddressHL, RegisterH, 8),
		op.LD_HL_L: newLoadRegister(AddressHL, RegisterL, 8),
		op.LD_HL_N: newLoadRegister(AddressHL, Immediate, 12),

		op.LD_B_A:  newLoadRegister(RegisterB, RegisterA, 4),
		op.LD_C_A:  newLoadRegister(RegisterC, RegisterA, 4),
		op.LD_D_A:  newLoadRegister(RegisterD, RegisterA, 4),
		op.LD_E_A:  newLoadRegister(RegisterE, RegisterA, 4),
		op.LD_H_A:  newLoadRegister(RegisterH, RegisterA, 4),
		op.LD_L_A:  newLoadRegister(RegisterL, RegisterA, 4),
		op.LD_BC_A: newLoadRegister(AddressBC, RegisterA, 8),
		op.LD_DE_A: newLoadRegister(AddressDE, RegisterA, 8),
		op.LD_HL_A: newLoadRegister(AddressHL, RegisterA, 8),
		op.LD_NN_A: newLoadRegister(AddressImmediate, RegisterA, 16),

		op.LD_A_CADDR: newLoadRegister(RegisterA, AddressC, 8),
		op.LD_CADDR_A: newLoadRegister(AddressC, RegisterA, 8),

		op.LD_A_HLD: newLoadRegister(RegisterA, AddressHLDec, 8),
		op.LD_HLD_A: newLoadRegister(AddressHLDec, RegisterA, 8),

		op.LD_A_HLI: newLoadRegister(RegisterA, AddressHLDec, 8),
		op.LD_HLI_A: newLoadRegister(AddressHLDec, RegisterA, 8),

		op.LDH_FF00N_A: newLoadRegister(AddressFF00N, RegisterA, 12),
		op.LDH_A_FF00N: newLoadRegister(RegisterA, AddressFF00N, 12),

		op.LD_BC_NN: newLoadRegister16(RegisterBC, Immediate16, 12),
		op.LD_DE_NN: newLoadRegister16(RegisterDE, Immediate16, 12),
		op.LD_HL_NN: newLoadRegister16(RegisterHL, Immediate16, 12),
		op.LD_SP_NN: newLoadRegister16(RegisterSP, Immediate16, 12),

		op.LD_SP_HL:   newLoadRegister16(RegisterSP, RegisterHL, 8),
		op.LD_HL_SP_N: newLoadRegister16(RegisterHL, AddressSPN, 12),
		op.LD_NN_SP:   newLoadRegister16(AddressImmediate16, RegisterSP, 20),
	})
}
