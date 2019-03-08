package cpu

import (
	"errors"
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
	resetDIV  bool

	interruptsMasterEnable    bool
	interruptsMasterEnabling  bool
	interruptsMasterDisabling bool

	MicroInstructions []MicroInstruction

	mem *memory.Memory
	io  *io.IO

	dmaCycles         uint16
	dmaBaseSrcAddress uint16
	resetDMA          bool
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
	if cpu.IsDMA() && cpu.dmaCycles > 0 && cpu.PC >= memory.OAMAddr && cpu.PC < 0xfea0 {
		cpu.AdvancePC()
		return 0xff
	}
	return cpu.mem.Read(cpu.AdvancePC())
}

func (cpu *CPU) FetchWord() uint16 {
	defer cpu.AdvancePC() // because differs are cool
	return cpu.mem.ReadWord(cpu.AdvancePC())
}

var ClearPipeline = errors.New("ClearPipeline")

func (cpu *CPU) execMicro(microInstruction MicroInstruction) error {
	err := microInstruction(cpu, cpu.mem)
	if err == ClearPipeline {
		cpu.MicroInstructions = cpu.MicroInstructions[:0]
	}
	return nil
}

func (cpu *CPU) NextCycle() (uint32, error) {
	if len(cpu.MicroInstructions) > 0 {
		microInstruction := cpu.MicroInstructions[0]
		cpu.MicroInstructions = cpu.MicroInstructions[1:]
		return 4, cpu.execMicro(microInstruction)
	} else {
		return cpu.ExecuteOpCode()
	}
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

	if len(instruction.MicroInstructions) > 0 {
		cpu.MicroInstructions = instruction.MicroInstructions[1:]
		return 4, cpu.execMicro(instruction.MicroInstructions[0])
	}

	return instruction.Cycles, instruction.Handler(cpu, cpu.mem)
}

func (cpu *CPU) Jump(nn uint16) {
	cpu.PC = nn
}

// AdvancePC returns PC value and increments it
func (cpu *CPU) AdvancePC() uint16 {
	n := cpu.PC
	cpu.PC++
	return n
}

func (cpu *CPU) Push(v uint16) {
	cpu.SP.hilo -= 2
	cpu.WriteWord(cpu.SP.Get(), v)
}

func (cpu *CPU) PushByte(v uint8) {
	cpu.SP.hilo -= 1
	cpu.Write(cpu.SP.Get(), v)
}

func (cpu *CPU) Pop() uint16 {
	value := cpu.mem.ReadWord(cpu.SP.Get())
	cpu.SP.hilo += 2
	return value
}

func (cpu *CPU) PopByte() uint8 {
	value := cpu.Read(cpu.SP.Get())
	cpu.SP.hilo += 1
	return value
}

func (cpu *CPU) readDIV() uint8 {
	return cpu.div
}

func (cpu *CPU) writeDIV(val uint8) {
	cpu.resetDIV = true
	cpu.io.Write(io.TIMA, cpu.io.Read(io.TMA))
}

func (cpu *CPU) UpdateDIV(cycles uint32) {
	cpu.divCycles += cycles

	for cpu.divCycles >= 256 {
		if cpu.resetDIV {
			cpu.io.Write(io.TIMA, cpu.io.Read(io.TIMA)+1)
		}
		cpu.divCycles -= 256
		cpu.div++
	}

	if cpu.resetDIV {
		cpu.div = 0
		cpu.divCycles = 4
		cpu.resetDIV = false
	}
}

func (cpu *CPU) RequestInterrupt(interrupt uint8) {
	if interrupt >= I_VBLANK && interrupt <= I_JOYPAD {
		cpu.io.SetBit(io.IF, interrupt)
	}
}

func (cpu *CPU) Write(addr uint16, value uint8) {
	if !cpu.canAccess(addr) {
		return
	}
	cpu.mem.Write(addr, value)
}

func (cpu *CPU) canAccess(addr uint16) bool {
	if cpu.IsDMA() && cpu.dmaCycles > 0 && addr >= memory.OAMAddr && addr < 0xfea0 {
		return false
	}
	return true
}

func (cpu *CPU) Read(addr uint16) uint8 {
	if !cpu.canAccess(addr) {
		return 0xff
	}
	return cpu.mem.Read(addr)
}

func (cpu *CPU) ReadWord(addr uint16) uint16 {
	if !cpu.canAccess(addr) {
		return 0xffff
	}
	return cpu.mem.ReadWord(addr)
}

func (cpu *CPU) WriteWord(addr, value uint16) {
	if !cpu.canAccess(addr) {
		return
	}
	cpu.mem.Write(addr, uint8(value&0xff))
	cpu.mem.Write(addr+1, uint8(value>>8))
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

func (cpu *CPU) IsDMA() bool {
	return cpu.dmaCycles < 0xa1
}

func (cpu *CPU) UpdateDMA(cycles uint32) {
	for cycles >= 4 && cpu.IsDMA() {
		if cpu.resetDMA {
			cpu.dmaCycles = 0
			cpu.resetDMA = false
			cycles -= 4
			continue
		}

		if cpu.dmaCycles > 0 && cpu.dmaCycles <= 0xa0 {
			srcAddress := (uint16(cpu.dmaBaseSrcAddress) << 8) + cpu.dmaCycles - 1
			dstAddress := 0xfe00 + cpu.dmaCycles - 1
			cpu.mem.Write(dstAddress, cpu.mem.Read(srcAddress))
		}

		cpu.dmaCycles++
		cycles -= 4
	}
}

func (cpu *CPU) DMATransfert(baseSrcAddress uint8) {
	if baseSrcAddress >= 0xe0 {
		baseSrcAddress = baseSrcAddress - 0xe0 + 0xc0
	}
	cpu.dmaBaseSrcAddress = uint16(baseSrcAddress)
	cpu.dmaCycles = 0
	cpu.resetDMA = true
}

func (cpu *CPU) ReadDMA() uint8 {
	return uint8(cpu.dmaBaseSrcAddress)
}

func (cpu *CPU) DMACycles() uint16 {
	return cpu.dmaCycles
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

		dmaCycles: 0xa3,

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
