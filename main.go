package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paulloz/ohboi/gameboy"
	"github.com/paulloz/ohboi/gui"
)

var (
	romFilename string
	gameBoy     *gameboy.GameBoy
)

func init() {
	flag.StringVar(&romFilename, "rom", "", "path to the rom file")
	flag.Parse()

	if len(romFilename) <= 0 {
		fmt.Println("No cardbridge inserted...")
	}
}

func main() {
	quitChan := make(chan int)

	gameBoy = gameboy.NewGameBoy(true)

	go func() {
		if len(romFilename) > 0 {
			gameBoy.InsertCartridgeFromFile(romFilename)
			gameBoy.PowerOn(quitChan)
		}
		quitChan <- 0
	}()

	guiOptions := gui.GUIOptions{VRAMViewer: false}
	os.Exit(gui.GUIStart(guiOptions, gameBoy, quitChan))
}
