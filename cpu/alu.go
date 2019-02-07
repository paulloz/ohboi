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

func newAnd(register Getter, cycles uint) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			result := register.Get(cpu) & cpu.A.Get()
			cpu.A.Set(result)

			cpu.SetZFlag(result == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(true)
			cpu.SetCFlag(false)

			return nil
		},
		Cycles: cycles,
	}
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

func newIncrementRegister16(register GetterSetter16) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			register.Set(cpu, register.Get(cpu)+1)
			return nil
		},
		Cycles: 8,
	}
}

func newDecrementRegister(register GetterSetter) Instruction {
	cycles := uint(4)
	if register == AddressHL {
		cycles = 12
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			initial := register.Get(cpu)
			final := initial - 1
			register.Set(cpu, final)

			cpu.SetZFlag(final == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(bits.HalfCarryCheck(initial, 1))

			return nil
		},
		Cycles: cycles,
	}
}

func newDecrementRegister16(register GetterSetter16) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			register.Set(cpu, register.Get(cpu)-1)
			return nil
		},
		Cycles: 8,
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

		op.AND_A_A:  newAnd(RegisterA, 4),
		op.AND_A_B:  newAnd(RegisterB, 4),
		op.AND_A_C:  newAnd(RegisterC, 4),
		op.AND_A_D:  newAnd(RegisterD, 4),
		op.AND_A_E:  newAnd(RegisterE, 4),
		op.AND_A_H:  newAnd(RegisterH, 4),
		op.AND_A_L:  newAnd(RegisterL, 4),
		op.AND_A_HL: newAnd(AddressHL, 8),
		op.AND_A_N:  newAnd(Immediate, 8),

		op.CP_A:  newCompareA(RegisterA),
		op.CP_B:  newCompareA(RegisterB),
		op.CP_C:  newCompareA(RegisterC),
		op.CP_D:  newCompareA(RegisterD),
		op.CP_E:  newCompareA(RegisterE),
		op.CP_H:  newCompareA(RegisterH),
		op.CP_L:  newCompareA(RegisterL),
		op.CP_HL: newCompareA(AddressHL),
		op.CP_N:  newCompareA(Immediate),

		op.INC_A:   newIncrementRegister(RegisterA),
		op.INC_B:   newIncrementRegister(RegisterB),
		op.INC_C:   newIncrementRegister(RegisterC),
		op.INC_D:   newIncrementRegister(RegisterD),
		op.INC_E:   newIncrementRegister(RegisterE),
		op.INC_H:   newIncrementRegister(RegisterH),
		op.INC_L:   newIncrementRegister(RegisterL),
		op.INC_HLA: newIncrementRegister(AddressHL),

		op.INC_BC: newIncrementRegister16(RegisterBC),
		op.INC_DE: newIncrementRegister16(RegisterDE),
		op.INC_HL: newIncrementRegister16(RegisterHL),
		op.INC_SP: newIncrementRegister16(RegisterSP),

		op.DEC_A:   newDecrementRegister(RegisterA),
		op.DEC_B:   newDecrementRegister(RegisterB),
		op.DEC_C:   newDecrementRegister(RegisterC),
		op.DEC_D:   newDecrementRegister(RegisterD),
		op.DEC_E:   newDecrementRegister(RegisterE),
		op.DEC_H:   newDecrementRegister(RegisterH),
		op.DEC_L:   newDecrementRegister(RegisterL),
		op.DEC_HLA: newDecrementRegister(AddressHL),

		op.DEC_BC: newDecrementRegister16(RegisterBC),
		op.DEC_DE: newDecrementRegister16(RegisterDE),
		op.DEC_HL: newDecrementRegister16(RegisterHL),
		op.DEC_SP: newDecrementRegister16(RegisterSP),

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
