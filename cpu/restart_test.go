package cpu

import (
	"testing"

	op "github.com/paulloz/ohboi/cpu/opcodes"
)

func TestOpcodeRST_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.RST_30H})
	cpu.SP.Set(0xffff)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0x30 {
		t.Errorf("Expected PC to contain 0x30, got %X", cpu.PC)
	}

	if cpu.SP.Get() != 0xfffd {
		t.Errorf("Expected SP to contain 0xfffd, got %X", cpu.SP.Get())
	}

	if cpu.mem.ReadWord(0xfffd) != 0x100 {
		t.Errorf("Expected memory 0xfffd to contain 0x100, got %X", cpu.mem.ReadWord(0xfffd))
	}
}
