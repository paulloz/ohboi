package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeSwapA(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.CB, op.SWAP_A},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.A.Set(0xf0)
		},
		checks: []check{
			newRegisterCheck("A", cpu.RegisterA, 0xf),
		},
	})(t)
}
