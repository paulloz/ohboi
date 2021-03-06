package cpu

import (
	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newAdd(src Getter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			a, b := cpu.A.Get(), src.Get(cpu)
			sum := uint16(a) + uint16(b)
			uint8Sum := uint8(sum)
			cpu.A.Set(uint8Sum)

			cpu.SetZFlag(uint8Sum == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(((a&0xf)+(b&0xf))&0x10 != 0)
			cpu.SetCFlag(sum >= 256)
			return nil
		},
		Cycles: cycles,
	}
}

func newAdd16(src Getter16, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			a, b := cpu.HL.Get(), src.Get(cpu)

			sum := uint32(a) + uint32(b)
			cpu.HL.Set(a + b)

			cpu.SetNFlag(false)
			cpu.SetHFlag(((a & 0xfff) + (b & 0xfff)) > 0xfff)
			cpu.SetCFlag(sum > 0xffff)
			return nil
		},
		Cycles: cycles,
	}
}

func newAddC(src Getter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			var carry uint8
			if cpu.F.Get()&CarryFlag != 0 {
				carry = 1
			}
			a, b := cpu.A.Get(), src.Get(cpu)
			sum := uint16(a) + uint16(b) + uint16(carry)
			uint8Sum := uint8(sum)
			cpu.A.Set(uint8Sum)

			cpu.SetZFlag(uint8Sum == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(((a&0xf)+(b&0xf)+carry)&0x10 != 0)
			cpu.SetCFlag(sum >= 256)
			return nil
		},
		Cycles: cycles,
	}
}

func newSub(src Getter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			a, b := cpu.A.Get(), src.Get(cpu)
			tot := a - b

			cpu.A.Set(tot)

			cpu.SetZFlag(tot == 0)
			cpu.SetNFlag(true)
			cpu.SetHFlag((a & 0x0f) < (b & 0x0f))
			cpu.SetCFlag(a < b)
			return nil
		},
		Cycles: cycles,
	}
}

func newSubC(src Getter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			carry := int16(bits.FromBool(cpu.GetCFlag()))
			a, b := int16(cpu.A.Get()), int16(src.Get(cpu))
			sub := a - b - carry
			cpu.A.Set(uint8(sub))

			cpu.SetZFlag(cpu.A.Get() == 0)
			cpu.SetNFlag(true)
			cpu.SetHFlag((a&0x0f)-(b&0x0f)-carry < 0)
			cpu.SetCFlag(a < b+carry)
			return nil
		},
		Cycles: cycles,
	}
}

func newAnd(register Getter, cycles uint32) Instruction {
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

func newOr(register Getter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			result := register.Get(cpu) | cpu.A.Get()
			cpu.A.Set(result)

			cpu.SetZFlag(result == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(false)

			return nil
		},
		Cycles: cycles,
	}
}

func newCompareA(src Getter) Instruction {
	cycles := uint32(4)
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
	cycles := uint32(4)
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
	cycles := uint32(4)
	if register == AddressHL {
		cycles = 12
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			initial := register.Get(cpu)
			final := initial - 1
			register.Set(cpu, final)

			cpu.SetZFlag(final == 0)
			cpu.SetNFlag(true)
			cpu.SetHFlag((initial & 0x0f) < 1)

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
	cycles := uint32(4)
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
	RegisterInstructions(map[uint8]Instruction{
		op.ADD_A_A:  newAdd(RegisterA, 4),
		op.ADD_A_B:  newAdd(RegisterB, 4),
		op.ADD_A_C:  newAdd(RegisterC, 4),
		op.ADD_A_D:  newAdd(RegisterD, 4),
		op.ADD_A_E:  newAdd(RegisterE, 4),
		op.ADD_A_H:  newAdd(RegisterH, 4),
		op.ADD_A_L:  newAdd(RegisterL, 4),
		op.ADD_A_HL: newAdd(AddressHL, 8),
		op.ADD_A_N:  newAdd(Immediate, 8),

		op.ADD_HL_BC: newAdd16(RegisterBC, 8),
		op.ADD_HL_DE: newAdd16(RegisterDE, 8),
		op.ADD_HL_HL: newAdd16(RegisterHL, 8),
		op.ADD_HL_SP: newAdd16(RegisterSP, 8),

		op.ADC_A_A:  newAddC(RegisterA, 4),
		op.ADC_A_B:  newAddC(RegisterB, 4),
		op.ADC_A_C:  newAddC(RegisterC, 4),
		op.ADC_A_D:  newAddC(RegisterD, 4),
		op.ADC_A_E:  newAddC(RegisterE, 4),
		op.ADC_A_H:  newAddC(RegisterH, 4),
		op.ADC_A_L:  newAddC(RegisterL, 4),
		op.ADC_A_HL: newAddC(AddressHL, 8),
		op.ADC_A_N:  newAddC(Immediate, 8),

		op.SUB_A_A:  newSub(RegisterA, 4),
		op.SUB_A_B:  newSub(RegisterB, 4),
		op.SUB_A_C:  newSub(RegisterC, 4),
		op.SUB_A_D:  newSub(RegisterD, 4),
		op.SUB_A_E:  newSub(RegisterE, 4),
		op.SUB_A_H:  newSub(RegisterH, 4),
		op.SUB_A_L:  newSub(RegisterL, 4),
		op.SUB_A_HL: newSub(AddressHL, 8),
		op.SUB_A_N:  newSub(Immediate, 8),

		op.SBC_A_A:  newSubC(RegisterA, 4),
		op.SBC_A_B:  newSubC(RegisterB, 4),
		op.SBC_A_C:  newSubC(RegisterC, 4),
		op.SBC_A_D:  newSubC(RegisterD, 4),
		op.SBC_A_E:  newSubC(RegisterE, 4),
		op.SBC_A_H:  newSubC(RegisterH, 4),
		op.SBC_A_L:  newSubC(RegisterL, 4),
		op.SBC_A_HL: newSubC(AddressHL, 8),
		op.SBC_A_N:  newSubC(Immediate, 8),

		op.AND_A_A:  newAnd(RegisterA, 4),
		op.AND_A_B:  newAnd(RegisterB, 4),
		op.AND_A_C:  newAnd(RegisterC, 4),
		op.AND_A_D:  newAnd(RegisterD, 4),
		op.AND_A_E:  newAnd(RegisterE, 4),
		op.AND_A_H:  newAnd(RegisterH, 4),
		op.AND_A_L:  newAnd(RegisterL, 4),
		op.AND_A_HL: newAnd(AddressHL, 8),
		op.AND_A_N:  newAnd(Immediate, 8),

		op.OR_A_A:  newOr(RegisterA, 4),
		op.OR_A_B:  newOr(RegisterB, 4),
		op.OR_A_C:  newOr(RegisterC, 4),
		op.OR_A_D:  newOr(RegisterD, 4),
		op.OR_A_E:  newOr(RegisterE, 4),
		op.OR_A_H:  newOr(RegisterH, 4),
		op.OR_A_L:  newOr(RegisterL, 4),
		op.OR_A_HL: newOr(AddressHL, 8),
		op.OR_A_N:  newOr(Immediate, 8),

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

		op.ADD_SP_N: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				in := int32(cpu.SP.hilo)
				rel := int32(int8(cpu.FetchByte()))

				result := in + rel
				cpu.SP.Set(uint16(result))

				overflowTest := (in ^ rel ^ (result & 0xffff))

				cpu.SetZFlag(false)
				cpu.SetNFlag(false)
				cpu.SetHFlag((overflowTest & 0x10) == 0x10)
				cpu.SetCFlag((overflowTest & 0x100) == 0x100)

				return nil
			},
			Cycles: 16,
		},
	})
}
