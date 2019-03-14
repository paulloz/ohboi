// +build android

package ppu

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/paulloz/ohboi/consts"
	"github.com/paulloz/ohboi/ppu/colors"
	"golang.org/x/mobile/exp/gl/glutil"
)

var Screen *glutil.Image

type gomobile struct {
}

func (g *gomobile) Render(pixels [consts.ScreenWidth * consts.ScreenHeight]colors.Color) {
	if Screen != nil {
		draw.Draw(Screen.RGBA, Screen.RGBA.Bounds(), image.Black, image.Point{}, draw.Src)
		for x := 0; x < consts.ScreenWidth; x++ {
			for y := 0; y < consts.ScreenHeight; y++ {
				pixel := pixels[y*consts.ScreenWidth+x]
				Screen.RGBA.SetRGBA(x, y, color.RGBA{R: pixel.R, G: pixel.G, B: pixel.B, A: 255})
			}
		}
	}
}

func (g *gomobile) Initialize(windowName string) {
}

func (g *gomobile) Destroy() {
}

func NewBackend() *gomobile {
	return &gomobile{}
}
