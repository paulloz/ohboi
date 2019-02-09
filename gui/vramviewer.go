package gui

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/paulloz/ohboi/gameboy"
)

type vramViewer struct {
	window   *sdl.Window
	windowID uint32

	renderer *sdl.Renderer

	windowRect  sdl.Rect
	previewRect sdl.Rect
	tileRect    sdl.Rect

	tiles []*sdl.Texture

	previewIndex int32

	scale       float32
	nTiles      point
	tileSize    point
	previewSize point
}

func (m *vramViewer) update(gb *gameboy.GameBoy, events []sdl.Event) {
	for _, event := range events {
		switch t := event.(type) {
		case *sdl.MouseMotionEvent:
			if t.WindowID == m.windowID {
				mX := t.X / int32(m.scale)
				mY := t.Y / int32(m.scale)
				if mX <= (m.nTiles.x * m.tileSize.x) {
					x := (mX / m.tileSize.x)
					y := (mY / m.tileSize.y)
					m.previewIndex = x + (y * m.nTiles.x)
				}
			}
		}
	}

	palette := gb.Memory.Read(0xff47)
	c := [4][3]uint8{}
	for i := uint8(0); i < 8; i += 2 {
		shade := (palette >> i) & 0xff
		c[i/2] = [3]uint8{shade, shade, shade}
	}

	sdl.Do(func() {
		tileAddr := uint16(0x8000)
		for tileIndex := range m.tiles {
			buffer, _, err := m.tiles[tileIndex].Lock(nil)
			if err != nil {
				continue
			}
			for i := uint16(0); i < 16; i += 2 {
				data := gb.Memory.ReadWord(tileAddr + uint16(tileIndex*16) + i)

				for j := uint8(0); j < 8; j++ {
					bit := uint8((int8(j%8) - 7) * -1)
					shade := ((data >> bit & 1) << 1) | (data >> (bit + 8) & 1)

					p := int((i/2)*8) + int(j)
					buffer[p*4] = c[shade][0]
					buffer[p*4+1] = c[shade][1]
					buffer[p*4+2] = c[shade][2]
				}
			}
			m.tiles[tileIndex].Unlock()

			tileIndex++
		}
	})
}

func (m *vramViewer) render() {
	m.renderer.Clear()

	m.renderer.SetDrawColor(255, 255, 255, 255)
	m.renderer.DrawRect(&m.windowRect)

	for i, tile := range m.tiles {
		m.renderer.Copy(tile, &m.tileRect, &sdl.Rect{
			X: (int32(i) % m.nTiles.x) * m.tileSize.x,
			Y: (int32(i) / m.nTiles.x) * m.tileSize.y,
			W: m.tileSize.x,
			H: m.tileSize.y,
		})
	}

	m.renderer.Copy(m.tiles[m.previewIndex], &m.tileRect, &m.previewRect)

	m.renderer.Present()
}

func (m *vramViewer) initialize(id int) {
	var err error

	sdl.Do(func() {
		m.windowRect = sdl.Rect{
			X: 0,
			Y: 0,
			W: (m.nTiles.x * m.tileSize.x) + (m.tileSize.x * 2) + m.previewSize.x,
			H: m.nTiles.y * m.tileSize.y,
		}

		m.previewRect = sdl.Rect{
			X: m.windowRect.W - m.tileSize.x - m.previewSize.x,
			Y: m.tileSize.y,
			W: m.previewSize.x,
			H: m.previewSize.y,
		}

		m.tileRect = sdl.Rect{
			X: 0,
			Y: 0,
			W: m.tileSize.x,
			H: m.tileSize.y,
		}
	})

	intScale := int32(m.scale)
	m.window, m.renderer, err = createWindowWithRenderer("VRAMViewer", m.windowRect.W*intScale, m.windowRect.H*intScale)
	if err != nil {
		panic(err)
	}
	m.windowID, err = m.window.GetID()
	if err != nil {
		panic(err)
	}

	m.renderer.SetScale(m.scale, m.scale)

	sdl.Do(func() {
		for i := int32(0); i < m.nTiles.x*m.nTiles.y; i++ {
			tile, err := m.renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, m.tileSize.x, m.tileSize.y)
			if err != nil {
				panic(err)
			}
			m.tiles = append(m.tiles, tile)
		}
	})
}

func (m *vramViewer) destroy() {
	sdl.Do(func() {
		m.renderer.Destroy()
		m.window.Destroy()
	})
}

func newVRAMViewer() *vramViewer {
	return &vramViewer{
		scale:       2.0,
		nTiles:      point{x: 16, y: 18},
		tileSize:    point{x: 8, y: 8},
		previewSize: point{x: 32, y: 32},
	}
}
