package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func CopyByte(dst Setter, src Getter) MicroInstruction {
	return func(cpu *CPU, mem *memory.Memory) error {
		dst.Set(cpu, src.Get(cpu))
		return nil
	}
}

func FetchByte(value *uint8) MicroInstruction {
	return func(cpu *CPU, mem *memory.Memory) error {
		*value = cpu.FetchByte()
		return nil
	}
}

func FetchRegister(dst Setter) MicroInstruction {
	return func(cpu *CPU, mem *memory.Memory) error {
		dst.Set(cpu, cpu.FetchByte())
		return nil
	}
}

func ReadByte(hi, lo *uint8, setter Setter) MicroInstruction {
	return func(cpu *CPU, mem *memory.Memory) error {
		address := uint16(*hi)<<8 | uint16(*lo)
		v := cpu.mem.Read(address)
		setter.Set(cpu, v)
		return nil
	}
}

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
		op.LD_A_N: NewInstruction(DecodeInstruction, CopyByte(RegisterA, Immediate)),
		op.LD_B_N: NewInstruction(DecodeInstruction, CopyByte(RegisterB, Immediate)),
		op.LD_C_N: NewInstruction(DecodeInstruction, CopyByte(RegisterC, Immediate)),
		op.LD_D_N: NewInstruction(DecodeInstruction, CopyByte(RegisterD, Immediate)),
		op.LD_E_N: NewInstruction(DecodeInstruction, CopyByte(RegisterE, Immediate)),
		op.LD_H_N: NewInstruction(DecodeInstruction, CopyByte(RegisterH, Immediate)),
		op.LD_L_N: NewInstruction(DecodeInstruction, CopyByte(RegisterL, Immediate)),

		op.LD_A_A:  NewInstruction(CopyByte(RegisterA, RegisterA)),
		op.LD_A_B:  NewInstruction(CopyByte(RegisterA, RegisterB)),
		op.LD_A_C:  NewInstruction(CopyByte(RegisterA, RegisterC)),
		op.LD_A_D:  NewInstruction(CopyByte(RegisterA, RegisterD)),
		op.LD_A_E:  NewInstruction(CopyByte(RegisterA, RegisterE)),
		op.LD_A_H:  NewInstruction(CopyByte(RegisterA, RegisterH)),
		op.LD_A_L:  NewInstruction(CopyByte(RegisterA, RegisterL)),
		op.LD_A_BC: NewInstruction(DecodeInstruction, CopyByte(RegisterA, AddressBC)),
		op.LD_A_DE: NewInstruction(DecodeInstruction, CopyByte(RegisterA, AddressDE)),
		op.LD_A_HL: NewInstruction(DecodeInstruction, CopyByte(RegisterA, AddressHL)),
		op.LD_A_NN: NewInstruction(DecodeInstruction, FetchByte(&AddressImmediateOperand.lo),
			FetchByte(&AddressImmediateOperand.hi), CopyByte(RegisterA, AddressImmediateOperand)),

		op.LD_B_A:  NewInstruction(CopyByte(RegisterB, RegisterA)),
		op.LD_B_B:  NewInstruction(CopyByte(RegisterB, RegisterB)),
		op.LD_B_C:  NewInstruction(CopyByte(RegisterB, RegisterC)),
		op.LD_B_D:  NewInstruction(CopyByte(RegisterB, RegisterD)),
		op.LD_B_E:  NewInstruction(CopyByte(RegisterB, RegisterE)),
		op.LD_B_H:  NewInstruction(CopyByte(RegisterB, RegisterH)),
		op.LD_B_L:  NewInstruction(CopyByte(RegisterB, RegisterL)),
		op.LD_B_HL: NewInstruction(DecodeInstruction, CopyByte(RegisterB, AddressHL)),

		op.LD_C_A:  NewInstruction(CopyByte(RegisterC, RegisterA)),
		op.LD_C_B:  NewInstruction(CopyByte(RegisterC, RegisterB)),
		op.LD_C_C:  NewInstruction(CopyByte(RegisterC, RegisterC)),
		op.LD_C_D:  NewInstruction(CopyByte(RegisterC, RegisterD)),
		op.LD_C_E:  NewInstruction(CopyByte(RegisterC, RegisterE)),
		op.LD_C_H:  NewInstruction(CopyByte(RegisterC, RegisterH)),
		op.LD_C_L:  NewInstruction(CopyByte(RegisterC, RegisterL)),
		op.LD_C_HL: NewInstruction(DecodeInstruction, CopyByte(RegisterC, AddressHL)),

		op.LD_D_A:  NewInstruction(CopyByte(RegisterD, RegisterA)),
		op.LD_D_B:  NewInstruction(CopyByte(RegisterD, RegisterB)),
		op.LD_D_C:  NewInstruction(CopyByte(RegisterD, RegisterC)),
		op.LD_D_D:  NewInstruction(CopyByte(RegisterD, RegisterD)),
		op.LD_D_E:  NewInstruction(CopyByte(RegisterD, RegisterE)),
		op.LD_D_H:  NewInstruction(CopyByte(RegisterD, RegisterH)),
		op.LD_D_L:  NewInstruction(CopyByte(RegisterD, RegisterL)),
		op.LD_D_HL: NewInstruction(DecodeInstruction, CopyByte(RegisterD, AddressHL)),

		op.LD_E_A:  NewInstruction(CopyByte(RegisterE, RegisterA)),
		op.LD_E_B:  NewInstruction(CopyByte(RegisterE, RegisterB)),
		op.LD_E_C:  NewInstruction(CopyByte(RegisterE, RegisterC)),
		op.LD_E_D:  NewInstruction(CopyByte(RegisterE, RegisterD)),
		op.LD_E_E:  NewInstruction(CopyByte(RegisterE, RegisterE)),
		op.LD_E_H:  NewInstruction(CopyByte(RegisterE, RegisterH)),
		op.LD_E_L:  NewInstruction(CopyByte(RegisterE, RegisterL)),
		op.LD_E_HL: NewInstruction(DecodeInstruction, CopyByte(RegisterE, AddressHL)),

		op.LD_H_A:  NewInstruction(CopyByte(RegisterH, RegisterA)),
		op.LD_H_B:  NewInstruction(CopyByte(RegisterH, RegisterB)),
		op.LD_H_C:  NewInstruction(CopyByte(RegisterH, RegisterC)),
		op.LD_H_D:  NewInstruction(CopyByte(RegisterH, RegisterD)),
		op.LD_H_E:  NewInstruction(CopyByte(RegisterH, RegisterE)),
		op.LD_H_H:  NewInstruction(CopyByte(RegisterH, RegisterH)),
		op.LD_H_L:  NewInstruction(CopyByte(RegisterH, RegisterL)),
		op.LD_H_HL: NewInstruction(DecodeInstruction, CopyByte(RegisterH, AddressHL)),

		op.LD_L_A:  NewInstruction(CopyByte(RegisterL, RegisterA)),
		op.LD_L_B:  NewInstruction(CopyByte(RegisterL, RegisterB)),
		op.LD_L_C:  NewInstruction(CopyByte(RegisterL, RegisterC)),
		op.LD_L_D:  NewInstruction(CopyByte(RegisterL, RegisterD)),
		op.LD_L_E:  NewInstruction(CopyByte(RegisterL, RegisterE)),
		op.LD_L_H:  NewInstruction(CopyByte(RegisterL, RegisterH)),
		op.LD_L_L:  NewInstruction(CopyByte(RegisterL, RegisterL)),
		op.LD_L_HL: NewInstruction(DecodeInstruction, CopyByte(RegisterL, AddressHL)),

		op.LD_HL_B: NewInstruction(DecodeInstruction, CopyByte(AddressHL, RegisterB)),
		op.LD_HL_C: NewInstruction(DecodeInstruction, CopyByte(AddressHL, RegisterC)),
		op.LD_HL_D: NewInstruction(DecodeInstruction, CopyByte(AddressHL, RegisterD)),
		op.LD_HL_E: NewInstruction(DecodeInstruction, CopyByte(AddressHL, RegisterE)),
		op.LD_HL_H: NewInstruction(DecodeInstruction, CopyByte(AddressHL, RegisterH)),
		op.LD_HL_L: NewInstruction(DecodeInstruction, CopyByte(AddressHL, RegisterL)),
		op.LD_HL_N: NewInstruction(DecodeInstruction, FetchByte(&ImmediateOperand.v), CopyByte(AddressHL, ImmediateOperand)),

		op.LD_BC_A: NewInstruction(DecodeInstruction, CopyByte(AddressBC, RegisterA)),
		op.LD_DE_A: NewInstruction(DecodeInstruction, CopyByte(AddressDE, RegisterA)),
		op.LD_HL_A: NewInstruction(DecodeInstruction, CopyByte(AddressHL, RegisterA)),
		op.LD_NN_A: NewInstruction(DecodeInstruction, FetchByte(&AddressImmediateOperand.lo),
			FetchByte(&AddressImmediateOperand.hi), CopyByte(AddressImmediateOperand, RegisterA)),

		op.LD_A_CADDR: NewInstruction(DecodeInstruction, CopyByte(RegisterA, AddressC)),
		op.LD_CADDR_A: NewInstruction(DecodeInstruction, CopyByte(AddressC, RegisterA)),

		op.LD_A_HLD: NewInstruction(DecodeInstruction, CopyByte(RegisterA, AddressHLDec)),
		op.LD_HLD_A: NewInstruction(DecodeInstruction, CopyByte(AddressHLDec, RegisterA)),

		op.LD_A_HLI: NewInstruction(DecodeInstruction, CopyByte(RegisterA, AddressHLInc)),
		op.LD_HLI_A: NewInstruction(DecodeInstruction, CopyByte(AddressHLInc, RegisterA)),

		op.LDH_FF00N_A: NewInstruction(DecodeInstruction, FetchByte(&AddressFF00NOperand.v), CopyByte(AddressFF00NOperand, RegisterA)),
		op.LDH_A_FF00N: NewInstruction(DecodeInstruction, FetchByte(&AddressFF00NOperand.v), CopyByte(RegisterA, AddressFF00NOperand)),

		op.LD_BC_NN: NewInstruction(DecodeInstruction, FetchRegister(RegisterC), FetchRegister(RegisterB)),
		op.LD_DE_NN: NewInstruction(DecodeInstruction, FetchRegister(RegisterE), FetchRegister(RegisterD)),
		op.LD_HL_NN: NewInstruction(DecodeInstruction, FetchRegister(RegisterL), FetchRegister(RegisterH)),
		op.LD_SP_NN: NewInstruction(DecodeInstruction, FetchRegister(RegisterP), FetchRegister(RegisterS)),

		op.LD_SP_HL: NewInstruction(CopyByte(RegisterS, RegisterH), CopyByte(RegisterP, RegisterL)),
		op.LD_NN_SP: NewInstruction(DecodeInstruction, FetchByte(&AddressImmediateOperand.lo),
			FetchByte(&AddressImmediateOperand.hi), CopyByte(AddressImmediateOperand.Lo(), RegisterP),
			CopyByte(AddressImmediateOperand.Hi(), RegisterS)),

		op.LD_HL_SP_N: NewInstruction(
			DecodeInstruction,
			FetchByte(&ImmediateOperand.v),
			func(cpu *CPU, mem *memory.Memory) error {
				in := int32(cpu.SP.hilo)
				rel := int32(int8(ImmediateOperand.v))

				result := in + rel
				cpu.HL.Set(uint16(result))

				overflowTest := (in ^ rel ^ (result & 0xffff))

				cpu.SetZFlag(false)
				cpu.SetNFlag(false)
				cpu.SetHFlag((overflowTest & 0x10) == 0x10)
				cpu.SetCFlag((overflowTest & 0x100) == 0x100)

				return nil
			},
		),
	})
}
