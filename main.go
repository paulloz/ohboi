package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paulloz/ohboi/gameboy"
	"github.com/paulloz/ohboi/gui"
	"github.com/paulloz/ohboi/lcd"
)

var (
	romFilename string
	colorTheme  string
	vramViewer  bool
	skipBoot    bool
	gameBoy     *gameboy.GameBoy
)

func init() {
	flag.StringVar(&romFilename, "rom", "", "path to the rom file")
	flag.BoolVar(&vramViewer, "vramViewer", false, "enable VRAM viewer")
	flag.BoolVar(&skipBoot, "skipBoot", true, "skip boot")
	flag.IntVar(&lcd.Scale, "scale", 2, "scale")
	flag.StringVar(&colorTheme, "theme", "green", "color theme (grey, green)")
	flag.Parse()

	if len(romFilename) <= 0 {
		fmt.Println("No cardbridge inserted...")
	}
}

func main() {
	quitChan := make(chan int)

	switch colorTheme {
	case "green":
		lcd.CurrentPalette = lcd.Greens
	case "sgb":
		lcd.CurrentPalette = lcd.SuperGameboy
	}

	gameBoy = gameboy.NewGameBoy(skipBoot)

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
