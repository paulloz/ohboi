package cpu

// Interrupt handler addresses
var (
	Interrupts map[uint8]uint16 = map[uint8]uint16{
		0: 0x0040, // V-Blank
		1: 0x0048, // LCD STAT
		2: 0x0050, // Timer
		3: 0x0058, // Serial
		4: 0x0060, // Joypad
	}
)
