// +build android

package gui

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"

	"github.com/paulloz/ohboi/consts"
	"github.com/paulloz/ohboi/gameboy"
	"github.com/paulloz/ohboi/ppu"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32
)

func GUIStart(options GUIOptions, gb *gameboy.GameBoy, quitChan chan int) int {
	exitCode := 0

	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)

				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
	return exitCode
}

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	ppu.Screen = images.NewImage(consts.ScreenWidth, consts.ScreenHeight)
	fps = debug.NewFPS(images)
}

func onStop(glctx gl.Context) {
	fps.Release()
	ppu.Screen.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	ppu.Screen.Upload()
	ppu.Screen.Draw(
		sz,
		geom.Point{0, 0},
		geom.Point{sz.WidthPt, 0},
		geom.Point{0, sz.HeightPt},
		ppu.Screen.RGBA.Bounds(),
	)
	fps.Draw(sz)
}
