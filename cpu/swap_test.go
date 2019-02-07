package cpu

import (
	"testing"

	op "github.com/paulloz/ohboi/cpu/opcodes"
)

func TestOpcodeSwapA(t *testing.T) {
	cpu := newTestCPU([]byte{op.CB, op.SWAP_A})
	cpu.A.Set(0xf0)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.A.Get() != 0xf {
		t.Errorf("Expected A to be 0xf, got %d", cpu.A.Get())
	}
}
