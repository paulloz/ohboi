package cpu

import (
	"testing"

	op "github.com/paulloz/ohboi/cpu/opcodes"
)

func TestOpcodeJP_NN(t *testing.T) {
	cpu := newTestCPU([]byte{op.JP_NN, 0x80, 0xff})
	cpu.SP.Set(0xabcd)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0xff80 {
		t.Errorf("Expected PC to contain 0xff00, got %x", cpu.PC)
	}
}

func TestOpcodeJP_Z_NN(t *testing.T) {
	cpu := newTestCPU([]byte{op.ADD_A_A, op.JP_Z_NN, 0x80, 0xff})
	cpu.A.Set(0)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	_, err = cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0xff80 {
		t.Errorf("Expected PC to contain 0xff80, got %x", cpu.PC)
	}
}

func TestOpcodeJP_NZ_NN(t *testing.T) {
	cpu := newTestCPU([]byte{op.ADD_A_A, op.JP_NZ_NN, 0x80, 0xff})
	cpu.A.Set(0)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	_, err = cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0x104 {
		t.Errorf("Expected PC to contain 0x104, got %X", cpu.PC)
	}
}

func TestOpcodeJP_HL(t *testing.T) {
	cpu := newTestCPU([]byte{op.JP_HL})
	cpu.HL.Set(0xff80)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0xff80 {
		t.Errorf("Expected PC to contain 0xff80, got %X", cpu.PC)
	}
}

func TestOpcodeJR_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.JR_N, 0xb0})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0xb2 {
		t.Errorf("Expected PC to contain 0xb2, got %X", cpu.PC)
	}
}

func TestOpcodeJR_Z_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.ADD_A_A, op.JR_Z_N, 0x7})
	cpu.A.Set(0)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	_, err = cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.PC != 0x10a {
		t.Errorf("Expected PC to contain 0x10a, got %X", cpu.PC)
	}
}
