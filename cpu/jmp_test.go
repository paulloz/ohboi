package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeJP_NN(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.JP_NN, 0x80, 0xff},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.SP.Set(0xabcd)
		},
		checks: []check{
			newRegister16Check("PC", cpu.RegisterPC, 0xff80),
		},
	})(t)
}

func TestOpcodeJP_Z_NN(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.ADD_A_A, op.JP_Z_NN, 0x80, 0xff},
		instr:    2,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0)
		},
		checks: []check{
			newRegister16Check("PC", cpu.RegisterPC, 0xff80),
		},
	})(t)
}

func TestOpcodeJP_NZ_NN(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.ADD_A_A, op.JP_NZ_NN, 0x80, 0xff},
		instr:    2,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0)
		},
		checks: []check{
			newRegister16Check("PC", cpu.RegisterPC, 0x104),
		},
	})(t)
}

func TestOpcodeJP_HL(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.JP_HL},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(0xff80)
		},
		checks: []check{
			newRegister16Check("PC", cpu.RegisterPC, 0xff80),
		},
	})(t)
}

func TestOpcodeJR_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.JR_N, 0xb0},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.HL.Set(0xff80)
		},
		checks: []check{
			newRegister16Check("PC", cpu.RegisterPC, 0xb2),
		},
	})(t)
}

func TestOpcodeJR_Z_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.ADD_A_A, op.JR_Z_N, 0x7},
		instr:    2,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0)
		},
		checks: []check{
			newRegister16Check("PC", cpu.RegisterPC, 0x10a),
		},
	})(t)
}
