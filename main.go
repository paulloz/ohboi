package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/paulloz/ohboi/config"
	"github.com/paulloz/ohboi/gameboy"
	"github.com/paulloz/ohboi/gui"
	"github.com/paulloz/ohboi/ppu/colors"
	"github.com/paulloz/ohboi/statics"
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
	flag.BoolVar(&conf.Audio.Enabled, "audio", conf.Audio.Enabled, "emulate audio")

	// Emulation options
	flag.BoolVar(&conf.Emulation.SkipBoot, "skipboot", conf.Emulation.SkipBoot, "skip boot")
	flag.StringVar(&romFilename, "rom", "", "path to the rom file")

	// Video options
	var theme string
	flag.Float64Var(&conf.Video.Scale, "scale", conf.Video.Scale, "scale")
	flag.StringVar(&theme, "theme", "", "color theme (sgb, green)")
	switch theme {
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
		var err error
		if len(romFilename) > 0 {
			_, err = gameBoy.InsertCartridgeFromPath(romFilename)
		} else if bundled, err := statics.Asset("bundledROM.gb"); err == nil {
			fmt.Println("Found embedded ROM")
			reader := bytes.NewReader([]byte(bundled))
			_, err = gameBoy.InsertCartridge("bundledROM.gb", reader)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}

		gameBoy.PowerOn(quitChan)
		quitChan <- 0
	}()

	guiOptions := gui.GUIOptions{VRAMViewer: vramViewer}
	os.Exit(gui.GUIStart(guiOptions, gameBoy, quitChan))
}
