package gui

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/paulloz/ohboi/gameboy"
)

type point struct {
	x int32
	y int32
}

type guiModule interface {
	initialize(int)
	destroy()
	update(*gameboy.GameBoy, []sdl.Event)
	render()
}

type GUIOptions struct {
	VRAMViewer bool
}

var (
	modules []guiModule
)

func createWindow(title string, w, h int32) (*sdl.Window, error) {
	var window *sdl.Window
	var err error

	sdl.Do(func() {
		window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	})

	return window, err
}

func createWindowWithRenderer(title string, w, h int32) (*sdl.Window, *sdl.Renderer, error) {
	window, err := createWindow(title, w, h)
	if err != nil {
		return window, nil, err
	}

	var renderer *sdl.Renderer
	sdl.Do(func() {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	})

	return window, renderer, err
}

func addModule(m guiModule) {
	modules = append(modules, m)
	id := len(modules) - 1
	modules[id].initialize(id)
}

func guiRun(options GUIOptions, gb *gameboy.GameBoy, quitChan chan int) int {
	if options.VRAMViewer {
		addModule(newVRAMViewer())
	}

	ticker := time.NewTicker(time.Second / 10).C
	for {
		select {
		case <-ticker:
			var events []sdl.Event
			sdl.Do(func() {
				for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
					events = append(events, event)

					switch event.(type) {
					case *sdl.QuitEvent:
						quitChan <- 0
					}
				}
			})

			for _, module := range modules {
				module.update(gb, events)
			}
			sdl.Do(func() {
				for _, module := range modules {
					module.render()
				}
			})
		case exitCode := <-quitChan:
			for _, module := range modules {
				module.destroy()
			}
			return exitCode
		}
	}
}

func GUIStart(options GUIOptions, gb *gameboy.GameBoy, quitChan chan int) int {
	exitCode := 0

	sdl.Main(func() {
		exitCode = guiRun(options, gb, quitChan)
	})

	return exitCode
}

func init() {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS | sdl.INIT_TIMER | sdl.INIT_AUDIO)
	if err != nil {
		panic(err)
	}
}
