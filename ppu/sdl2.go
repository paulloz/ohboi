package ppu

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/paulloz/ohboi/config"
	"github.com/paulloz/ohboi/consts"
	"github.com/paulloz/ohboi/ppu/colors"
)

type sdl2 struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	texture    *sdl.Texture
	screenRect *sdl.Rect
}

func (sdl2 *sdl2) Render(pixels [consts.ScreenWidth * consts.ScreenHeight]colors.Color) {
	buffer := [len(pixels)]uint32{}
	for i, pixel := range pixels {
		buffer[i] = (uint32(pixel.R) << 24) | (uint32(pixel.G) << 16) | (uint32(pixel.B) << 8) | 0xff
	}
	err := sdl2.texture.UpdateRGBA(sdl2.screenRect, buffer[:], consts.ScreenWidth)
	if err != nil {
		panic(err)
	}

	sdl2.renderer.Clear()
	sdl2.renderer.Copy(sdl2.texture, sdl2.screenRect, sdl2.screenRect)
	sdl2.renderer.Present()
}

func (sdl2 *sdl2) Initialize(windowName string) {
	scale := float32(config.Get().Video.Scale)

	width := int32(consts.ScreenWidth * scale)
	height := int32(consts.ScreenHeight * scale)
	window, err := sdl.CreateWindow(windowName, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer.SetScale(scale, scale)

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, consts.ScreenWidth, consts.ScreenHeight)
	if err != nil {
		panic(err)
	}

	sdl2.window = window
	sdl2.renderer = renderer
	sdl2.texture = texture
	sdl2.screenRect = &sdl.Rect{W: consts.ScreenWidth, H: consts.ScreenHeight}
}

func (sdl2 *sdl2) Destroy() {
	sdl2.texture.Destroy()
	sdl2.renderer.Destroy()
	sdl2.window.Destroy()
}

func NewSDL2() *sdl2 {
	return &sdl2{}
}
