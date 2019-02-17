package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paulloz/ohboi/apu"
	"github.com/paulloz/ohboi/gameboy"
	"github.com/paulloz/ohboi/gui"
	"github.com/paulloz/ohboi/lcd"
)

var (
	romFilename string
	colorTheme  string
	vramViewer  bool
	skipBoot    bool
	audio       bool
	gameBoy     *gameboy.GameBoy
	breakpoint  string
)

func init() {
	flag.StringVar(&romFilename, "rom", "", "path to the rom file")
	flag.BoolVar(&vramViewer, "vramviewer", false, "enable VRAM viewer")
	flag.BoolVar(&skipBoot, "skipboot", true, "skip boot")
	flag.BoolVar(&audio, "audio", false, "emulate audio")
	flag.IntVar(&lcd.Scale, "scale", 2, "scale")
	flag.StringVar(&colorTheme, "theme", "green", "color theme (grey, green)")
	flag.StringVar(&breakpoint, "breakpoint", "", "breakpoint")
	flag.Parse()

	if len(romFilename) <= 0 {
		fmt.Println("No cardbridge inserted...")
	}
}

type Const struct {
	value uint8
}

func (c *Const) Get() uint8 {
	return c.value
}

func main() {
	quitChan := make(chan int)

	switch colorTheme {
	case "green":
		lcd.CurrentPalette = lcd.Greens
	case "sgb":
		lcd.CurrentPalette = lcd.SuperGameboy
	}

	if audio {
		apu.Backend = "sdl2"
	} else {
		apu.Backend = ""
	}

	gameBoy = gameboy.NewGameBoy(skipBoot)

	if breakpoint != "" {
		gameboy.AddBreakpoint(breakpoint)
		gameboy.StepByStep(false)
	}

	go func() {
		if len(romFilename) > 0 {
			gameBoy.InsertCartridgeFromFile(romFilename)
			gameBoy.PowerOn(quitChan)
		}
		quitChan <- 0
	}()

	guiOptions := gui.GUIOptions{VRAMViewer: vramViewer}
	os.Exit(gui.GUIStart(guiOptions, gameBoy, quitChan))
}
