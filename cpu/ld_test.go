package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeLD_B_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_B_N, 123},
		instr:    1,
		checks: []check{
			newRegisterCheck("B", cpu.RegisterB, 123),
		},
	})(t)
}

func TestOpcodeLD_C_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_C_N, 123},
		instr:    1,
		checks: []check{
			newRegisterCheck("C", cpu.RegisterC, 123),
		},
	})(t)
}

func TestOpcodeLD_D_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_D_N, 123},
		instr:    1,
		checks: []check{
			newRegisterCheck("D", cpu.RegisterD, 123),
		},
	})(t)
}

func TestOpcodeLD_E_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_E_N, 123},
		instr:    1,
		checks: []check{
			newRegisterCheck("E", cpu.RegisterE, 123),
		},
	})(t)
}

func TestOpcodeLD_H_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_H_N, 123},
		instr:    1,
		checks: []check{
			newRegisterCheck("H", cpu.RegisterH, 123),
		},
	})(t)
}

func TestOpcodeLD_L_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_L_N, 123},
		instr:    1,
		checks: []check{
			newRegisterCheck("L", cpu.RegisterL, 123),
		},
	})(t)
}

func TestOpcodeLD_HL_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_HL_N, 123},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(memory.InternalRAM2Addr)
		},
		checks: []check{
			newMemoryCheck(memory.InternalRAM2Addr, 123),
		},
	})(t)
}

func TestOpcodeLD_A_NN(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_A_NN, uint8(memory.InternalRAM2Addr & 0xff), uint8(memory.InternalRAM2Addr >> 8)},
		instr:    1,
		cycles:   16,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			mem.Write(memory.InternalRAM2Addr, 123)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 123),
		},
	})(t)
}

func TestOpcodeLD_NN_A(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_NN_A, uint8(memory.InternalRAM2Addr & 0xff), uint8(memory.InternalRAM2Addr >> 8)},
		instr:    1,
		cycles:   16,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(123)
		},
		checks: []check{
			newMemoryCheck(memory.InternalRAM2Addr, 123),
		},
	})(t)
}

func TestOpcodeLD_A_CADDR(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_A_CADDR},
		instr:    1,
		cycles:   8,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.C.Set(0x80)
			mem.Write(memory.InternalRAM2Addr, 123)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 123),
		},
	})(t)
}

func TestOpcodeLD_A_HLD(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_A_HLD},
		instr:    1,
		cycles:   8,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(memory.InternalRAM2Addr)
			mem.Write(memory.InternalRAM2Addr, 123)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 123),
			newRegister16Check("HL", cpu.RegisterHL, memory.InternalRAM2Addr-1),
		},
	})(t)
}

func TestOpcodeLD_HLD_A(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_HLD_A},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(memory.InternalRAM2Addr)
			cpu.A.Set(123)
		},
		checks: []check{
			newMemoryCheck(memory.InternalRAM2Addr, 123),
			newRegister16Check("HL", cpu.RegisterHL, memory.InternalRAM2Addr-1),
		},
	})(t)
}

func TestOpcodeLDH_FF00N_A(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LDH_FF00N_A, 128},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(123)
		},
		checks: []check{
			newMemoryCheck(memory.InternalRAM2Addr, 123),
		},
	})(t)
}

func TestOpcodeLD_BC_NN(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_BC_NN, 123, 00},
		instr:    1,
		checks: []check{
			newRegister16Check("BC", cpu.RegisterBC, 123),
		},
	})(t)
}

func TestOpcodeLD_SP_HL(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_SP_HL},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(123)
		},
		checks: []check{
			newRegister16Check("SP", cpu.RegisterSP, 123),
		},
	})(t)
}

func TestOpcodeLD_HL_SP_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_HL_SP_N, 128},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.SP.Set(0xff00)
		},
		checks: []check{
			newRegister16Check("HL", cpu.RegisterHL, 0xfe80),
		},
	})(t)
}

func TestOpcodeLD_NN_SP(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.LD_NN_SP, 0x80, 0xff},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.SP.Set(0xabcd)
		},
		checks: []check{
			newMemoryWordCheck(0xff80, 0xabcd),
		},
	})(t)
}
