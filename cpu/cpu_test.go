package cpu

import (
	"testing"

	"github.com/paulloz/ohboi/cartridge"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newTestCPU(bytecode []byte) *CPU {
	data := make([]byte, 256+len(bytecode))
	copy(data[256:], bytecode)

	rom := cartridge.NewROM(data)

	memory := memory.NewMemory()
	memory.LoadCartridge(&cartridge.Cartridge{
		MBC:   rom,
		Title: "TestROM",
	})

	return NewCPU(memory)
}

func TestOpcodeNoop(t *testing.T) {
	cpu := newTestCPU([]byte{op.NOOP})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}
}

func TestOpcodeIncA(t *testing.T) {
	cpu := newTestCPU([]byte{op.INC_A})
	previousA := cpu.AF.Hi()

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.AF.Hi() != previousA+1 {
		t.Errorf("Expected A to be %d, got %d", previousA, cpu.AF.Hi())
	}
}
