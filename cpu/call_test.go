package cpu

import (
	"testing"

	op "github.com/paulloz/ohboi/cpu/opcodes"
)

func TestOpcodeCALL_NN(t *testing.T) {
	cpu := newTestCPU([]byte{op.CALL_NN, 0xbb, 0xaa})
	cpu.A.Set(0)
	cpu.SP.Set(0xffff)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0xaabb {
		t.Errorf("Expected PC to contain 0xaabb, got %X", cpu.PC)
	}

	if cpu.SP.Get() != 0xfffd {
		t.Errorf("Expected SP to contain 0xfffd, got %X", cpu.SP.Get())
	}

	if cpu.mem.ReadWord(0xfffd) != 0x103 {
		t.Errorf("Expected memory 0xfffd to contain 0x103, got %X", cpu.mem.ReadWord(0xfffd))
	}
}

func TestOpcodeCALL_Z_NN(t *testing.T) {
	cpu := newTestCPU([]byte{op.ADD_A_A, op.CALL_Z_NN, 0xbb, 0xaa})
	cpu.A.Set(0)
	cpu.SP.Set(0xffff)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	_, err = cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0xaabb {
		t.Errorf("Expected PC to contain 0xaabb, got %X", cpu.PC)
	}

	if cpu.SP.Get() != 0xfffd {
		t.Errorf("Expected SP to contain 0xfffd, got %X", cpu.SP.Get())
	}

	if cpu.mem.ReadWord(0xfffd) != 0x104 {
		t.Errorf("Expected memory 0xfffd to contain 0x104, got %X", cpu.mem.ReadWord(0xfffd))
	}
}
