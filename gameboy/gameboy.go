package gameboy

// Speed constants
const (
	ClockSpeed     = 4194304          // Cycles per second
	FPS            = 60               // We want to run at 60FPS
	CyclesPerFrame = ClockSpeed / FPS // # of cycles executed each frame
)

// GameBoy ...
type GameBoy struct {
}

// ExecuteNextOpCode ...
func (gb *GameBoy) ExecuteNextOpCode() uint {
	return 1
}

// Update ...
func (gb *GameBoy) Update() uint {
	var cycles uint

	for cycles = 0; cycles < CyclesPerFrame; {
		_cycles := gb.ExecuteNextOpCode()

		// UpdateTimers
		// UpdateGraphics
		// Interrupts

		cycles += _cycles
	}

	// RenderScreen

	return cycles
}

// NewGameBoy ...
func NewGameBoy() *GameBoy {
	return &GameBoy{}
}
