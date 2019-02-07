package cpu

import (
	"testing"

	op "github.com/paulloz/ohboi/cpu/opcodes"
)

func TestOpcodeADD_A_A(t *testing.T) {
	cpu := newTestCPU([]byte{op.ADD_A_A})
	cpu.A.Set(10)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.A.Get() != 20 {
		t.Errorf("Expected A to be 20, got %d", cpu.A.Get())
	}
}

func TestOpcodeADD_A_CarryFlag(t *testing.T) {
	cpu := newTestCPU([]byte{op.ADD_A_A})
	cpu.A.Set(200)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.F.Get()&CarryFlag == 0 {
		t.Errorf("Expected carry to be set")
	}
}

func TestOpcodeADD_A_ZFlag(t *testing.T) {
	cpu := newTestCPU([]byte{op.ADD_A_A})
	cpu.A.Set(0)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.F.Get()&ZFlag == 0 {
		t.Errorf("Expected zero flag to be set")
	}
}
