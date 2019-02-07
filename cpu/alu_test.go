package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeADD_A_A(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.ADD_A_A},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(10)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 20),
		},
	})(t)
}

func TestOpcodeADD_A_CarryFlag(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.ADD_A_A},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(200)
		},
		checks: []check{
			carryFlagSetCheck{},
		},
	})(t)
}

func TestOpcodeADD_A_ZFlag(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.ADD_A_A},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0)
		},
		checks: []check{
			zeroFlagSetCheck{},
		},
	})(t)
}

func TestOpcodeINC_A(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.INC_A},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(11)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 12),
		},
	})(t)
}

func TestOpCodeXOR_A(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.XOR_B},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0x55)
			cpu.B.Set(0xaa)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 0xff),
			zeroFlagResetCheck{},
			carryFlagResetCheck{},
		},
	})(t)
}
func TestOpCodeXOR_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.XOR_N, 0x55},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0x55)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 0),
			zeroFlagSetCheck{},
			carryFlagResetCheck{},
		},
	})(t)
}
func TestOpCodeXOR_HL(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.XOR_HL},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0x0f)
			cpu.HL.Set(memory.InternalRAMAddr)
			mem.Write(memory.InternalRAMAddr, 0x1f)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 0x10),
			zeroFlagResetCheck{},
			carryFlagResetCheck{},
		},
	})(t)
}

func TestOpCodeAND_A_B(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.AND_A_B},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0xab)
			cpu.B.Set(0xf0)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 0xa0),
		},
	})(t)
}
