package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paulloz/ohboi/config"
	"github.com/paulloz/ohboi/gameboy"
	"github.com/paulloz/ohboi/gui"
	"github.com/paulloz/ohboi/ppu/colors"
)

var (
	romFilename string
	vramViewer  bool
	gameBoy     *gameboy.GameBoy
	breakpoint  string
)

func init() {
	// Audio options
	config.Get().Audio.Enabled = *flag.Bool("audio", false, "emulate audio")

	// Emulation options
	config.Get().Emulation.SkipBoot = *flag.Bool("skipboot", true, "skip boot")
	flag.StringVar(&romFilename, "rom", "", "path to the rom file")

	// Video options
	config.Get().Video.Scale = *flag.Float64("scale", 2, "scale")
	switch *flag.String("theme", "green", "color theme (grey, green)") {
	case "green":
		config.Get().Video.ColorTheme = colors.Greens
	case "sgb":
		config.Get().Video.ColorTheme = colors.SuperGameboy
	}

	flag.BoolVar(&vramViewer, "vramviewer", false, "enable VRAM viewer")
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

	gameBoy = gameboy.NewGameBoy()

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
