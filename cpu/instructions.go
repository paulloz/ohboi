package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

type MicroInstruction func(cpu *CPU, mem *memory.Memory) error

func (m MicroInstruction) Combine(m2 MicroInstruction) MicroInstruction {
	return func(cpu *CPU, mem *memory.Memory) error {
		if err := m(cpu, mem); err != nil {
			return err
		}
		return m2(cpu, mem)
	}
}

type Instruction struct {
	Handler           func(cpu *CPU, mem *memory.Memory) error
	Cycles            uint32
	MicroInstructions []MicroInstruction
}

func (i Instruction) Handle(cpu *CPU, mem *memory.Memory) error {
	return i.Handler(cpu, mem)
}

func NewInstruction(micros ...MicroInstruction) Instruction {
	return Instruction{
		MicroInstructions: micros,
	}
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

func DecodeInstruction(cpu *CPU, mem *memory.Memory) error {
	return nil
}

func NoopInstruction(cpu *CPU, mem *memory.Memory) error {
	return nil
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.NOOP: Instruction{
			Handler: NoopInstruction,
			Cycles:  4,
		},
		op.HALT: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.isHalted = true
				return nil
			},
			Cycles: 4,
		},
		op.STOP: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				return nil
			},
			Cycles: 4,
		},
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
		op.DAA: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				a := int16(cpu.A.Get())

				if !cpu.GetNFlag() {
					if cpu.GetHFlag() || ((a & 0x0f) > 9) {
						a += 6
					}
					if cpu.GetCFlag() || (a > 0x9f) {
						a += 0x60
					}
				} else {
					if cpu.GetHFlag() {
						a -= 6
						if !cpu.GetCFlag() {
							a &= 0xff
						}
					}
					if cpu.GetCFlag() {
						a -= 0x60
					}
				}

				cpu.A.Set(uint8(a & 0xff))

				cpu.SetZFlag(cpu.A.Get() == 0)
				if a&0x100 != 0 {
					cpu.SetCFlag(true)
				}
				cpu.SetHFlag(false)

				return nil
			},
			Cycles: 4,
		},
	})
}
