package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newCallNN(condition func(cpu *CPU, mem *memory.Memory) bool) Instruction {
	return NewInstruction(
		DecodeInstruction,
		FetchByte(&AddressImmediateOperand.lo),
		FetchByte(&AddressImmediateOperand.hi).Combine(SkipIf(condition)),
		NoopInstruction,
		func(cpu *CPU, mem *memory.Memory) error {
			cpu.PushByte(uint8(cpu.PC >> 8))
			return nil
		},
		func(cpu *CPU, mem *memory.Memory) error {
			cpu.PushByte(uint8(cpu.PC))
			cpu.PC = AddressImmediateOperand.Address()
			return nil
		},
	)
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.CALL_NN: newCallNN(func(cpu *CPU, mem *memory.Memory) bool {
			return true
		}),

		op.CALL_C_NN: newCallNN(func(cpu *CPU, mem *memory.Memory) bool {
			return cpu.GetCFlag()
		}),

		op.CALL_NC_NN: newCallNN(func(cpu *CPU, mem *memory.Memory) bool {
			return !cpu.GetCFlag()
		}),

		op.CALL_Z_NN: newCallNN(func(cpu *CPU, mem *memory.Memory) bool {
			return cpu.GetZFlag()
		}),

		op.CALL_NZ_NN: newCallNN(func(cpu *CPU, mem *memory.Memory) bool {
			return !cpu.GetZFlag()
		}),
	})
}
