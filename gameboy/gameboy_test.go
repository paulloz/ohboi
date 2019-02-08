package gameboy_test

import (
	"testing"

	"github.com/paulloz/ohboi/cartridge"
	"github.com/paulloz/ohboi/cpu"
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/memory"
)

func TestDisableROM(t *testing.T) {
	data := []byte{123}
	rom := cartridge.NewROM(data)

	memory := memory.NewMemory(io.NewIO())
	memory.LoadCartridge(&cartridge.Cartridge{MBC: rom})

	cpu := cpu.NewCPU(memory)

	if memory.Read(0) != op.LD_SP_NN {
		t.Error("Expected first byte to be LD_SP_NN")
	}

	cpu.PC = 0xfc

	cpu.ExecuteOpCode()
	cpu.ExecuteOpCode()

	if cpu.PC != 0x100 {
		t.Error("Expected PC to be 0x100")
	}

	if memory.Read(0) != 123 {
		t.Error("Expected first byte to be 123")
	}
}
