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
	})
}

func TestOpcodeIncA(t *testing.T) {
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
