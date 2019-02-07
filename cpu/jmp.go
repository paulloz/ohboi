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
				cpu.Jump(cpu.FetchWord())
				return nil
			},
			Cycles: 12,
		},

		op.JP_Z_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if cpu.GetZFlag() {
					cpu.Jump(addr)
				}
				return nil
			},
			Cycles: 12,
		},

		op.JP_NZ_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if !cpu.GetZFlag() {
					cpu.Jump(addr)
				}
				return nil
			},
			Cycles: 12,
		},

		op.JP_C_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if cpu.GetCFlag() {
					cpu.Jump(addr)
				}
				return nil
			},
			Cycles: 12,
		},

		op.JP_NC_NN: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				addr := cpu.FetchWord()
				if !cpu.GetCFlag() {
					cpu.Jump(addr)
				}
				return nil
			},
			Cycles: 12,
		},

		op.JP_HL: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Jump(cpu.HL.Get())
				return nil
			},
			Cycles: 12,
		},

		op.JR_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.Jump(uint16(int(cpu.PC) + int(int8(cpu.FetchByte()))))
				return nil
			},
			Cycles: 8,
		},

		op.JR_C_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				rel := int(int8(cpu.FetchByte()))
				if cpu.GetCFlag() {
					cpu.Jump(uint16(int(cpu.PC) + rel))
				}
				return nil
			},
			Cycles: 12,
		},

		op.JR_NC_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				rel := int(int8(cpu.FetchByte()))
				if !cpu.GetCFlag() {
					cpu.Jump(uint16(int(cpu.PC) + rel))
				}
				return nil
			},
			Cycles: 12,
		},

		op.JR_Z_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				rel := int(int8(cpu.FetchByte()))
				if cpu.GetZFlag() {
					cpu.Jump(uint16(int(cpu.PC) + rel))
				}
				return nil
			},
			Cycles: 12,
		},

		op.JR_NZ_N: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				rel := int(int8(cpu.FetchByte()))
				if !cpu.GetCFlag() {
					cpu.Jump(uint16(int(cpu.PC) + rel))
				}
				return nil
			},
			Cycles: 12,
		},
	})
}
