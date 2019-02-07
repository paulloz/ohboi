package cpu

import (
	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newRL(getset GetterSetter, quick bool) Instruction {
	cycles := uint(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			out := (in << 1) | bits.FromBool(cpu.GetCFlag())

			getset.Set(cpu, out)
			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(bits.Test(7, in))

			return nil
		},
		Cycles: cycles,
	}
}

func newRLC(getset GetterSetter, quick bool) Instruction {
	cycles := uint(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			oldSeventh := in >> 7
			out := (in << 1) | oldSeventh

			getset.Set(cpu, out)
			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(oldSeventh == 1)

			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.RLA: newRL(RegisterA, true),

		op.RLCA: newRLC(RegisterA, true),
	})
	RegisterExtInstructions(map[uint8]Instruction{
		op.RL_A:  newRL(RegisterA, false),
		op.RL_B:  newRL(RegisterB, false),
		op.RL_C:  newRL(RegisterC, false),
		op.RL_D:  newRL(RegisterD, false),
		op.RL_E:  newRL(RegisterE, false),
		op.RL_H:  newRL(RegisterH, false),
		op.RL_L:  newRL(RegisterL, false),
		op.RL_HL: newRL(AddressHL, false),

		op.RLC_A:  newRLC(RegisterA, false),
		op.RLC_B:  newRLC(RegisterB, false),
		op.RLC_C:  newRLC(RegisterC, false),
		op.RLC_D:  newRLC(RegisterD, false),
		op.RLC_E:  newRLC(RegisterE, false),
		op.RLC_H:  newRLC(RegisterH, false),
		op.RLC_L:  newRLC(RegisterL, false),
		op.RLC_HL: newRLC(AddressHL, false),
	})
}
