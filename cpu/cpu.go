package cpu

import (
	"fmt"

	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/io"
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

	isHalted bool

	divCycles uint32
	div       uint8

	interruptsMasterEnable    bool
	interruptsMasterEnabling  bool
	interruptsMasterDisabling bool

	mem *memory.Memory
	io  *io.IO
}

func (cpu *CPU) Dump() string {
	return fmt.Sprintf("PC:0x%04x\n", cpu.PC) +
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

func (cpu *CPU) ExecuteOpCode() (uint32, error) {
	if cpu.isHalted {
		return InstructionSet[op.NOOP].Cycles, nil
	}

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

func (cpu *CPU) readDIV() uint8 {
	return cpu.div
}

func (cpu *CPU) writeDIV(val uint8) {
	cpu.div = 0
	cpu.divCycles = 0
	cpu.mem.Write((0xff00 | io.TIMA), cpu.mem.Read(io.TMA))
}

func (cpu *CPU) UpdateDIV(cycles uint32) {
	cpu.divCycles += cycles

	for cpu.divCycles >= 256 {
		cpu.divCycles -= 256
		cpu.div++
	}
}

func (cpu *CPU) RequestInterrupt(interrupt uint8) {
	if interrupt >= I_VBLANK && interrupt <= I_JOYPAD {
		cpu.io.SetBit(io.IF, interrupt)
	}
}

func (cpu *CPU) EnableInterrupts() {
	cpu.interruptsMasterEnabling = true
}

func (cpu *CPU) DisableInterrupts() {
	cpu.interruptsMasterDisabling = true
}

func (cpu *CPU) ManageInterrupts() uint32 {
	defer func() {
		if cpu.interruptsMasterDisabling {
			cpu.interruptsMasterEnable = false
			cpu.interruptsMasterDisabling = false
		}
	}()

	if cpu.interruptsMasterEnabling {
		cpu.interruptsMasterEnabling = false
		if !cpu.interruptsMasterEnable {
			cpu.interruptsMasterEnable = true
			return 0
		}
	}

	if !cpu.interruptsMasterEnable && !cpu.isHalted {
		return 0
	}

	interruptEnable := cpu.io.Read(io.IE)
	interruptFlag := cpu.io.Read(io.IF)

	// Interrupts priority is from least to most significant bits
	for b := uint8(0); b <= 4; b++ {
		// If interrupt was requested
		if bits.Test(b, interruptFlag) {
			// And this interrupt is enabled
			if bits.Test(b, interruptEnable) {
				if cpu.isHalted {
					cpu.isHalted = false
					return 0
				}

				// Service the interrupt:
				cpu.interruptsMasterEnable = false // Clear IME

				// Reset bit in IF
				cpu.io.ResetBit(io.IF, b)

				// Push PC to Stack (8 cycles)
				cpu.Push(cpu.PC)
				// Set PC to interrupt handler address (4 cycles)
				cpu.PC = interrupts[b]
				// Execute 2 NOP (8 cycles)
				return 20
			}
		}
	}

	return 0
}

func (cpu *CPU) DMATransfert(baseSrcAddress uint8) {
	for i := uint16(0); i < 0xa0; i++ {
		srcAddress := uint16(baseSrcAddress) << 8 + i
		dstAddress := 0xfe00 + i
		cpu.mem.Write(dstAddress, cpu.mem.Read(srcAddress))
	}
}

func (cpu *CPU) ReadDMA() uint8 {
	return cpu.mem.LastWrittenValue()
}

func NewCPU(mem *memory.Memory, io_ *io.IO) *CPU {
	cpu := &CPU{
		PC: 0x0,
		AF: NewRegister(0x01b0),
		BC: NewRegister(0x01b0),
		DE: NewRegister(0x01b0),
		HL: NewRegister(0x01b0),
		SP: NewRegister(0xfffe),

		isHalted: false,

		divCycles: 0,
		div:       0,

		interruptsMasterEnable:    true,
		interruptsMasterEnabling:  false,
		interruptsMasterDisabling: false,

		mem: mem,
		io:  io_,
	}

	cpu.A = PseudoRegisterHigh{hwRegister: &cpu.AF}
	cpu.F = PseudoRegisterLow{hwRegister: &cpu.AF}
	cpu.B = PseudoRegisterHigh{hwRegister: &cpu.BC}
	cpu.C = PseudoRegisterLow{hwRegister: &cpu.BC}
	cpu.D = PseudoRegisterHigh{hwRegister: &cpu.DE}
	cpu.E = PseudoRegisterLow{hwRegister: &cpu.DE}
	cpu.H = PseudoRegisterHigh{hwRegister: &cpu.HL}
	cpu.L = PseudoRegisterLow{hwRegister: &cpu.HL}

	io_.MapRegister(io.DIV, cpu.readDIV, cpu.writeDIV)
	io_.MapRegister(io.DMA, cpu.ReadDMA, cpu.DMATransfert)

	return cpu
}
