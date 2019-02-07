package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func init() {
	RegisterIntructions(map[uint8]Instruction{
		op.RET: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.PC = cpu.Pop()
				return nil
			},
			Cycles: 8,
		},

		op.RET_C: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.Pop()
				if cpu.GetCFlag() {
					cpu.PC = addr
				}
				return nil
			},
			Cycles: 8,
		},

		op.RET_NC: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.Pop()
				if !cpu.GetCFlag() {
					cpu.PC = addr
				}
				return nil
			},
			Cycles: 8,
		},

		op.RET_Z: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.Pop()
				if cpu.GetZFlag() {
					cpu.PC = addr
				}
				return nil
			},
			Cycles: 8,
		},

		op.RET_NZ: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.Pop()
				if !cpu.GetZFlag() {
					cpu.PC = addr
				}
				return nil
			},
			Cycles: 8,
		},

		op.RETI: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.Pop()
				cpu.PC = addr
				cpu.EnableInterrupts()
				return nil
			},
			Cycles: 8,
		},
	})
}
