package cpu

import (
	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newAdd(src Getter, cycles uint) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			sum := uint16(cpu.A.Get()) + uint16(src.Get(cpu))
			cpu.A.Set(uint8(sum))

			cpu.SetZFlag(sum == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(sum&0x8 != 0)
			cpu.SetCFlag(sum >= 256)
			return nil
		},
		Cycles: cycles,
	}
}

func newAddC(src Getter, cycles uint) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			var carry uint8
			if cpu.F.Get()&CarryFlag != 0 {
				carry = 1
			}
			sum := uint16(cpu.A.Get()) + uint16(src.Get(cpu)+carry)
			cpu.A.Set(uint8(sum))

			cpu.SetZFlag(sum == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(sum&0x8 != 0)
			cpu.SetCFlag(sum >= 256)
			return nil
		},
		Cycles: cycles,
	}
}

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
		op.ADD_A_A:  newAdd(RegisterA, 4),
		op.ADD_A_B:  newAdd(RegisterB, 4),
		op.ADD_A_C:  newAdd(RegisterC, 4),
		op.ADD_A_D:  newAdd(RegisterD, 4),
		op.ADD_A_E:  newAdd(RegisterE, 4),
		op.ADD_A_H:  newAdd(RegisterH, 4),
		op.ADD_A_L:  newAdd(RegisterL, 4),
		op.ADD_A_HL: newAdd(AddressHL, 8),
		op.ADD_A_N:  newAdd(Immediate, 8),

		op.ADC_A_A:  newAddC(RegisterA, 4),
		op.ADC_A_B:  newAddC(RegisterB, 4),
		op.ADC_A_C:  newAddC(RegisterC, 4),
		op.ADC_A_D:  newAddC(RegisterD, 4),
		op.ADC_A_E:  newAddC(RegisterE, 4),
		op.ADC_A_H:  newAddC(RegisterH, 4),
		op.ADC_A_L:  newAddC(RegisterL, 4),
		op.ADC_A_HL: newAddC(AddressHL, 8),
		op.ADC_A_N:  newAddC(Immediate, 8),

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
