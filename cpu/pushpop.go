package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newPushRegister16(src Getter16, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			cpu.Push(src.Get(cpu))
			return nil
		},
		Cycles: cycles,
	}
}

func newPopRegister16(dst Setter16, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			dst.Set(cpu, cpu.Pop())
			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.PUSH_AF: newPushRegister16(RegisterAF, 16),
		op.PUSH_BC: newPushRegister16(RegisterBC, 16),
		op.PUSH_DE: newPushRegister16(RegisterDE, 16),
		op.PUSH_HL: newPushRegister16(RegisterHL, 16),

		op.POP_AF: newPopRegister16(RegisterAF, 12),
		op.POP_BC: newPopRegister16(RegisterBC, 12),
		op.POP_DE: newPopRegister16(RegisterDE, 12),
		op.POP_HL: newPopRegister16(RegisterHL, 12),
	})
}
