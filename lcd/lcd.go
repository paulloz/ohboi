package lcd

import (
	"os"

	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/memory"
)

const (
	Width  = 160
	Height = 144

	Scale = 2

	ScanlineFrequency = 456
)

type backend interface {
	Initialize(string)
	Render([Width * Height]color)
}

type LCD struct {
	backend backend

	cpu    *cpu.CPU
	memory *memory.Memory
	io     *io.IO

	scanlineCounter   uint32
	lastDrawnScanline uint8

	workData   [Width * Height]color
	renderData [Width * Height]color
}

func (lcd *LCD) setLCDSTAT() {
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

func (lcd *LCD) getPalette(ioAddr uint8) [4]color {
	palette := lcd.io.Read(ioAddr)
	colorPalette := [4]color{}
	for i := uint8(0); i < 8; i += 2 {
		shade := (palette >> i) & 0xff
		colorPalette[i/2] = newColor(shade, shade, shade)
	}
	return colorPalette
}

func (lcd *LCD) getBackgroundConf(scanline uint8) (uint16, uint16, bool, uint16) {
	lcdc := lcd.io.Read(io.LDCD)

	bgData := uint16(0x9800)
	if bits.Test(3, lcdc) {
		bgData = 0x9c00
	}

	window := bits.Test(5, lcdc) && scanline >= lcd.io.Read(io.WY)
	windowData := uint16(0x9800)
	if bits.Test(6, lcdc) {
		windowData = 0x9c00
	}

	if bits.Test(4, lcdc) {
		return uint16(0x8000), bgData, window, windowData
	}

	return uint16(0x8800), bgData, window, windowData
}

func (lcd *LCD) drawBackgroundTiles(scanline uint8) {
	tileDataBaseAddr, bgData, window, windowData := lcd.getBackgroundConf(scanline)
	winX := lcd.io.Read(io.WX) - 7
	winY := lcd.io.Read(io.WY)

	colorPalette := lcd.getPalette(io.BGP)

	y := func() uint8 {
		if window {
			return scanline - winY
		}
		return lcd.io.Read(io.SCY) + scanline
	}()
	tileY := uint16(y / 8)
	line := y % 8 * 2
	for i := uint16(0); i < Width; i++ {
		x, tileAddress := func() (uint16, uint16) {
			if window && uint8(i) >= winX {
				return i - uint16(winX), windowData
			}
			return uint16(lcd.io.Read(io.SCX)) + i, bgData
		}()

		tileX := x / 8

		tileAddress += (tileY * 32) + tileX

		// TODO BG WRAP

		tileDataAddress := func() uint16 {
			if tileDataBaseAddr == 0x8800 {
				// TODO signed shenanigans
				// tileNumber := int16(int8(lcd.memory.VRAM[tileAddress-memory.VRAMAddr]))
				// tileDataAddress = uint16(int32(tileDataAddress) + int32)
			}
			// unsigned addressing
			return tileDataBaseAddr + (uint16(lcd.memory.Read(tileAddress)) * 16)
		}()

		tileData := lcd.memory.ReadWord(tileDataAddress + uint16(line))

		bit := uint8((int8(x%8) - 7) * -1)
		shade := ((tileData >> bit & 1) << 1) | (tileData >> (bit + 8) & 1)

		lcd.workData[int(scanline)*Width+int(i)] = colorPalette[shade]
	}
}

func (lcd *LCD) drawScanline(scanline uint8) {
	lcdc := lcd.io.Read(io.LDCD)

	if bits.Test(0, lcdc) {
		lcd.drawBackgroundTiles(scanline)
	}

	if bits.Test(5, lcdc) {
	}
}

func (lcd *LCD) clearScreen() {
	palette := lcd.io.Read(io.BGP)
	shade := palette & 0xff
	for i := 0; i < Width*Height; i++ {
		lcd.workData[i] = newColor(shade, shade, shade)
	}
}

func (lcd *LCD) Update(cycles uint32) {
	lcd.setLCDSTAT()

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
			lcd.io.Write(io.LY, 0)
			copy(lcd.renderData[:], lcd.workData[:])
			lcd.clearScreen()
		} else if ly >= 144 {
			lcd.cpu.RequestInterrupt(cpu.I_VBLANK)
		} else {
			if lcd.lastDrawnScanline != ly {
				// TODO Maybe we should draw at the beginning of H-Blank?
				lcd.drawScanline(ly)
				lcd.lastDrawnScanline = ly
			}
		}
	}
}

func (lcd *LCD) RenderFrame() {
	lcd.backend.Render(lcd.renderData)
}

func NewLCD(cpu *cpu.CPU, mem *memory.Memory, io_ *io.IO) *LCD {
	backend := NewSDL2()
	backend.Initialize(os.Args[0])

	lcd := &LCD{
		backend: backend,

		cpu:    cpu,
		memory: mem,
		io:     io_,

		scanlineCounter:   0,
		lastDrawnScanline: Height,
	}

	lcd.clearScreen()
	copy(lcd.renderData[:], lcd.workData[:])

	return lcd
}
