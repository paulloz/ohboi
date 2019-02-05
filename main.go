package main

import (
	"flag"
	"os"

)

var (
	romFilename string
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

}
