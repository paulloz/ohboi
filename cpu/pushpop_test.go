package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodePUSH_AF(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.PUSH_AF},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.SP.Set(0xffff)
			cpu.AF.Set(123)
		},
		checks: []check{
			newMemoryWordCheck(0xfffd, 123),
		},
	})(t)
}

func TestOpcodePOP_AF(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.POP_AF},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			mem.WriteWord(0xff80, 0xabcd)
			cpu.SP.Set(0xff80)
		},
		checks: []check{
			newRegister16Check("AF", cpu.RegisterAF, 0xabcd),
		},
	})(t)
}
