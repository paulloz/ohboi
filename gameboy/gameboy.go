package gameboy

import (
	"fmt"
	goio "io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/paulloz/ohboi/cartridge"

	"github.com/paulloz/ohboi/apu"
	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/config"
	"github.com/paulloz/ohboi/consts"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/joypad"
	"github.com/paulloz/ohboi/memory"
	"github.com/paulloz/ohboi/ppu"
)

type GameBoy struct {
	apu    *apu.APU
	cpu    *cpu.CPU
	io     *io.IO
	Memory *memory.Memory
	ppu    *ppu.PPU
	joypad *joypad.Joypad

	timaClock uint32
	tac       uint8
	tima      uint8
}

func (gb *GameBoy) Panic(err error) {
	debuggerStop()
	panic(err)
}

func (gb *GameBoy) Update(pendingCycles uint32) (uint32, uint32) {
	var cycles uint32
	var currentInstrCycles = pendingCycles

	gb.joypad.Update()

	for cycles = 0; cycles < consts.CPUCyclesPerFrame; {
		debuggerStep()

		// Execute instruction
		opCycles, err := gb.cpu.NextCycle()
		if err != nil {
			gb.Panic(err)
		}

		gb.ppu.Update(opCycles)
		gb.apu.Update(opCycles)

		gb.UpdateTimers(opCycles)
		cycles += opCycles

		gb.cpu.UpdateDIV(opCycles)
		gb.cpu.UpdateDMA(opCycles)

		if len(gb.cpu.MicroInstructions) == 0 {
			currentInstrCycles = gb.cpu.ManageInterrupts()
			gb.cpu.UpdateDIV(currentInstrCycles)
			gb.cpu.UpdateDMA(currentInstrCycles)
			cycles += currentInstrCycles
			gb.UpdateTimers(currentInstrCycles)
		}
	}

	gb.ppu.RenderFrame()

	return cycles, currentInstrCycles
}

func (gb *GameBoy) UpdateTimers(cycles uint32) {
	if bits.Test(2, gb.tac) {
		timaClockFrequency := gb.getTACFrequency()

		gb.timaClock += cycles
		for gb.timaClock >= timaClockFrequency {
			gb.timaClock -= timaClockFrequency

			if gb.tima < 0xff {
				gb.tima++
			} else {
				gb.tima = gb.io.Read(io.TMA)
				gb.cpu.RequestInterrupt(cpu.I_TIMER)
			}
		}
	}
}

// DEBUG
func (gb *GameBoy) GETCLOCK() uint32 {
	return gb.timaClock
}

func (gb *GameBoy) GetTAC() uint8 {
	return gb.tac
}

func (gb *GameBoy) getTACFrequency() uint32 {
	switch gb.tac & 0x03 {
	case 1:
		return 16
	case 2:
		return 64
	case 3:
		return 256
	default:
		return 1024
	}
}

func (gb *GameBoy) SetTAC(value uint8) {
	oldFreq := gb.getTACFrequency()
	gb.tac = 0xf8 | value
	if gb.getTACFrequency() != oldFreq {
		gb.timaClock = 0
	}
}

func (gb *GameBoy) GetTIMA() uint8 {
	return gb.tima
}

func (gb *GameBoy) SetTIMA(value uint8) {
	gb.tima = value
	gb.timaClock = 0
}

func (gb *GameBoy) InsertCartridgeFromPath(filename string) (*cartridge.Cartridge, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return gb.InsertCartridge(filename, f)
}

func (gb *GameBoy) InsertCartridge(filename string, reader goio.Reader) (*cartridge.Cartridge, error) {
	cartridge, err := gb.Memory.LoadCartridgeFromFile(filename, reader)
	if err != nil {
		return nil, err
	}

	sav := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".sav"
	if _, err := os.Stat(sav); err == nil {
		if err := cartridge.Load(sav); err != nil {
			return cartridge, fmt.Errorf("Failed to load cartridge state from %s: %s", err, sav)
		} else {
			return cartridge, fmt.Errorf("Cartridge state successfully loaded from %s", sav)
		}
	}

	return cartridge, nil
}

func (gb *GameBoy) PowerOn(stop chan int) {
	ticker := time.NewTicker(time.Second / consts.FPS).C

	debuggerStart(gb)

	start := time.Now()

	pendingCycles := uint32(0)

	for {
		select {
		case <-ticker:
			_, pendingCycles = gb.Update(pendingCycles)

			if time.Since(start) > time.Second {
				start = time.Now()
			}
		case <-stop:
			gb.PowerOff()
			return
		}
	}
}

func (gb *GameBoy) PowerOff() {
	gb.ppu.Destroy()
	gb.apu.Destroy()

	if cartridge := gb.Memory.Cartridge(); cartridge != nil {
		filename := cartridge.Filename()
		cartridge.Save(strings.TrimSuffix(filename, filepath.Ext(filename)) + ".sav")
	}
}

func (gb *GameBoy) GetCPU() *cpu.CPU {
	return gb.cpu
}

func NewGameBoy() (*GameBoy, error) {
	io_ := io.NewIO()
	apu, err := apu.NewAPU(io_)
	if err != nil {
		return nil, err
	}

	memory := memory.NewMemory(io_)
	cpu := cpu.NewCPU(memory, io_)

	ppu := ppu.NewPPU(cpu, memory, io_)

	joypad := joypad.NewJoypad(cpu, io_)

	gb := &GameBoy{
		apu:    apu,
		cpu:    cpu,
		io:     io_,
		Memory: memory,
		ppu:    ppu,
		joypad: joypad,
	}

	io_.MapRegister(io.TAC, gb.GetTAC, gb.SetTAC)
	io_.MapRegister(io.TIMA, gb.GetTIMA, gb.SetTIMA)

	if config.Get().Emulation.SkipBoot {
		io_.Write(io.BOOTROM, 1)
		cpu.PC = 0x100
	}

	return gb, nil
}

func NewSerialTextGameBoy(f func(uint8)) (*GameBoy, error) {
	gb, err := NewGameBoy()
	if err != nil {
		return nil, err
	}

	gb.io.MapRegister(io.SC, func() uint8 { return 0xff }, func(v uint8) {
		if v == 0x81 {
			f(gb.io.Read(io.SB))
		}
	})

	return gb, nil
}
