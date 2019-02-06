package gameboy

import (
	"github.com/paulloz/ohboi/gameboy/cpu"
	"github.com/paulloz/ohboi/gameboy/memory"
)

// Speed constants
const (
	ClockSpeed     = 4194304          // Cycles per second
	FPS            = 60               // We want to run at 60FPS
	CyclesPerFrame = ClockSpeed / FPS // # of cycles executed each frame
)

// GameBoy ...
type GameBoy struct {
	cpu    *cpu.CPU
	memory *memory.Memory
}

// ExecuteNextOpCode ...
func (gb *GameBoy) ExecuteNextOpCode() uint {
	opCode := gb.memory.Read(gb.cpu.AdvancePC())

	cycles := gb.cpu.ExecuteOpCode(opCode, gb.memory)

	return cycles
}

// Update ...
func (gb *GameBoy) Update() uint {
	var cycles uint

	for cycles = 0; cycles < CyclesPerFrame; {
		_cycles := gb.ExecuteNextOpCode()

		// UpdateTimers
		// UpdateGraphics
		// Interrupts

		cycles += _cycles
	}

	// RenderScreen

	return cycles
}

// InsertCartridgeFromFile ...
func (gb *GameBoy) InsertCartridgeFromFile(filename string) {
	gb.memory.LoadCartridgeFromFile(filename)
}

// NewGameBoy ...
func NewGameBoy() *GameBoy {
	return &GameBoy{
		cpu:    cpu.NewCPU(),
		memory: memory.NewMemory(),
	}
}
