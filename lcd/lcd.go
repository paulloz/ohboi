package lcd

import (
	"os"
	"sort"

	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/consts"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/memory"
)

var (
	Scale = 2
)

const (
	Width  = 160
	Height = 144

	ScanlineFrequency   = 456
	SpritesCount        = 40
	MaxDisplayedSprites = 10
)

type backend interface {
	Initialize(string)
	Render([consts.ScreenWidth * consts.ScreenHeight]Color)
	Destroy()
}

type LCD struct {
	backend backend

	cpu    *cpu.CPU
	memory *memory.Memory
	io     *io.IO

	scanlineCounter   uint32
	lastDrawnScanline uint8

	workData   [consts.ScreenWidth * consts.ScreenHeight]Color
	renderData [consts.ScreenWidth * consts.ScreenHeight]Color
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

	return uint16(0x9000), bgData, window, windowData
}

type Sprite struct {
	X, Y    uint8
	Pattern uint8
	Flags   uint8
}

type SpriteList [40]Sprite

func (l *SpriteList) Len() int {
	return 40
}

func (l *SpriteList) Swap(i, j int) {
	tmp := l[i]
	l[i] = l[j]
	l[j] = tmp
}

func (l *SpriteList) Less(i, j int) bool {
	switch {
	case l[i].X < l[j].X:
		return true
	case l[i].X == l[j].X:
		return i < j
	default:
		return false
	}
}
func (lcd *LCD) drawSprites(scanline uint8) {
	var palette [4]Color
	var sprites SpriteList

	for i := uint16(0); i < SpritesCount; i++ {
		sprites[i] = Sprite{
			Y:       lcd.memory.Read(memory.OAMAddr + i*4),
			X:       lcd.memory.Read(memory.OAMAddr + i*4 + 1),
			Pattern: lcd.memory.Read(memory.OAMAddr + i*4 + 2),
			Flags:   lcd.memory.Read(memory.OAMAddr + i*4 + 3),
		}
	}
	sort.Sort(&sprites)

	spriteHeight := uint8(8)
	patternMask := uint8(0xff)
	if lcd.io.ReadBit(io.LDCD, 2) {
		spriteHeight = 16
		patternMask = 0xfe
	}

	x, displayed := int16(0), 0
	for i := uint16(0); i < SpritesCount && displayed < MaxDisplayedSprites; i++ {
		sprite := &sprites[i]
		spriteY := int16(sprite.Y) - 16

		if (sprite.X == 0 && sprite.Y == 0) ||
			(int16(scanline) < spriteY) ||
			(int16(scanline) >= spriteY+int16(spriteHeight)) {
			continue
		}

		if x = int16(sprite.X - 8); x >= consts.ScreenWidth {
			break
		}

		if bits.Test(4, sprite.Flags) {
			palette = lcd.getPalette(io.OBP1)
		} else {
			palette = lcd.getPalette(io.OBP0)
		}

		tileDataAddress := memory.VRAMAddr + uint16(sprite.Pattern&patternMask)*16
		flipX := bits.Test(5, sprite.Flags)
		flipY := bits.Test(6, sprite.Flags)

		line := scanline - (sprite.Y - 16)
		if flipY {
			line = spriteHeight - line - 1
		}

		for j := uint8(0); j < 8 && x >= 0 && x < consts.ScreenWidth; j++ {
			var bit uint8
			if flipX {
				bit = j
			} else {
				bit = 7 - j
			}

			tileData := lcd.memory.ReadWord(tileDataAddress + (uint16(line) % 8 * 2))
			shade := ((tileData >> bit & 1) << 1) | (tileData >> (bit + 8) & 1)
			if shade > 0 {
				index := (int(scanline) * consts.ScreenWidth) + int(x)
				lcd.workData[index] = palette[shade]
			}
			x++
		}

		displayed++
	}
}

func (lcd *LCD) drawBackgroundTiles(scanline uint8) {
	tileDataBaseAddr, bgData, window, windowData := lcd.getBackgroundConf(scanline)
	winX := lcd.io.Read(io.WX) - 7
	winY := lcd.io.Read(io.WY)
	colorPalette := lcd.getPalette(io.BGP)

	for i := uint16(0); i < consts.ScreenWidth; i++ {
		var x, y, tileAddress, tileDataAddress uint16
		var tileNumber int16

		if window && uint8(i) >= winX {
			x, tileAddress = i-uint16(winX), windowData
		} else {
			x, tileAddress = uint16(lcd.io.Read(io.SCX))+i, bgData
		}

		if window {
			y = uint16(scanline - winY)
		} else {
			y = uint16(lcd.io.Read(io.SCY) + scanline)
		}

		tileX := x / 8
		tileY := uint16(y / 8)
		line := y % 8 * 2

		tileAddress += (tileY * 32) + tileX

		// TODO BG WRAP

		if tileDataBaseAddr == 0x9000 {
			// signed addressing
			tileNumber = int16(int8(lcd.memory.Read(tileAddress)))
		} else {
			// unsigned addressing
			tileNumber = int16(lcd.memory.Read(tileAddress))
		}

		tileDataAddress = uint16(int32(tileDataBaseAddr) + (int32(tileNumber) * 16))
		tileData := lcd.memory.ReadWord(tileDataAddress + uint16(line))

		bit := uint8((int8(x%8) - 7) * -1)
		shade := ((tileData >> bit & 1) << 1) | (tileData >> (bit + 8) & 1)

		index := (int(scanline) * consts.ScreenWidth) + int(i)
		lcd.workData[index] = colorPalette[shade]
	}
}

func (lcd *LCD) drawScanline(scanline uint8) {
	lcdc := lcd.io.Read(io.LDCD)

	if bits.Test(0, lcdc) {
		lcd.drawBackgroundTiles(scanline)
	}

	if bits.Test(1, lcdc) {
		lcd.drawSprites(scanline)
	}
}

func (lcd *LCD) clearScreen() {
	for i := 0; i < (consts.ScreenWidth * consts.ScreenHeight); i++ {
		lcd.workData[i] = Greys[0]
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
		lcd.scanlineCounter -= ScanlineFrequency

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

func (lcd *LCD) Destroy() {
	lcd.backend.Destroy()
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
		lastDrawnScanline: consts.ScreenHeight,
	}

	lcd.clearScreen()
	copy(lcd.renderData[:], lcd.workData[:])

	return lcd
}
