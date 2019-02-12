package cpu

import (
	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func rl(v uint8, c bool) uint8 {
	return (v << 1) | bits.FromBool(c)
}

func newRL(getset GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)

			out := rl(in, cpu.GetCFlag())
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

func rr(v uint8, c bool) uint8 {
	return (bits.FromBool(c) << 7) | (v >> 1)
}

func newRR(getset GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)

			out := rr(in, cpu.GetCFlag())
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

func rlc(v uint8) uint8 {
	return (v << 1) | (v >> 7)
}

func newRLC(getset GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)

			out := rlc(in)
			getset.Set(cpu, out)

			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(in > 0x7f)

			return nil
		},
		Cycles: cycles,
	}
}

func rrc(v uint8) uint8 {
	return ((v & 0x01) << 7) | (v >> 1)
}

func newRRC(getset GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)

			out := rrc(in)
			getset.Set(cpu, out)

			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag((in & 0x01) == 1)

			return nil
		},
		Cycles: cycles,
	}
}

func newSLA(getset GetterSetter, cycles uint32) Instruction {
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

func newSRA(getset GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)

			out := (in >> 1) | (in & 0x80)
			getset.Set(cpu, out)

			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(bits.Test(0, in))

			return nil
		},
		Cycles: cycles,
	}
}

func newSRL(getset GetterSetter, cycles uint32) Instruction {
	return Instruction{
		Handler: func(cpu *CPU, mem *memory.Memory) error {
			in := getset.Get(cpu)

			out := in >> 1
			getset.Set(cpu, out)

			cpu.SetZFlag(out == 0)
			cpu.SetNFlag(false)
			cpu.SetHFlag(false)
			cpu.SetCFlag(bits.Test(0, in))

			return nil
		},
		Cycles: cycles,
	}
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.RLA: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				in := cpu.A.Get()

				cpu.A.Set(rl(in, cpu.GetCFlag()))

				cpu.SetZFlag(false)
				cpu.SetNFlag(false)
				cpu.SetHFlag(false)
				cpu.SetCFlag(bits.Test(7, in))

				return nil
			},
			Cycles: 4,
		},
		op.RRA: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				in := cpu.A.Get()

				cpu.A.Set(rr(in, cpu.GetCFlag()))

				cpu.SetZFlag(false)
				cpu.SetNFlag(false)
				cpu.SetHFlag(false)
				cpu.SetCFlag(in&1 != 0)

				return nil
			},
			Cycles: 4,
		},
		op.RLCA: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				in := cpu.A.Get()

				cpu.A.Set(rlc(in))

				cpu.SetZFlag(false)
				cpu.SetNFlag(false)
				cpu.SetHFlag(false)
				cpu.SetCFlag(in > 0x7f)

				return nil
			},
			Cycles: 4,
		},
		op.RRCA: Instruction{
			Handler: func(cpu *CPU, mem *memory.Memory) error {
				in := cpu.A.Get()

				cpu.A.Set(rrc(in))

				cpu.SetZFlag(false)
				cpu.SetNFlag(false)
				cpu.SetHFlag(false)
				cpu.SetCFlag((in & 0x01) == 1)

				return nil
			},
			Cycles: 4,
		},
	})
	RegisterExtInstructions(map[uint8]Instruction{
		op.RL_A:  newRL(RegisterA, 8),
		op.RL_B:  newRL(RegisterB, 8),
		op.RL_C:  newRL(RegisterC, 8),
		op.RL_D:  newRL(RegisterD, 8),
		op.RL_E:  newRL(RegisterE, 8),
		op.RL_H:  newRL(RegisterH, 8),
		op.RL_L:  newRL(RegisterL, 8),
		op.RL_HL: newRL(AddressHL, 16),

		op.RR_A:  newRR(RegisterA, 8),
		op.RR_B:  newRR(RegisterB, 8),
		op.RR_C:  newRR(RegisterC, 8),
		op.RR_D:  newRR(RegisterD, 8),
		op.RR_E:  newRR(RegisterE, 8),
		op.RR_H:  newRR(RegisterH, 8),
		op.RR_L:  newRR(RegisterL, 8),
		op.RR_HL: newRR(AddressHL, 16),

		op.RLC_A:  newRLC(RegisterA, 8),
		op.RLC_B:  newRLC(RegisterB, 8),
		op.RLC_C:  newRLC(RegisterC, 8),
		op.RLC_D:  newRLC(RegisterD, 8),
		op.RLC_E:  newRLC(RegisterE, 8),
		op.RLC_H:  newRLC(RegisterH, 8),
		op.RLC_L:  newRLC(RegisterL, 8),
		op.RLC_HL: newRLC(AddressHL, 16),

		op.RRC_A:  newRRC(RegisterA, 8),
		op.RRC_B:  newRRC(RegisterB, 8),
		op.RRC_C:  newRRC(RegisterC, 8),
		op.RRC_D:  newRRC(RegisterD, 8),
		op.RRC_E:  newRRC(RegisterE, 8),
		op.RRC_H:  newRRC(RegisterH, 8),
		op.RRC_L:  newRRC(RegisterL, 8),
		op.RRC_HL: newRRC(AddressHL, 16),

		op.SLA_A:  newSLA(RegisterA, 8),
		op.SLA_B:  newSLA(RegisterB, 8),
		op.SLA_C:  newSLA(RegisterC, 8),
		op.SLA_D:  newSLA(RegisterD, 8),
		op.SLA_E:  newSLA(RegisterE, 8),
		op.SLA_H:  newSLA(RegisterH, 8),
		op.SLA_L:  newSLA(RegisterL, 8),
		op.SLA_HL: newSLA(AddressHL, 16),

		op.SRA_A:  newSRA(RegisterA, 8),
		op.SRA_B:  newSRA(RegisterB, 8),
		op.SRA_C:  newSRA(RegisterC, 8),
		op.SRA_D:  newSRA(RegisterD, 8),
		op.SRA_E:  newSRA(RegisterE, 8),
		op.SRA_H:  newSRA(RegisterH, 8),
		op.SRA_L:  newSRA(RegisterL, 8),
		op.SRA_HL: newSRA(AddressHL, 16),

		op.SRL_A:  newSRL(RegisterA, 8),
		op.SRL_B:  newSRL(RegisterB, 8),
		op.SRL_C:  newSRL(RegisterC, 8),
		op.SRL_D:  newSRL(RegisterD, 8),
		op.SRL_E:  newSRL(RegisterE, 8),
		op.SRL_H:  newSRL(RegisterH, 8),
		op.SRL_L:  newSRL(RegisterL, 8),
		op.SRL_HL: newSRL(AddressHL, 16),
	})
}
