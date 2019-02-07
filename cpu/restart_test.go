package cpu_test

import (
	"testing"

	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func TestOpcodeRST_N(t *testing.T) {
	newTestCPU(testScenario{
		bytecode: []byte{op.RST_30H},
		instr:    1,
		setup: func(cpu *cpu.CPU, mem *memory.Memory) {
			cpu.SP.Set(0xffff)
		},
		checks: []check{
			newRegister16Check("PC", cpu.RegisterPC, 0x30),
			newRegister16Check("SP", cpu.RegisterSP, 0xfffd),
			newMemoryWordCheck(0xfffd, 0x100),
		},
	})(t)
}
