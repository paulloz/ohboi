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
func (gb *GameBoy) ExecuteNextOpCode() (uint, error) {
	opCode := gb.memory.Read(gb.cpu.AdvancePC())
	return gb.cpu.ExecuteOpCode(opCode)
}

func (gb *GameBoy) Panic(err error) {
	panic(err)
}

// Update ...
func (gb *GameBoy) Update() uint {
	var cycles uint

	for cycles = 0; cycles < CyclesPerFrame; {
		_cycles, err := gb.ExecuteNextOpCode()
		if err != nil {
			gb.Panic(err)
		}

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
	memory := memory.NewMemory()
	cpu := cpu.NewCPU(memory)

	return &GameBoy{
		cpu:    cpu,
		memory: memory,
	}
}
