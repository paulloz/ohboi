package cartridge

import (
	"errors"
	"fmt"
	"io/ioutil"

	. "github.com/paulloz/gb-emulator/types"
)

type MBC interface {
	Read(address Word) Byte
}

type Cartridge struct {
	MBC
	Title string
}

func NewCartridge(filename string) (*Cartridge, error) {
	cartridge := &Cartridge{}

	// Read data from disk
	_data, err := ioutil.ReadFile(filename)
	data := make([]Byte, len(_data))
	if err != nil {
		return nil, err
	}
	for i, b := range _data {
		data[i] = Byte(b)
	}

	// Check cartridge type and init MBC accordingly
	cartridgeType := data[0x0147]
	if cartridgeType == 0x00 {
		cartridge.MBC = &ROM{rom: data}
	} else {
		if cartridgeType <= 0x03 {
			cartridge.MBC = &MBC1{rom: data}
		} else {
			return nil, errors.New(fmt.Sprintf("Unknown cartridge type: %X.", cartridgeType))
		}
	}

	// Validate checksum
	x := Byte(0)
	for addr := Word(0x0134); addr <= 0x014C; addr++ {
		x = x - cartridge.Read(addr) - 1
	}
	if checksum := cartridge.Read(0x014D); x != checksum {
		return nil, errors.New(fmt.Sprintf("Could not validate checksum, %X != %X.", x, checksum))
	}

	// Title
	cartridge.Title = ""
	for address := Word(0x0134); address <= 0x143; address++ {
		c := cartridge.Read(address)
		if c == 0 {
			break
		}
		cartridge.Title += string(c)
	}

	return cartridge, nil
}
