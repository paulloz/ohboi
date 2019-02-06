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

func newCompareA(src Getter) Instruction {
	cycles := uint(4)
	if src == AddressHL || src == Immediate {
		cycles = 8
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			a := cpu.A.Get()
			n := src.Get(cpu)
			result := a - n

			cpu.SetZFlag(result == 0)
			cpu.SetNFlag(true)
			cpu.SetHFlag((a & 0x0f) < (n & 0x0f))
			cpu.SetCFlag(a < n)

			return nil
		},
		Cycles: cycles,
	}
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

func newXorA(src Getter) Instruction {
	cycles := uint(4)
	if src == AddressHL || src == Immediate {
		cycles = 8
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			final := src.Get(cpu) ^ cpu.A.Get()
			cpu.A.Set(final)

			cpu.SetZFlag(final == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(false)

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

		op.CP_A:  newCompareA(RegisterA),
		op.CP_B:  newCompareA(RegisterB),
		op.CP_C:  newCompareA(RegisterC),
		op.CP_D:  newCompareA(RegisterD),
		op.CP_E:  newCompareA(RegisterE),
		op.CP_H:  newCompareA(RegisterH),
		op.CP_L:  newCompareA(RegisterL),
		op.CP_HL: newCompareA(AddressHL),
		op.CP_N:  newCompareA(Immediate),

		op.INC_A:  newIncrementRegister(RegisterA),
		op.INC_B:  newIncrementRegister(RegisterB),
		op.INC_C:  newIncrementRegister(RegisterC),
		op.INC_D:  newIncrementRegister(RegisterD),
		op.INC_E:  newIncrementRegister(RegisterE),
		op.INC_H:  newIncrementRegister(RegisterH),
		op.INC_L:  newIncrementRegister(RegisterL),
		op.INC_HL: newIncrementRegister(AddressHL),

		op.XOR_A:  newXorA(RegisterA),
		op.XOR_B:  newXorA(RegisterB),
		op.XOR_C:  newXorA(RegisterC),
		op.XOR_D:  newXorA(RegisterD),
		op.XOR_E:  newXorA(RegisterE),
		op.XOR_H:  newXorA(RegisterH),
		op.XOR_L:  newXorA(RegisterL),
		op.XOR_HL: newXorA(AddressHL),
		op.XOR_N:  newXorA(Immediate),
	})
}
