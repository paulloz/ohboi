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
	conf := config.Get()

	// Audio options
	conf.Audio.Enabled = *flag.Bool("audio", conf.Audio.Enabled, "emulate audio")

	// Emulation options
	conf.Emulation.SkipBoot = *flag.Bool("skipboot", conf.Emulation.SkipBoot, "skip boot")
	flag.StringVar(&romFilename, "rom", "", "path to the rom file")

	// Video options
	conf.Video.Scale = *flag.Float64("scale", conf.Video.Scale, "scale")
	switch *flag.String("theme", "", "color theme (sgb, green)") {
	case "green":
		conf.Video.ColorTheme = colors.Greens
	case "sgb":
		conf.Video.ColorTheme = colors.SuperGameboy
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
