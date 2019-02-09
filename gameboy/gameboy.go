package gameboy

import (
	"time"

	"github.com/paulloz/ohboi/apu"
	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/lcd"
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
	Memory *memory.Memory
	lcd    *lcd.LCD

	timaClock uint32
}

func (gb *GameBoy) Panic(err error) {
	debuggerStop()
	panic(err)
}

func (gb *GameBoy) Update(pendingCycles uint32) (uint32, uint32) {
	var cycles uint32
	var currentInstrCycles = pendingCycles

	for cycles = 0; cycles < CyclesPerFrame; {
		debuggerStep()

		// Execute instruction
		opCycles, err := gb.cpu.ExecuteOpCode()
		if err != nil {
			gb.Panic(err)
		}
		currentInstrCycles += opCycles

		gb.lcd.Update(currentInstrCycles)

		gb.cpu.UpdateDIV(currentInstrCycles, CyclesPerDIV)
		gb.UpdateTimers(currentInstrCycles)
		cycles += currentInstrCycles

		currentInstrCycles = gb.cpu.ManageInterrupts()
		cycles += currentInstrCycles
	}

	gb.lcd.RenderFrame()

	return cycles, currentInstrCycles
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
				gb.cpu.RequestInterrupt(cpu.I_TIMER)
			}
		}
	}
}

func (gb *GameBoy) InsertCartridgeFromFile(filename string) {
	gb.Memory.LoadCartridgeFromFile(filename)
}

func (gb *GameBoy) PowerOn() {
	ticker := time.NewTicker(time.Second / FPS).C

	debuggerStart(gb)

	start := time.Now()
	frames := 0

	pendingCycles := uint32(0)

	for {
		select {
		case <-ticker:
			_, pendingCycles = gb.Update(pendingCycles)

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
	apu := apu.NewAPU(io_)

	memory := memory.NewMemory(io_)
	cpu := cpu.NewCPU(memory, io_)

	lcd := lcd.NewLCD(cpu, memory, io_)

	return &GameBoy{
		apu:    apu,
		cpu:    cpu,
		io:     io_,
		Memory: memory,
		lcd:    lcd,
	}
}
