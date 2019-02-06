package cpu

import (
	"fmt"

	"github.com/paulloz/ohboi/memory"
)

type PseudoRegister interface {
	Get() uint8
	Set(value uint8)
}

// CPU describes the Gameboy processor
type CPU struct {
	AF Register
	BC Register
	DE Register
	HL Register

	A PseudoRegister
	B PseudoRegister
	C PseudoRegister
	D PseudoRegister
	E PseudoRegister
	F PseudoRegister
	H PseudoRegister
	L PseudoRegister

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

type PseudoRegisterHigh struct {
	hwRegister *Register
}

func (r PseudoRegisterHigh) Get() uint8 {
	return r.hwRegister.Hi()
}

func (r PseudoRegisterHigh) Set(v uint8) {
	r.hwRegister.SetHi(v)
}

type PseudoRegisterLow struct {
	hwRegister *Register
}

func (r PseudoRegisterLow) Get() uint8 {
	return r.hwRegister.Lo()
}

func (r PseudoRegisterLow) Set(v uint8) {
	r.hwRegister.SetLo(v)
}

func NewCPU(mem *memory.Memory) *CPU {
	cpu := &CPU{
		PC:  0x0100,
		AF:  NewRegister(0x01b0),
		BC:  NewRegister(0x01b0),
		DE:  NewRegister(0x01b0),
		HL:  NewRegister(0x01b0),
		SP:  NewRegister(0xfffe),
		mem: mem,
	}
	cpu.A = PseudoRegisterHigh{hwRegister: &cpu.AF}
	cpu.F = PseudoRegisterLow{hwRegister: &cpu.AF}
	cpu.B = PseudoRegisterHigh{hwRegister: &cpu.BC}
	cpu.C = PseudoRegisterLow{hwRegister: &cpu.BC}
	cpu.D = PseudoRegisterHigh{hwRegister: &cpu.DE}
	cpu.E = PseudoRegisterLow{hwRegister: &cpu.DE}
	cpu.H = PseudoRegisterHigh{hwRegister: &cpu.HL}
	cpu.L = PseudoRegisterLow{hwRegister: &cpu.HL}
	return cpu
}
