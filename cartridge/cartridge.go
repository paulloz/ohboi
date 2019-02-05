package cartridge

import (
	"fmt"
	"io/ioutil"
)

// MBC ...
type MBC interface {
	Read(address uint16) uint8
	Write(address uint16, data uint8)
}

// Cartridge ...
type Cartridge struct {
	MBC
	Title string
}

// NewCartridge ...
func NewCartridge(filename string) (*Cartridge, error) {
	cartridge := &Cartridge{}

	// Read data from disk
	_data, err := ioutil.ReadFile(filename)
	data := make([]uint8, len(_data))
	if err != nil {
		return nil, err
	}
	for i, b := range _data {
		data[i] = uint8(b)
	}

	// Check cartridge type and init MBC accordingly
	cartridgeType := data[0x0147]
	if cartridgeType == 0x00 {
		cartridge.MBC = &ROM{rom: data}
	} else {
		ramSize := uint16(0)
		switch data[0x0149] {
		case 0x01:
			ramSize = 2 * 1024
		case 0x02:
			ramSize = 8 * 1024
		case 0x03:
			ramSize = 32 * 1024
		}

		if cartridgeType <= 0x03 {
			cartridge.MBC = NewMBC1(data, ramSize)
		} else {
			return nil, fmt.Errorf("unknown cartridge type: %X", cartridgeType)
		}
	}

	// Validate checksum
	x := uint8(0)
	for addr := uint16(0x0134); addr <= 0x014C; addr++ {
		x = x - cartridge.Read(addr) - 1
	}
	if checksum := cartridge.Read(0x014D); x != checksum {
		return nil, fmt.Errorf("could not validate checksum, %X != %X", x, checksum)
	}

	// Title
	cartridge.Title = ""
	for address := uint16(0x0134); address <= 0x143; address++ {
		c := cartridge.Read(address)
		if c == 0 {
			break
		}
		cartridge.Title += string(c)
	}

	return cartridge, nil
}
