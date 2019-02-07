package cpu

import (
	"github.com/paulloz/ohboi/bits"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

type bitopcode struct {
	name  string
	start uint8
	end   uint8
}

func newBitHandler(src Getter, bit uint8) func(*CPU, *memory.Memory) error {
	return func(cpu *CPU, mem *memory.Memory) error {
		isSet := bits.Test(bit, src.Get(cpu))
		cpu.SetZFlag(!isSet)
		cpu.SetNFlag(false)
		cpu.SetHFlag(true)
		return nil
	}
}

func newResetHandler(register GetterSetter, bit uint8) func(*CPU, *memory.Memory) error {
	return func(cpu *CPU, mem *memory.Memory) error {
		register.Set(cpu, bits.Reset(bit, register.Get(cpu)))
		return nil
	}
}

func newSetHandler(register GetterSetter, bit uint8) func(*CPU, *memory.Memory) error {
	return func(cpu *CPU, mem *memory.Memory) error {
		register.Set(cpu, bits.Set(bit, register.Get(cpu)))
		return nil
	}
}

func init() {
	instructions := make(map[uint8]Instruction)

	bitopcodes := [3]bitopcode{
		bitopcode{name: "BIT", start: op.BIT_0_B, end: op.BIT_7_A},
		bitopcode{name: "RES", start: op.RES_0_B, end: op.RES_7_A},
		bitopcode{name: "SET", start: op.SET_0_B, end: op.SET_7_A},
	}

	for _, bitopcode := range bitopcodes {
		for hi := uint8(bitopcode.start >> 4); hi <= (bitopcode.end >> 4); hi++ {
			for lo := uint8(0x0); lo <= 0xf; lo++ {
				opCode := (hi << 4) | lo
				bit := (opCode - bitopcode.start) / 8

				var register GetterSetter

				switch {
				case lo == 0x0 || lo == 0x8:
					register = RegisterB
				case lo == 0x1 || lo == 0x9:
					register = RegisterC
				case lo == 0x2 || lo == 0xa:
					register = RegisterD
				case lo == 0x3 || lo == 0xb:
					register = RegisterE
				case lo == 0x4 || lo == 0xc:
					register = RegisterH
				case lo == 0x5 || lo == 0xd:
					register = RegisterL
				case lo == 0x6 || lo == 0xe:
					register = AddressHL
				case lo == 0x7 || lo == 0xf:
					register = RegisterA
				}

				cycles := uint(8)
				if register == AddressHL {
					cycles = 16
				}

				var handler func(*CPU, *memory.Memory) error

				switch bitopcode.name {
				case "BIT":
					handler = newBitHandler(Getter(register), bit)
				case "RES":
					handler = newResetHandler(register, bit)
				case "SET":
					handler = newSetHandler(register, bit)
				}

				if handler != nil {
					instructions[opCode] = Instruction{Handler: handler, Cycles: cycles}
				}
			}
		}
	}

	RegisterExtInstructions(instructions)
}
