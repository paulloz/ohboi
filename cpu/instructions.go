package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

type Instruction struct {
	Handler func(cpu *CPU, mem *memory.Memory) error
	Cycles  uint32
}

var InstructionSet map[uint8]Instruction
var ExtInstructionSet map[uint8]Instruction

func RegisterInstructions(instructions map[uint8]Instruction) {
	if InstructionSet == nil {
		InstructionSet = make(map[uint8]Instruction)
	}

	for opcode, instruction := range instructions {
		InstructionSet[opcode] = instruction
	}
}

func RegisterExtInstructions(instructions map[uint8]Instruction) {
	if ExtInstructionSet == nil {
		ExtInstructionSet = make(map[uint8]Instruction)
	}

	for opcode, instruction := range instructions {
		ExtInstructionSet[opcode] = instruction
	}
}

var NoopInstruction = Instruction{
	Handler: func(cpu *CPU, mem *memory.Memory) error {
		return nil
	},
	Cycles: 4,
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.NOOP: NoopInstruction,
		op.DI: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.DisableInterrupts()
				return nil
			},
			Cycles: 4,
		},
		op.EI: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.EnableInterrupts()
				return nil
			},
			Cycles: 4,
		},
		op.CCF: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.SetNFlag(false)
				cpu.SetHFlag(false)
				cpu.SetCFlag(!cpu.GetCFlag())

				return nil
			},
			Cycles: 4,
		},
		op.SCF: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.SetNFlag(false)
				cpu.SetHFlag(false)
				cpu.SetCFlag(true)

				return nil
			},
			Cycles: 4,
		},
		op.CPL: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.A.Set(^cpu.A.Get())

				cpu.SetNFlag(true)
				cpu.SetHFlag(true)

				return nil
			},
			Cycles: 4,
		},
	})
}
