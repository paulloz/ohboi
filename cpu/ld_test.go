package cpu

import (
	"testing"

	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeLD_B_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_B_N, 123})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.BC.Hi() != 123 {
		t.Errorf("Expected A to be 123, got %d", cpu.BC.Hi())
	}
}

func TestOpcodeLD_C_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_C_N, 123})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.BC.Lo() != 123 {
		t.Errorf("Expected A to be 123, got %d", cpu.BC.Lo())
	}
}

func TestOpcodeLD_D_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_D_N, 123})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.DE.Hi() != 123 {
		t.Errorf("Expected A to be 123, got %d", cpu.DE.Hi())
	}
}

func TestOpcodeLD_E_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_E_N, 123})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.DE.Lo() != 123 {
		t.Errorf("Expected A to be 123, got %d", cpu.DE.Lo())
	}
}

func TestOpcodeLD_H_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_H_N, 123})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.HL.Hi() != 123 {
		t.Errorf("Expected A to be 123, got %d", cpu.HL.Lo())
	}
}

func TestOpcodeLD_L_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_L_N, 123})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.HL.Lo() != 123 {
		t.Errorf("Expected A to be 123, got %d", cpu.HL.Lo())
	}
}

func TestOpcodeLD_HL_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_HL_N, 123})
	cpu.HL.Set(memory.InternalRAM2Addr)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.mem.Read(memory.InternalRAM2Addr) != 123 {
		t.Errorf("Expected byte 51 to be 123, got %d", cpu.mem.Read(50))
	}
}
