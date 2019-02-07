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

func TestOpcodeLD_A_NN(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_A_NN, uint8(memory.InternalRAM2Addr & 0xff), uint8(memory.InternalRAM2Addr >> 8)})
	cpu.mem.Write(memory.InternalRAM2Addr, 123)

	cycles, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.A.Get() != 123 {
		t.Errorf("Expected A to contain 123, got %d", cpu.mem.Read(50))
	}

	if cycles != 16 {
		t.Errorf("Expected LD_A_NN to take 16 cycles, got %d", cycles)
	}
}

func TestOpcodeLD_NN_A(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_NN_A, uint8(memory.InternalRAM2Addr & 0xff), uint8(memory.InternalRAM2Addr >> 8)})
	cpu.A.Set(123)

	cycles, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.mem.Read(memory.InternalRAM2Addr) != 123 {
		t.Errorf("Expected A to contain 123, got %d", cpu.mem.Read(memory.InternalRAM2Addr))
	}

	if cycles != 16 {
		t.Errorf("Expected LD_NN_A to take 16 cycles, got %d", cycles)
	}
}
func TestOpcodeLD_A_CADDR(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_A_CADDR})
	cpu.C.Set(0x80)
	cpu.mem.Write(memory.InternalRAM2Addr, 123)

	cycles, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.A.Get() != 123 {
		t.Errorf("Expected A to contain 123, got %d", cpu.mem.Read(memory.IOPortsAddr))
	}

	if cycles != 8 {
		t.Errorf("Expected LD_A_CADDR to take 8 cycles, got %d", cycles)
	}
}

func TestOpcodeLD_A_HLD(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_A_HLD})
	cpu.HL.Set(memory.InternalRAM2Addr)
	cpu.mem.Write(memory.InternalRAM2Addr, 123)

	cycles, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.A.Get() != 123 {
		t.Errorf("Expected A to contain 123, got %d", cpu.mem.Read(memory.IOPortsAddr))
	}

	if cpu.HL.Get() != memory.InternalRAM2Addr-1 {
		t.Errorf("Expected A to contain InternalRAM2Addr, got %d", cpu.HL.Get())
	}

	if cycles != 8 {
		t.Errorf("Expected LD_A_HLD to take 8 cycles, got %d", cycles)
	}
}

func TestOpcodeLD_HLD_A(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_HLD_A})
	cpu.HL.Set(memory.InternalRAM2Addr)
	cpu.A.Set(123)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.mem.Read(memory.InternalRAM2Addr) != 123 {
		t.Errorf("Expected InternalRAM2Addr to contain 123, got %d", cpu.mem.Read(memory.InternalRAM2Addr))
	}

	if cpu.HL.Get() != memory.InternalRAM2Addr-1 {
		t.Errorf("Expected A to contain InternalRAM2Addr, got %d", cpu.HL.Get())
	}
}

func TestOpcodeLDH_FF00N_A(t *testing.T) {
	cpu := newTestCPU([]byte{op.LDH_FF00N_A, 128})
	cpu.A.Set(123)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.mem.Read(memory.InternalRAM2Addr) != 123 {
		t.Errorf("Expected InternalRAM2Addr to contain 123, got %d", cpu.mem.Read(memory.InternalRAM2Addr))
	}

}

func TestOpcodeLD_BC_NN(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_BC_NN, 123, 00})

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.BC.Get() != 123 {
		t.Errorf("Expected BC to contain 123, got %d", cpu.BC.Get())
	}

}

func TestOpcodeLD_SP_HL(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_SP_HL})
	cpu.HL.Set(123)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.SP.Get() != 123 {
		t.Errorf("Expected SP to contain 123, got %d", cpu.SP.Get())
	}
}

func TestOpcodeLD_HL_SP_N(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_HL_SP_N, 128})
	cpu.SP.Set(0xff00)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.HL.Get() != 0xff80 {
		t.Errorf("Expected SP to contain 0xff80, got %x", cpu.HL.Get())
	}
}

func TestOpcodeLD_NN_SP(t *testing.T) {
	cpu := newTestCPU([]byte{op.LD_NN_SP, 0x80, 0xff})
	cpu.SP.Set(0xabcd)

	_, err := cpu.ExecuteOpCode()
	if err != nil {
		t.Error(err)
	}

	if cpu.mem.ReadWord(0xff80) != 0xabcd {
		t.Errorf("Expected address 0xff80 to contain 0xabcd, got %x", cpu.mem.ReadWord(0xff80))
	}
}
