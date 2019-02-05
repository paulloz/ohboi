package main

import (
	"flag"
	"fmt"
	"os"

	. "github.com/paulloz/ohboi/cartridge"
)

var (
	filename string
)

func init() {
	flag.StringVar(&filename, "rom", "", "path to the rom file")
	flag.Parse()
}

func main() {
	cartridge, err := NewCartridge(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(cartridge.Title)
}
