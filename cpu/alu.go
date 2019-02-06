package cpu

import (
	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func (cpu *CPU) And(out func(uint8), a uint8, b uint8) {
	result := b & a
	out(result)

	cpu.SetZFlag(result == 0)
	cpu.SetNFlag(false)
	cpu.SetHFlag(true)
	cpu.SetCFlag(false)
}

func newIncrementRegister(register GetterSetter) Instruction {
	cycles := uint(4)
	if register == AddressHL {
		cycles = 12
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			initial := register.Get(cpu)
			final := initial + 1
			register.Set(cpu, final)

			cpu.SetZFlag(final == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(bits.HalfCarryCheck(initial, 1))

			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterIntructions(map[uint8]Instruction{
		op.AND_A_E: {
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				cpu.And(cpu.AF.SetHi, cpu.DE.Lo(), cpu.AF.Hi())
				return nil
			},
			Cycles: 4,
		},

		op.INC_A:  newIncrementRegister(RegisterA),
		op.INC_B:  newIncrementRegister(RegisterB),
		op.INC_C:  newIncrementRegister(RegisterC),
		op.INC_D:  newIncrementRegister(RegisterD),
		op.INC_E:  newIncrementRegister(RegisterE),
		op.INC_H:  newIncrementRegister(RegisterH),
		op.INC_L:  newIncrementRegister(RegisterL),
		op.INC_HL: newIncrementRegister(AddressHL),
	})
}
