package cpu

const (
	I_VBLANK  = uint8(0)
	I_LCDSTAT = uint8(1)
	I_TIMER   = uint8(2)
	I_SERIAL  = uint8(3)
	I_JOYPAD  = uint8(4)
)

// Interrupt handler addresses
var (
	interrupts map[uint8]uint16 = map[uint8]uint16{
		0: 0x0040, // V-Blank
		1: 0x0048, // LCD STAT
		2: 0x0050, // Timer
		3: 0x0058, // Serial
		4: 0x0060, // Joypad
	}
)
