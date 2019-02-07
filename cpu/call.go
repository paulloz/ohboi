package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func (cpu *CPU) Call(nn uint16) {
	cpu.Push(cpu.PC)
	cpu.PC = nn
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.CALL_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Call(cpu.FetchWord())
				return nil
			},
			Cycles: 12,
		},

		op.CALL_C_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if cpu.GetCFlag() {
					cpu.Call(addr)
				}
				return nil
			},
			Cycles: 12,
		},

		op.CALL_NC_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if !cpu.GetCFlag() {
					cpu.Call(addr)
				}
				return nil
			},
			Cycles: 12,
		},

		op.CALL_Z_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if cpu.GetZFlag() {
					cpu.Call(addr)
				}
				return nil
			},
			Cycles: 12,
		},

		op.CALL_NZ_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if !cpu.GetZFlag() {
					cpu.Call(addr)
				}
				return nil
			},
			Cycles: 12,
		},
	})
}
