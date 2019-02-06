package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func (cpu *CPU) And(out func(uint8), a uint8, b uint8) {
	result := b & a
	out(result)

	cpu.SetZFlag(result == 0)
	cpu.SetNFlag(false)
	cpu.SetHFlag(true)
	cpu.SetCFlag(false)
}

func (cpu *CPU) Inc(out func(uint8), in uint8) {
	new := in + 1
	out(new)

	cpu.SetZFlag(new == 0)
	cpu.SetNFlag(false)
	cpu.SetHFlag(false) // TODO: Implement HalfCarry Flag
}

func init() {
	RegisterIntructions(map[uint8]Instruction{
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
	})
}
