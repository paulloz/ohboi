package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeCALL_NN(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.CALL_NN, 0xbb, 0xaa},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0)
			cpu.SP.Set(0xffff)
		},
		checks: []check{
			newPCCheck(0xaabb),
			newRegister16Check("SP", cpu.RegisterSP, 0xfffd),
			newMemoryWordCheck(0xfffd, 0x0103),
		},
	})(t)
}

func TestOpcodeCALL_Z_NN(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.ADD_A_A, op.CALL_Z_NN, 0xbb, 0xaa},
		instr:    2,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0)
			cpu.SP.Set(0xffff)
		},
		checks: []check{
			newPCCheck(0xaabb),
			newRegister16Check("SP", cpu.RegisterSP, 0xfffd),
			newMemoryWordCheck(0xfffd, 0x0104),
		},
	})(t)
}
