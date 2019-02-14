package consts

const (
	FPS = 60 // We want to run at 60FPS

	CPUClockSpeed     = uint32(4194304)     // Cycles per second
	CPUCyclesPerFrame = CPUClockSpeed / FPS // Cycles per frame

	APUSampleRate         = 44032
	CPUCyclesPerAPUSample = CPUClockSpeed / APUSampleRate

	ScreenWidth  = 160
	ScreenHeight = 144

	RenderScale = 2
)
