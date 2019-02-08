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
	Handler: NoopHandler,
	Cycles:  4,
}

func NoopHandler(cpu *CPU, mem *memory.Memory) error {
	return nil
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.NOOP: NoopInstruction,
	})
}
