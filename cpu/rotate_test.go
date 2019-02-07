package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeRL_A(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.CB, op.RL_A},
		instr:    1,
		cycles:   8,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0x55)
			cpu.SetCFlag(true)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 0xab),
			zeroFlagResetCheck{},
			carryFlagResetCheck{},
		},
	})(t)
}

func TestOpcodeRL_HL(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.CB, op.RL_HL},
		instr:    1,
		cycles:   16,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(memory.InternalRAMAddr)
			mem.Write(memory.InternalRAMAddr, 0x80)
			cpu.SetCFlag(false)
		},
		checks: []check{
			newMemoryCheck(memory.InternalRAMAddr, 0),
			carryFlagSetCheck{},
			zeroFlagSetCheck{},
		},
	})(t)
}

func TestOpcodeRLCA(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.RLCA},
		instr:    1,
		cycles:   4,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0x55)
			cpu.SetCFlag(true)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 0xaa),
			zeroFlagResetCheck{},
			carryFlagResetCheck{},
		},
	})(t)
}

func TestOpcodeRLC_HL(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.CB, op.RLC_HL},
		instr:    1,
		cycles:   16,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(memory.InternalRAMAddr)
			mem.Write(memory.InternalRAMAddr, 0x80)
			cpu.SetCFlag(false)
		},
		checks: []check{
			newMemoryCheck(memory.InternalRAMAddr, 1),
			carryFlagSetCheck{},
			zeroFlagResetCheck{},
		},
	})(t)
}
