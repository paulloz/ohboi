package cpu

import (
	"fmt"

	"github.com/paulloz/ohboi/memory"
)

// CPU describes the Gameboy processor
type CPU struct {
	AF Register
	BC Register
	DE Register
	HL Register

	SP Register
	PC uint16

	mem *memory.Memory
}

func (cpu *CPU) Dump() {
	fmt.Printf("PC:%X\n", cpu.PC)
	fmt.Printf("A: %X, F: %X\n", cpu.AF.Hi(), cpu.AF.Lo())
	fmt.Printf("B: %X, C: %X\n", cpu.BC.Hi(), cpu.BC.Lo())
	fmt.Printf("D: %X, E: %X\n", cpu.DE.Hi(), cpu.DE.Lo())
}

func (cpu *CPU) FetchByte() uint8 {
	opCode := cpu.mem.Read(cpu.AdvancePC())
	return opCode
}

// ExecuteOpCode ...
// TODO: Maybe an array of func would be better?
// TODO: There's probably a better way to handle CPU cycles count
func (cpu *CPU) ExecuteOpCode() (uint, error) {
	opcode := cpu.FetchByte()

	instruction, ok := InstructionSet[opcode]
	if !ok {
		return 0, fmt.Errorf("opcode %X not implemented", opcode)
	}

	return instruction.Cycles, instruction.Handler(cpu, cpu.mem)
}

// AdvancePC returns PC value and increments it
func (cpu *CPU) AdvancePC() uint16 {
	n := cpu.PC
	cpu.PC++
	return n
}

// NewCPU ...
func NewCPU(mem *memory.Memory) *CPU {
	return &CPU{
		PC:  0x0100,
		AF:  NewRegister(0x01b0),
		BC:  NewRegister(0x01b0),
		DE:  NewRegister(0x01b0),
		HL:  NewRegister(0x01b0),
		SP:  NewRegister(0xfffe),
		mem: mem,
	}
}
