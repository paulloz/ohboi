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

func isB(l uint8) bool {
	return l == 0x0 || l == 0x8
}

func isC(l uint8) bool {
	return l == 0x1 || l == 0x9
}

func isD(l uint8) bool {
	return l == 0x2 || l == 0xa
}

func isE(l uint8) bool {
	return l == 0x3 || l == 0xb
}

func isH(l uint8) bool {
	return l == 0x4 || l == 0xc
}

func isL(l uint8) bool {
	return l == 0x5 || l == 0xd
}

func isHL(l uint8) bool {
	return l == 0x6 || l == 0xe
}

func isA(l uint8) bool {
	return l == 0x7 || l == 0xf
}

func newBitHandler(src Getter, bit uint8) func(*CPU, *memory.Memory) error {
	return func(cpu *CPU, mem *memory.Memory) error {
		isSet := bits.Test(bit, src.Get(cpu))
		cpu.SetZFlag(isSet)
		cpu.SetNFlag(false)
		cpu.SetHFlag(true)
		return nil
	}
}

func newResetHandler(dst Setter, src Getter, bit uint8) func(*CPU, *memory.Memory) error {
	return func(cpu *CPU, mem *memory.Memory) error {
		dst.Set(cpu, bits.Reset(bit, src.Get(cpu)))
		return nil
	}
}

func newSetHandler(dst Setter, src Getter, bit uint8) func(*CPU, *memory.Memory) error {
	return func(cpu *CPU, mem *memory.Memory) error {
		dst.Set(cpu, bits.Set(bit, src.Get(cpu)))
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

				var src Getter
				var dst Setter

				switch {
				case isB(lo):
					src = RegisterB
					dst = RegisterB
				case isC(lo):
					src = RegisterC
					dst = RegisterC
				case isD(lo):
					src = RegisterD
					dst = RegisterD
				case isE(lo):
					src = RegisterE
					dst = RegisterE
				case isH(lo):
					src = RegisterH
					dst = RegisterH
				case isL(lo):
					src = RegisterL
					dst = RegisterL
				case isHL(lo):
					src = AddressHL
					dst = AddressHL
				case isA(lo):
					src = RegisterA
					dst = RegisterA
				}

				cycles := uint(8)
				if src == AddressHL {
					cycles = 16
				}

				var handler func(*CPU, *memory.Memory) error

				switch bitopcode.name {
				case "BIT":
					handler = newBitHandler(src, bit)
				case "RES":
					handler = newResetHandler(dst, src, bit)
				case "SET":
					handler = newSetHandler(dst, src, bit)
				}

				if handler != nil {
					instructions[opCode] = Instruction{Handler: handler, Cycles: cycles}
				}
			}
		}
	}

	RegisterExtInstruction(instructions)
}
