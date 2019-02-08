package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cartridge"
	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/memory"
)

type testScenario struct {
	bytecode []byte
	instr    uint
	cycles   uint32
	setup    func(*cpu.CPU, *memory.Memory)
	checks   []check
}

func newTestCPU(scenario testScenario) func(t *testing.T) {
	return func(t *testing.T) {
		data := make([]byte, 256+len(scenario.bytecode))
		copy(data[256:], scenario.bytecode)

		rom := cartridge.NewROM(data)

		io := io.NewIO()

		memory := memory.NewMemory(io)
		memory.LoadCartridge(&cartridge.Cartridge{MBC: rom})

		cpu := cpu.NewCPU(memory, io)
		cpu.PC = 0x100

		if scenario.setup != nil {
			scenario.setup(cpu, memory)
		}

		cycles := uint32(0)
		for instr := uint(0); instr < scenario.instr; instr++ {
			c, err := cpu.ExecuteOpCode()
			if err != nil {
				t.Error(err)
			}
			cycles += c
		}

		if scenario.cycles != 0 {
			if cycles != scenario.cycles {
				t.Errorf("Expected to take %d cycles, got %d", scenario.cycles, cycles)
			}
		}

		for _, check := range scenario.checks {
			check.Check(t, cpu, memory)
		}
	}
}

type check interface {
	Check(*testing.T, *cpu.CPU, *memory.Memory)
}

type memoryCheck struct {
	address uint16
	value   uint8
}

func (c memoryCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	b := mem.Read(c.address)
	if b != c.value {
		t.Errorf("Expected memory 0x%x to contain 0x%x. got 0x%x", c.address, c.value, b)
	}
}

func newMemoryCheck(address uint16, value uint8) memoryCheck {
	return memoryCheck{address: address, value: value}
}

type memoryWordCheck struct {
	address uint16
	value   uint16
}

func (c memoryWordCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	word := mem.ReadWord(c.address)
	if word != c.value {
		t.Errorf("Expected memory 0x%x to contain 0x%x. got 0x%x", c.address, c.value, word)
	}
}

func newMemoryWordCheck(address uint16, value uint16) memoryWordCheck {
	return memoryWordCheck{address: address, value: value}
}

type pcCheck struct{ value uint16 }

func (c pcCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	if cpu.PC != c.value {
		t.Errorf("Expected PC to contain 0x%x, got 0x%x", c.value, cpu.PC)
	}
}

func newPCCheck(value uint16) pcCheck {
	return pcCheck{value: value}
}

type registerCheck struct {
	name  string
	g     cpu.Getter
	value uint8
}

func (c registerCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	if c.g.Get(cpu) != c.value {
		t.Errorf("Expected %s to contain 0x%x, got 0x%x", c.name, c.value, c.g.Get(cpu))
	}
}

func newRegisterCheck(name string, g cpu.Getter, value uint8) registerCheck {
	return registerCheck{value: value, g: g, name: name}
}

type register16Check struct {
	name  string
	g     cpu.Getter16
	value uint16
}

func (c register16Check) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	if c.g.Get(cpu) != c.value {
		t.Errorf("Expected %s to contain 0x%x, got 0x%x", c.name, c.value, c.g.Get(cpu))
	}
}

func newRegister16Check(name string, g cpu.Getter16, value uint16) register16Check {
	return register16Check{value: value, g: g, name: name}
}

type carryFlagSetCheck struct{}

func (c carryFlagSetCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	if !cpu.GetCFlag() {
		t.Errorf("Expected C flag to be set")
	}
}

type carryFlagResetCheck struct{}

func (c carryFlagResetCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	if cpu.GetCFlag() {
		t.Errorf("Expected C flag to be reset")
	}
}

type zeroFlagSetCheck struct{}

func (c zeroFlagSetCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	if !cpu.GetZFlag() {
		t.Errorf("Expected Z flag to be set")
	}
}

type zeroFlagResetCheck struct{}

func (c zeroFlagResetCheck) Check(t *testing.T, cpu *cpu.CPU, mem *memory.Memory) {
	if cpu.GetZFlag() {
		t.Errorf("Expected Z flag to be reset")
	}
}

func TestOpcodeNoop(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.NOOP},
		instr:    1,
	})(t)
}
