package cpu

import (
	"fmt"

	op "github.com/paulloz/ohboi/cpu/opcodes"
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

func (cpu *CPU) Dump() string {
	return fmt.Sprintf("PC:0x%02x\n", cpu.PC) +
		fmt.Sprintf("\n") +
		fmt.Sprintf("A: 0x%02x, F: 0x%02x\n", cpu.AF.Hi(), cpu.AF.Lo()) +
		fmt.Sprintf("B: 0x%02x, C: 0x%02x\n", cpu.BC.Hi(), cpu.BC.Lo()) +
		fmt.Sprintf("D: 0x%02x, E: 0x%02x\n", cpu.DE.Hi(), cpu.DE.Lo()) +
		fmt.Sprintf("\n") +
		fmt.Sprintf("HL: 0x%04x\n", cpu.HL.Get()) +
		fmt.Sprintf("\n") +
		fmt.Sprintf("Z: %5t, H: %5t\n", cpu.GetZFlag(), cpu.GetHFlag()) +
		fmt.Sprintf("N: %5t, C: %5t\n", cpu.GetNFlag(), cpu.GetCFlag())
}

func (cpu *CPU) FetchByte() uint8 {
	return cpu.mem.Read(cpu.AdvancePC())
}

func (cpu *CPU) FetchWord() uint16 {
	defer cpu.AdvancePC() // because differs are cool
	return cpu.mem.ReadWord(cpu.AdvancePC())
}

func (cpu *CPU) ExecuteOpCode() (uint, error) {
	opcode := cpu.FetchByte()

	var instruction Instruction
	var ok bool

	if opcode == op.CB {
		opcode = cpu.FetchByte()
		instruction, ok = ExtInstructionSet[opcode]
		if !ok {
			return 0, fmt.Errorf("extended opcode %X not implemented", opcode)
		}
	} else {
		instruction, ok = InstructionSet[opcode]
		if !ok {
			return 0, fmt.Errorf("opcode %X not implemented", opcode)
		}
	}

	return instruction.Cycles, instruction.Handler(cpu, cpu.mem)
}

// AdvancePC returns PC value and increments it
func (cpu *CPU) AdvancePC() uint16 {
	n := cpu.PC
	cpu.PC++
	return n
}

func (cpu *CPU) Push(v uint16) {
	cpu.SP.hilo -= 2
	cpu.mem.WriteWord(cpu.SP.Get(), v)
}

func (cpu *CPU) Pop() uint16 {
	value := cpu.mem.ReadWord(cpu.SP.Get())
	cpu.SP.hilo += 2
	return value
}

func (cpu *CPU) EnableInterrupts() {
	// TODO
}

func (cpu *CPU) DisableInterrupts() {
	// TODO
}

func NewCPU(mem *memory.Memory) *CPU {
	cpu := &CPU{
		PC:  0x0,
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
