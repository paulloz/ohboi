package gameboy

import (
	"time"

	"github.com/paulloz/ohboi/apu"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/memory"
)

// Speed constants
const (
	ClockSpeed     = uint32(4194304)              // Cycles per second
	DIVSpeed       = 16384                        // DIV increments per second
	FPS            = 60                           // We want to run at 60FPS
	CyclesPerFrame = ClockSpeed / FPS             // Cycles per frame
	DIVPerFrame    = DIVSpeed / FPS               // DIV increments per frame
	CyclesPerDIV   = CyclesPerFrame / DIVPerFrame // Cycles to wait between DIV increments
)

type GameBoy struct {
	apu    *apu.APU
	cpu    *cpu.CPU
	io     *io.IO
	memory *memory.Memory
}

func (gb *GameBoy) Panic(err error) {
	debuggerStop()
	panic(err)
}

func (gb *GameBoy) Update(pendingDIVCycles uint32) (uint32, uint32) {
	var cycles uint32
	var cyclesDIV uint32 = pendingDIVCycles

	for cycles = 0; cycles < CyclesPerFrame; {
		debuggerStep()

		// Execute instruction
		_cycles, err := gb.cpu.ExecuteOpCode()
		if err != nil {
			gb.Panic(err)
		}

		// Update DIV register
		cyclesDIV += uint32(_cycles)
		if cyclesDIV >= CyclesPerDIV {
			cyclesDIV -= CyclesPerDIV
			gb.cpu.IncrementDIV()
		}

		cycles += uint32(_cycles)
	}

	return cycles, cyclesDIV
}

func (gb *GameBoy) InsertCartridgeFromFile(filename string) {
	gb.memory.LoadCartridgeFromFile(filename)
}

func (gb *GameBoy) PowerOn() {
	ticker := time.NewTicker(time.Second / FPS).C

	debuggerStart(gb)

	start := time.Now()
	frames := 0

	var pendingDIVCycles uint32

	for {
		select {
		case <-ticker:
			_, pendingDIVCycles = gb.Update(pendingDIVCycles)

			frames++
			if time.Since(start) > time.Second {
				start = time.Now()
				frames = 0
			}
		}
	}
}

func NewGameBoy() *GameBoy {
	io := io.NewIO()
	memory := memory.NewMemory(io)
	cpu := cpu.NewCPU(memory, io)
	apu := apu.NewAPU(io)

	return &GameBoy{
		apu:    apu,
		cpu:    cpu,
		io:     io,
		memory: memory,
	}
}
