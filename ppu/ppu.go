package ppu

import (
	"os"
	"sort"

	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/config"
	"github.com/paulloz/ohboi/consts"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
	"github.com/paulloz/ohboi/memory"
	"github.com/paulloz/ohboi/ppu/colors"
)

var (
	Scale = 2
)

const (
	ScanlineFrequency   = 456
	SpritesCount        = 40
	MaxDisplayedSprites = 10
)

type backend interface {
	Initialize(string)
	Render([consts.ScreenWidth * consts.ScreenHeight]colors.Color)
	Destroy()
}

type PPU struct {
	backend backend

	cpu    *cpu.CPU
	memory *memory.Memory
	io     *io.IO

	scanlineCounter   uint32
	lastDrawnScanline uint8

	workData   [consts.ScreenWidth * consts.ScreenHeight]colors.Color
	renderData [consts.ScreenWidth * consts.ScreenHeight]colors.Color
	pixels     [consts.ScreenWidth * consts.ScreenHeight]uint8
}

func (ppu *PPU) setLCDSTAT() {
	if !ppu.io.ReadBit(io.LDCD, 7) {
		// Reset everything
		ppu.scanlineCounter = 0
		ppu.io.Write(io.LY, 0)
		ppu.io.Write(io.STAT, 1<<2)
		return
	}

	ly := ppu.io.Read(io.LY)

	stat := ppu.io.Read(io.STAT)
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
			ppu.cpu.RequestInterrupt(cpu.I_LCDSTAT)
		}
	} else {
		// The end of V-Blank to begin of next V-Blank period takes 456 CPU cycles
		// It is timed like this: Mode2 (80 cycles) -> Mode3 (170~240 cycles) -> Mode0 (remaining cycles)
		if ppu.scanlineCounter < 80 {
			// We're in Mode2, set mode bits to 11
			stat = bits.Set(1, bits.Set(0, stat))
			if bits.Test(5, stat) && modeChanged() {
				// We changed mode and Mode2 interrupt is enabled
				ppu.cpu.RequestInterrupt(cpu.I_LCDSTAT)
			}
		} else if ppu.scanlineCounter < 80+170 {
			// We're in Mode3, set mode bits to 10, no interrupt for Mode3
			stat = bits.Set(1, bits.Reset(0, stat))
		} else {
			// We're in Mode0, set mode bits to 00
			stat = bits.Reset(1, bits.Reset(0, stat))
			if bits.Test(3, stat) && modeChanged() {
				// We changed mode and Mode0 interrupt is enabled
				ppu.cpu.RequestInterrupt(cpu.I_LCDSTAT)
			}
		}
	}

	// if LYC == LYC, must set bit 2 and request interrupt if bit 6 is set
	// must reset bit 2 otherwise
	lyc := ppu.io.Read(io.LYC)
	if lyc == ly {
		bits.Set(2, stat)
		if bits.Test(6, stat) {
			ppu.cpu.RequestInterrupt(cpu.I_LCDSTAT)
		}
	} else {
		bits.Reset(2, stat)
	}

	ppu.io.Write(io.STAT, stat)
}

func (ppu *PPU) getBackgroundConf(scanline uint8) (uint16, uint16, bool, uint16) {
	lcdc := ppu.io.Read(io.LDCD)

	bgData := uint16(0x9800)
	if bits.Test(3, lcdc) {
		bgData = 0x9c00
	}

	window := bits.Test(5, lcdc) && scanline >= ppu.io.Read(io.WY)
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
func (ppu *PPU) drawSprites(scanline uint8) {
	var palette colors.Palette
	var sprites SpriteList

	for i := uint16(0); i < SpritesCount; i++ {
		sprites[i] = Sprite{
			Y:       ppu.memory.Read(memory.OAMAddr + i*4),
			X:       ppu.memory.Read(memory.OAMAddr + i*4 + 1),
			Pattern: ppu.memory.Read(memory.OAMAddr + i*4 + 2),
			Flags:   ppu.memory.Read(memory.OAMAddr + i*4 + 3),
		}
	}
	sort.Sort(&sprites)

	spriteHeight := uint8(8)
	patternMask := uint8(0xff)
	if ppu.io.ReadBit(io.LDCD, 2) {
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
			palette = ppu.getPalette(io.OBP1)
		} else {
			palette = ppu.getPalette(io.OBP0)
		}

		tileDataAddress := memory.VRAMAddr + uint16(sprite.Pattern&patternMask)*16
		priority := bits.Test(7, sprite.Flags)
		flipY := bits.Test(6, sprite.Flags)
		flipX := bits.Test(5, sprite.Flags)

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

			tileData := ppu.memory.ReadWord(tileDataAddress + (uint16(line) * 2))
			shade := ((tileData >> bit & 1) | (tileData>>(bit+8)&1)<<1)

			if shade != 0 && (!priority || ppu.pixels[int(scanline)*consts.ScreenWidth+int(x)] == 0) {
				offset := int(scanline)*consts.ScreenWidth + int(x)
				ppu.workData[offset] = palette[shade]
			}

			x++
		}

		displayed++
	}
}

func (ppu *PPU) drawBackgroundTiles(scanline uint8) {
	tileDataBaseAddr, bgData, window, windowData := ppu.getBackgroundConf(scanline)
	winX := ppu.io.Read(io.WX) - 7
	winY := ppu.io.Read(io.WY)
	colorPalette := ppu.getPalette(io.BGP)

	for i := uint16(0); i < consts.ScreenWidth; i++ {
		var x, y, tileAddress, tileDataAddress uint16
		var tileNumber int16

		if window && uint8(i) >= winX {
			x, tileAddress = i-uint16(winX), windowData
		} else {
			x, tileAddress = (uint16(ppu.io.Read(io.SCX))+i)%256, bgData
		}

		if window {
			y = uint16(scanline - winY)
		} else {
			y = uint16(ppu.io.Read(io.SCY) + scanline)
		}

		tileX := x / 8
		tileY := uint16(y / 8)
		line := y % 8 * 2

		tileAddress += (tileY * 32) + tileX

		// TODO: BG WRAP

		if tileDataBaseAddr == 0x9000 {
			// signed addressing
			tileNumber = int16(int8(ppu.memory.Read(tileAddress)))
		} else {
			// unsigned addressing
			tileNumber = int16(ppu.memory.Read(tileAddress))
		}

		tileDataAddress = uint16(int32(tileDataBaseAddr) + (int32(tileNumber) * 16))
		tileData := ppu.memory.ReadWord(tileDataAddress + uint16(line))

		bit := uint8((int8(x%8) - 7) * -1)
		shade := (tileData >> bit & 1) | ((tileData >> (bit + 8) & 1) << 1)

		ppu.workData[int(scanline)*consts.ScreenWidth+int(i)] = colorPalette[shade]
		ppu.pixels[int(scanline)*consts.ScreenWidth+int(i)] = uint8(shade)
	}
}

func (ppu *PPU) drawScanline(scanline uint8) {
	lcdc := ppu.io.Read(io.LDCD)

	if bits.Test(0, lcdc) {
		ppu.drawBackgroundTiles(scanline)
	}

	if bits.Test(1, lcdc) {
		ppu.drawSprites(scanline)
	}
}

func (ppu *PPU) getPalette(ioAddr uint8) colors.Palette {
	palette := ppu.io.Read(ioAddr)
	colorPalette := colors.Palette{}
	for i := uint8(0); i < 8; i += 2 {
		shade := (palette >> i) & 3
		colorPalette[i/2] = config.Get().Video.ColorTheme[shade]
	}
	return colorPalette
}

func (ppu *PPU) clearScreen() {
	for i := 0; i < (consts.ScreenWidth * consts.ScreenHeight); i++ {
		ppu.workData[i] = config.Get().Video.ColorTheme[0]
	}
}

func (ppu *PPU) Update(cycles uint32) {
	ppu.setLCDSTAT()

	if !ppu.io.ReadBit(io.LDCD, 7) {
		// If LCD is disabled
		return
	}

	ppu.scanlineCounter += cycles
	if ppu.scanlineCounter >= ScanlineFrequency {
		ppu.scanlineCounter -= ScanlineFrequency

		ly := ppu.io.Read(io.LY) + 1
		ppu.io.Write(io.LY, ly)

		if ly > 153 {
			ppu.io.Write(io.LY, 0)
			copy(ppu.renderData[:], ppu.workData[:])
			ppu.clearScreen()
		} else if ly >= 144 {
			ppu.cpu.RequestInterrupt(cpu.I_VBLANK)
		} else {
			if ppu.lastDrawnScanline != ly {
				// TODO: Maybe we should draw at the beginning of H-Blank?
				ppu.drawScanline(ly)
				ppu.lastDrawnScanline = ly
			}
		}
	}
}

func (ppu *PPU) RenderFrame() {
	ppu.backend.Render(ppu.renderData)
}

func (ppu *PPU) Destroy() {
	ppu.backend.Destroy()
}

func NewPPU(cpu *cpu.CPU, mem *memory.Memory, io_ *io.IO) *PPU {
	backend := NewSDL2()
	backend.Initialize(os.Args[0])

	ppu := &PPU{
		backend: backend,

		cpu:    cpu,
		memory: mem,
		io:     io_,

		scanlineCounter:   0,
		lastDrawnScanline: consts.ScreenHeight,
	}

	ppu.clearScreen()
	copy(ppu.renderData[:], ppu.workData[:])

	return ppu
}
