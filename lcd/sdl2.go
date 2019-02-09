package lcd

import (
	"github.com/veandco/go-sdl2/sdl"
)

type sdl2 struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	texture    *sdl.Texture
	screenRect *sdl.Rect
}

func (sdl2 *sdl2) Render(pixels [Width * Height]color) {
	buffer, _, err := sdl2.texture.Lock(sdl2.screenRect)
	if err != nil {
		panic(err)
	}

	for i, pixel := range pixels {
		buffer[i*4] = pixel.r
		buffer[i*4+1] = pixel.g
		buffer[i*4+2] = pixel.b
	}

	sdl2.texture.Unlock()

	sdl2.renderer.Clear()
	sdl2.renderer.Copy(sdl2.texture, sdl2.screenRect, sdl2.screenRect)
	sdl2.renderer.Present()
}

func (sdl2 *sdl2) Initialize(windowName string) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(windowName, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, Width*Scale, Height*Scale, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	fScale := float32(Scale)
	renderer.SetScale(fScale, fScale)

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, Width, Height)
	if err != nil {
		panic(err)
	}

	sdl2.window = window
	sdl2.renderer = renderer
	sdl2.texture = texture
	sdl2.screenRect = &sdl.Rect{W: Width, H: Height}
}

func NewSDL2() *sdl2 {
	return &sdl2{}
}
