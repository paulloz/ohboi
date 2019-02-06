package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func (cpu *CPU) Call(nn uint16) {
	// TODO: push PC on stack
	cpu.PC = nn
}

func init() {
	RegisterIntructions(map[uint8]Instruction{
		op.CALL_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Call(mem.ReadWord(cpu.AdvancePC(), cpu.AdvancePC()))
				return nil
			},
			Cycles: 12,
		},
	})
}
