package main

import (
	"flag"
	"fmt"

	"github.com/paulloz/ohboi/gameboy"
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
	gameBoy = gameboy.NewGameBoy()
	if len(romFilename) > 0 {
		gameBoy.InsertCartridgeFromFile(romFilename)
		gameBoy.PowerOn()
	}
}
