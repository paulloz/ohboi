package main

import (
	"flag"
	"fmt"
	"time"

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
	gameBoy.InsertCartridgeFromFile(romFilename)

	ticker := time.NewTicker(time.Second / gameboy.FPS)

	start := time.Now()
	frames := 0
	for range ticker.C {
		gameBoy.Update()

		frames++
		if time.Since(start) > time.Second {
			start = time.Now()
			frames = 0
		}
	}
}
