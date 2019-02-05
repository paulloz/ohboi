package main

import (
	"flag"
	"os"
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
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	gameBoy = gameboy.NewGameBoy()

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
