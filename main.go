package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paulloz/ohboi/config"
	"github.com/paulloz/ohboi/gameboy"
	"github.com/paulloz/ohboi/gui"
	"github.com/paulloz/ohboi/ppu"
)

var (
	romFilename string
	colorTheme  string
	vramViewer  bool
	skipBoot    bool
	gameBoy     *gameboy.GameBoy
	breakpoint  string
)

func init() {
	// Audio options
	config.Get().Audio.Enabled = *flag.Bool("audio", false, "emulate audio")

	flag.StringVar(&romFilename, "rom", "", "path to the rom file")
	flag.BoolVar(&vramViewer, "vramviewer", false, "enable VRAM viewer")
	flag.BoolVar(&skipBoot, "skipboot", true, "skip boot")
	flag.IntVar(&ppu.Scale, "scale", 2, "scale")
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
		ppu.CurrentPalette = ppu.Greens
	case "sgb":
		ppu.CurrentPalette = ppu.SuperGameboy
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
