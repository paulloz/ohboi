package lcd

import (
	"os"

	"math/rand"
	"time"

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
	Render([Width * Height]*color)
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

	ly := lcd.io.Read(io.LY)

	stat := lcd.io.Read(io.STAT)
	modeChanged := func(currentMode uint8) func() bool {
		return func() bool {
			return currentMode != stat&0x3
		}
	}(stat & 0x03)

	if ly >= 144 {
		// We're in V-BLANK, set mode bits to 01
		stat = bits.Reset(1, bits.Set(0, stat))
		if bits.Test(4, stat) && modeChanged() {
			// We changed mode and V-Blank STAT interrupt is enabled
			lcd.cpu.RequestInterrupt(cpu.I_LCDSTAT)
		}
	} else {
		// The end of V-Blank to begin of next V-Blank period takes 456 CPU cycles
		// It is timed like this: Mode2 (80 cycles) -> Mode3 (170~240 cycles) -> Mode0 (remaining cycles)
		if lcd.scanlineCounter < 80 {
			// We're in Mode2, set mode bits to 11
			stat = bits.Set(1, bits.Set(0, stat))
			if bits.Test(5, stat) && modeChanged() {
				// We changed mode and Mode2 interrupt is enabled
				lcd.cpu.RequestInterrupt(cpu.I_LCDSTAT)
			}
		} else if lcd.scanlineCounter < 80+170 {
			// We're in Mode3, set mode bits to 10, no interrupt for Mode3
			stat = bits.Set(1, bits.Reset(0, stat))
		} else {
			// We're in Mode0, set mode bits to 00
			stat = bits.Reset(1, bits.Reset(0, stat))
			if bits.Test(3, stat) && modeChanged() {
				// We changed mode and Mode0 interrupt is enabled
				lcd.cpu.RequestInterrupt(cpu.I_LCDSTAT)
			}
		}
	}

	// if LYC == LYC, must set bit 2 and request interrupt if bit 6 is set
	// must reset bit 2 otherwise
	lyc := lcd.io.Read(io.LYC)
	if lyc == ly {
		bits.Set(2, stat)
		if bits.Test(6, stat) {
			lcd.cpu.RequestInterrupt(cpu.I_LCDSTAT)
		}
	} else {
		bits.Reset(2, stat)
	}

	lcd.io.Write(io.STAT, stat)
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
	// This is a test
	pixels := [Width * Height]*color{}
	for i := range pixels {
		pixels[i] = newColor(0xb6, 0xb6, 0xb6)
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < r.Intn(1000); i++ {
		pixels[r.Intn(len(pixels))] = newColor(0, 0, 0)
	}
	lcd.backend.Render(pixels)
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

	return lcd
}
