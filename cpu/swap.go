package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newSwap(reg GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			lowerNibble := reg.Get(cpu) & 0x0f
			upperNibble := ((reg.Get(cpu) & 0xf0) >> 4) & 0x0f

			reg.Set(cpu, ((lowerNibble << 4) | upperNibble))

			cpu.SetZFlag(reg.Get(cpu) == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(false)

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
