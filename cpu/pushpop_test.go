package cpu

import (
	"testing"

	op "github.com/paulloz/ohboi/cpu/opcodes"
)

func TestOpcodePUSH_AF(t *testing.T) {
	cpu := newTestCPU([]byte{op.PUSH_AF})
	cpu.SP.Set(0xff80)
	cpu.AF.Set(123)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.mem.ReadWord(0xff80) != 123 {
		t.Errorf("Expected address 0xff80 to contain 123, got %x", cpu.mem.ReadWord(0xff80))
	}
}

func TestOpcodePOP_AF(t *testing.T) {
	cpu := newTestCPU([]byte{op.POP_AF})
	cpu.mem.WriteWord(0x80, 0xff, 0xabcd)
	cpu.SP.Set(0xff80)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.AF.Get() != 0xabcd {
		t.Errorf("Expected AF to contain 0xabcd, got %x", cpu.AF.Get())
	}
}
