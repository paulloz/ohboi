package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func (cpu *CPU) Jump(nn uint16) {
	cpu.PC = nn
}

func init() {
	RegisterIntructions(map[uint8]Instruction{
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
	})
}
