package gameboy

import (
	"time"

	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/memory"
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

func (gb *GameBoy) Panic(err error) {
	if debug && debugger != nil {
		debugger.close()
	}
	panic(err)
}

// Update ...
func (gb *GameBoy) Update() uint {
	var cycles uint

	for cycles = 0; cycles < CyclesPerFrame; {
		if debug && debugger != nil && debugger.stepByStep {
			<-debugger.stepper
		}

		_cycles, err := gb.cpu.ExecuteOpCode()
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

func (gb *GameBoy) PowerOn() {
	ticker := time.NewTicker(time.Second / FPS).C

	if debugger != nil {
		debugger.start(gb)
	}

	start := time.Now()
	frames := 0

	for {
		select {
		case <-ticker:
			gb.Update()

			frames++
			if time.Since(start) > time.Second {
				start = time.Now()
				frames = 0
			}
		}
	}
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
