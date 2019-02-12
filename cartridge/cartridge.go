package cartridge

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

// MBC interface describes a memory bank controller
type MBC interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}

// Header describes the GameBoy cartBridge header
type Header struct {
	Title                 string
	Designation           string
	ColorCompatibility    uint8
	SGBColorCompatibility uint8
	Type                  uint8
	NewLicenseeCode       uint16
	ROMSize               uint8
	RAMSize               uint8
	DestinationCode       uint8
	OldLicenseeCode       uint8
	MaskROMVersion        uint8
	HeaderChecksum        uint8
	GlobalChecksum        uint16
}

// Cartridge describes a GameBoy cartridge
type Cartridge struct {
	MBC
	Header
}

func (h *Header) readHeader(data []byte) error {
	h.HeaderChecksum = data[0x4d]
	if err := h.validateChecksum(data); err != nil {
		return err
	}

	h.Title = string(bytes.Trim(data[0x34:0x3e], "\x00"))
	h.Designation = string(bytes.Trim(data[0x3f:0x42], "\x00"))
	h.ColorCompatibility = data[0x43]
	h.NewLicenseeCode = binary.BigEndian.Uint16(data[0x44:0x46])
	h.SGBColorCompatibility = data[0x46]
	h.Type = data[0x47]
	h.ROMSize = data[0x48]
	h.RAMSize = data[0x49]
	h.DestinationCode = data[0x4a]
	h.OldLicenseeCode = data[0x4b]
	h.MaskROMVersion = data[0x4c]
	h.GlobalChecksum = binary.BigEndian.Uint16(data[0x4e:0x50])

	return nil
}

func (h *Header) validateChecksum(data []byte) error {
	x := uint8(0)
	for addr := uint16(0x34); addr <= 0x4C; addr++ {
		x = x - data[addr] - 1
	}

	if checksum := data[0x4d]; x != checksum {
		return fmt.Errorf("could not validate checksum, %X != %X", x, checksum)
	}

	return nil
}

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

	if err := cartridge.readHeader(data[0x100:0x150]); err != nil {
		return nil, err
	}

	// Check cartridge type and init MBC accordingly
	switch cartridge.Type {
	case ROMOnlyType:
		cartridge.MBC = &ROM{rom: data}
	case MBC1Type, MBC1RAMType, MBC1RAMBatteryType:
		ramSize := uint16(0)
		switch cartridge.RAMSize {
		case RAMNone:
			ramSize = 2 * 1024
		case RAM2k:
			ramSize = 2 * 1024
		case RAM8k:
			ramSize = 8 * 1024
		case RAM32k:
			ramSize = 32 * 1024
		default:
			return nil, fmt.Errorf("Unsupported RAM size: %d\n", cartridge.RAMSize)
		}
		cartridge.MBC = NewMBC1(data, ramSize)
	default:
		return nil, fmt.Errorf("Unsupported cartridge type: %X", cartridge.Type)
	}

	return cartridge, nil
}
