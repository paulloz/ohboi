package lcd

import (
	"os"

	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
)

const (
	Width  = 160
	Height = 144

	ScanlineFrequency = 456
)

type backend interface {
	Initialize(string)
	Render([Width * Height]color)
}

type LCD struct {
	backend backend

	cpu *cpu.CPU
	io  *io.IO

	scanlineCounter uint32
}

func (lcd *LCD) Update(cycles uint32) {
	lcd.SetLCDSTAT()

	if !lcd.io.ReadBit(io.LDCD, 7) {
		// If LCD is disabled
		return
	}

	lcd.scanlineCounter += cycles
	if lcd.scanlineCounter >= ScanlineFrequency {
		lcd.scanlineCounter = 0

		ly := lcd.io.Read(io.LY) + 1
		lcd.io.Write(io.LY, ly)

		if ly > 153 {
			// Can send screen data to backend
			lcd.io.Write(io.LY, 0)
		} else if ly >= 144 {
			lcd.cpu.RequestInterrupt(cpu.I_VBLANK)
		} else {
			lcd.DrawScanline()
		}
	}
}

func (lcd *LCD) SetLCDSTAT() {
	if !lcd.io.ReadBit(io.LDCD, 7) {
		// Reset everything
		lcd.scanlineCounter = 0
		lcd.io.Write(io.LY, 0)
		lcd.io.Write(io.STAT, 1<<2)
		return
	}
}

func (lcd *LCD) DrawScanline() {
	lcdc := lcd.io.Read(io.LDCD)

	if bits.Test(0, lcdc) {
		// Render Tiles
	}

	if bits.Test(5, lcdc) {
		// Render Sprites
	}
}

func (lcd *LCD) RenderFrame() {
	// lcd.backend.Render()
}

func NewLCD(cpu *cpu.CPU, io_ *io.IO) *LCD {
	backend := NewSDL2()
	backend.Initialize(os.Args[0])

	lcd := &LCD{
		backend: backend,

		cpu: cpu,
		io:  io_,

		scanlineCounter: 0,
	}

	// This is a test
	pixels := [Width * Height]color{}
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if x == y {
				pixels[x+Width*y] = newColor(255, 0, 255)
			} else {
				pixels[x+Width*y] = newColor(0, 0, 0)
			}
		}
	}
	backend.Render(pixels)

	return lcd
}
