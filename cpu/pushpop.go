package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newPushRegister16(src Getter16) Instruction {
	return Instruction{
		MicroInstructions: []MicroInstruction{
			DecodeInstruction,
			NoopInstruction,
			func(cpu *CPU, mem *memory.Memory) error {
				cpu.PushByte(uint8(src.Get(cpu) >> 8))
				return nil
			},
			func(cpu *CPU, mem *memory.Memory) error {
				cpu.PushByte(uint8(src.Get(cpu)))
				return nil
			},
		},
	}
}

func newPopRegister16(dst Setter16) Instruction {
	return Instruction{
		MicroInstructions: []MicroInstruction{
			DecodeInstruction,
			func(cpu *CPU, mem *memory.Memory) error {
				dst.SetLow(cpu, cpu.PopByte())
				return nil
			},
			func(cpu *CPU, mem *memory.Memory) error {
				dst.SetHigh(cpu, cpu.PopByte())
				return nil
			},
		},
	}
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.PUSH_AF: newPushRegister16(RegisterAF),
		op.PUSH_BC: newPushRegister16(RegisterBC),
		op.PUSH_DE: newPushRegister16(RegisterDE),
		op.PUSH_HL: newPushRegister16(RegisterHL),

		op.POP_AF: newPopRegister16(RegisterAF),
		op.POP_BC: newPopRegister16(RegisterBC),
		op.POP_DE: newPopRegister16(RegisterDE),
		op.POP_HL: newPopRegister16(RegisterHL),
	})
}
