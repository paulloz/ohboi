package cpu

import (
	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newRL(getset GetterSetter, quick bool) Instruction {
	cycles := uint32(8)
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

func newRR(getset GetterSetter, quick bool) Instruction {
	cycles := uint32(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			out := (in >> 1) | (bits.FromBool(cpu.GetCFlag()) << 7)

			getset.Set(cpu, out)
			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(in&1 != 0)

			return nil
		},
		Cycles: cycles,
	}
}

func newRLC(getset GetterSetter, quick bool) Instruction {
	cycles := uint32(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			oldHighest := in >> 7
			out := (in << 1) | oldHighest

			getset.Set(cpu, out)
			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(oldHighest == 1)

			return nil
		},
		Cycles: cycles,
	}
}

func newRRC(getset GetterSetter, quick bool) Instruction {
	cycles := uint32(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			oldLowest := in & 1
			out := (in >> 1) | (oldLowest << 7)

			getset.Set(cpu, out)
			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(oldLowest == 1)

			return nil
		},
		Cycles: cycles,
	}
}

func newSLA(getset GetterSetter, quick bool) Instruction {
	cycles := uint32(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			out := in << 1

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

func newSRA(getset GetterSetter, quick bool) Instruction {
	cycles := uint32(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			out := in>>1 | in&0x80

			getset.Set(cpu, out)
			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(in&1 != 0)

			return nil
		},
		Cycles: cycles,
	}
}

func newSRL(getset GetterSetter, quick bool) Instruction {
	cycles := uint32(8)
	if getset == AddressHL {
		cycles = 16
	}
	if quick {
		cycles /= 2
	}

	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)
			out := in >> 1

			getset.Set(cpu, out)
			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(in&1 != 0)

			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.RLA:  newRL(RegisterA, true),
		op.RRA:  newRR(RegisterA, true),
		op.RLCA: newRLC(RegisterA, true),
		op.RRCA: newRRC(RegisterA, true),
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

		op.RR_A:  newRR(RegisterA, false),
		op.RR_B:  newRR(RegisterB, false),
		op.RR_C:  newRR(RegisterC, false),
		op.RR_D:  newRR(RegisterD, false),
		op.RR_E:  newRR(RegisterE, false),
		op.RR_H:  newRR(RegisterH, false),
		op.RR_L:  newRR(RegisterL, false),
		op.RR_HL: newRR(AddressHL, false),

		op.RLC_A:  newRLC(RegisterA, false),
		op.RLC_B:  newRLC(RegisterB, false),
		op.RLC_C:  newRLC(RegisterC, false),
		op.RLC_D:  newRLC(RegisterD, false),
		op.RLC_E:  newRLC(RegisterE, false),
		op.RLC_H:  newRLC(RegisterH, false),
		op.RLC_L:  newRLC(RegisterL, false),
		op.RLC_HL: newRLC(AddressHL, false),

		op.RRC_A:  newRRC(RegisterA, false),
		op.RRC_B:  newRRC(RegisterB, false),
		op.RRC_C:  newRRC(RegisterC, false),
		op.RRC_D:  newRRC(RegisterD, false),
		op.RRC_E:  newRRC(RegisterE, false),
		op.RRC_H:  newRRC(RegisterH, false),
		op.RRC_L:  newRRC(RegisterL, false),
		op.RRC_HL: newRRC(AddressHL, false),

		op.SLA_A:  newSLA(RegisterA, false),
		op.SLA_B:  newSLA(RegisterB, false),
		op.SLA_C:  newSLA(RegisterC, false),
		op.SLA_D:  newSLA(RegisterD, false),
		op.SLA_E:  newSLA(RegisterE, false),
		op.SLA_H:  newSLA(RegisterH, false),
		op.SLA_L:  newSLA(RegisterL, false),
		op.SLA_HL: newSLA(AddressHL, false),

		op.SRA_A:  newSRA(RegisterA, false),
		op.SRA_B:  newSRA(RegisterB, false),
		op.SRA_C:  newSRA(RegisterC, false),
		op.SRA_D:  newSRA(RegisterD, false),
		op.SRA_E:  newSRA(RegisterE, false),
		op.SRA_H:  newSRA(RegisterH, false),
		op.SRA_L:  newSRA(RegisterL, false),
		op.SRA_HL: newSRA(AddressHL, false),

		op.SRL_A:  newSRL(RegisterA, false),
		op.SRL_B:  newSRL(RegisterB, false),
		op.SRL_C:  newSRL(RegisterC, false),
		op.SRL_D:  newSRL(RegisterD, false),
		op.SRL_E:  newSRL(RegisterE, false),
		op.SRL_H:  newSRL(RegisterH, false),
		op.SRL_L:  newSRL(RegisterL, false),
		op.SRL_HL: newSRL(AddressHL, false),
	})
}
