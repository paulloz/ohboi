package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func SkipIf(condition func(cpu *CPU, mem *memory.Memory) bool) MicroInstruction {
	return func(cpu *CPU, mem *memory.Memory) error {
		if condition(cpu, mem) {
			return nil
		}
		return ClearPipeline
	}
}

func newJumpNN(condition func(cpu *CPU, mem *memory.Memory) bool) Instruction {
	return NewInstruction(
		DecodeInstruction,
		FetchByte(&AddressImmediateOperand.lo),
		FetchByte(&AddressImmediateOperand.hi).Combine(SkipIf(condition)),
		func(cpu *CPU, mem *memory.Memory) error {
			cpu.Jump(AddressImmediateOperand.Address())
			return nil
		},
	)
}

func newJumpN(condition func(cpu *CPU, mem *memory.Memory) bool) Instruction {
	return NewInstruction(
		DecodeInstruction,
		FetchByte(&ImmediateOperand.v).Combine(SkipIf(condition)),
		func(cpu *CPU, mem *memory.Memory) error {
			cpu.Jump(uint16(int(cpu.PC) + int(int8(ImmediateOperand.v))))
			return nil
		},
	)
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.JP_NN: newJumpNN(func(cpu *CPU, mem *memory.Memory) bool {
			return true
		}),

		op.JP_Z_NN: newJumpNN(func(cpu *CPU, mem *memory.Memory) bool {
			return cpu.GetZFlag()
		}),

		op.JP_NZ_NN: newJumpNN(func(cpu *CPU, mem *memory.Memory) bool {
			return !cpu.GetZFlag()
		}),

		op.JP_C_NN: newJumpNN(func(cpu *CPU, mem *memory.Memory) bool {
			return cpu.GetCFlag()
		}),

		op.JP_NC_NN: newJumpNN(func(cpu *CPU, mem *memory.Memory) bool {
			return !cpu.GetCFlag()
		}),

		op.JP_HL: NewInstruction(func(cpu *CPU, mem *memory.Memory) error {
			cpu.Jump(cpu.HL.Get())
			return nil
		}),

		op.JR_N: newJumpN(func(cpu *CPU, mem *memory.Memory) bool {
			return true
		}),

		op.JR_C_N: newJumpN(func(cpu *CPU, mem *memory.Memory) bool {
			return cpu.GetCFlag()
		}),

		op.JR_NC_N: newJumpN(func(cpu *CPU, mem *memory.Memory) bool {
			return !cpu.GetCFlag()
		}),

		op.JR_Z_N: newJumpN(func(cpu *CPU, mem *memory.Memory) bool {
			return cpu.GetZFlag()
		}),

		op.JR_NZ_N: newJumpN(func(cpu *CPU, mem *memory.Memory) bool {
			return !cpu.GetZFlag()
		}),
	})
}
