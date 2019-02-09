package lcd

import (
	"github.com/veandco/go-sdl2/sdl"
)

type sdl2 struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func (sdl2 *sdl2) Render(pixels [Width * Height]color) {
	// Test
	for i, pixel := range pixels {
		sdl2.renderer.SetDrawColor(pixel.r, pixel.g, pixel.b, 255)
		sdl2.renderer.DrawPoint(int32(i%Width), int32(i/Width))
	}
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

	sdl2.window = window
	sdl2.renderer = renderer
}

func NewSDL2() *sdl2 {
	return &sdl2{}
}
