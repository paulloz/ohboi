package gameboy

import (
	"time"

	"github.com/paulloz/ohboi/apu"
	"github.com/paulloz/ohboi/bits"
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

	timaClock uint32
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

		gb.UpdateTimers(_cycles)

		cycles += uint32(_cycles)
	}

	return cycles, cyclesDIV
}

func (gb *GameBoy) UpdateTimers(cycles uint32) {
	// If Timer Enable (bit 2 of TAC is set)
	if bits.Test(2, gb.io.Read(io.TAC)) {
		// Retrieve frequency
		timaClockFrequency := uint32(1024)
		switch gb.io.Read(io.TAC) & 0x03 {
		case 1:
			timaClockFrequency = 16
		case 2:
			timaClockFrequency = 64
		case 3:
			timaClockFrequency = 256
		}

		gb.timaClock += cycles
		for gb.timaClock >= timaClockFrequency {
			gb.timaClock -= timaClockFrequency

			tima := gb.io.Read(io.TIMA)
			if tima < 0xff {
				gb.io.Write(io.TIMA, tima+1)
			} else {
				gb.io.Write(io.TIMA, gb.io.Read(io.TMA))
				// Request Timer Interrupt (IF bit 2)
			}
		}
	}
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
	io_ := io.NewIO()
	memory := memory.NewMemory(io_)
	cpu := cpu.NewCPU(memory, io_)
	apu := apu.NewAPU(io_)

	return &GameBoy{
		apu:    apu,
		cpu:    cpu,
		io:     io_,
		memory: memory,
	}
}
