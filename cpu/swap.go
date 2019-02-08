package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newSwap(reg GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			reg.Set(cpu, ((reg.Get(cpu)&0xf)<<4)|((reg.Get(cpu)&0xf0)>>4))
			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterExtInstructions(map[uint8]Instruction{
		op.SWAP_A:  newSwap(RegisterA, 8),
		op.SWAP_B:  newSwap(RegisterB, 8),
		op.SWAP_C:  newSwap(RegisterC, 8),
		op.SWAP_D:  newSwap(RegisterD, 8),
		op.SWAP_E:  newSwap(RegisterE, 8),
		op.SWAP_H:  newSwap(RegisterH, 8),
		op.SWAP_L:  newSwap(RegisterL, 8),
		op.SWAP_HL: newSwap(AddressHL, 16),
	})
}
