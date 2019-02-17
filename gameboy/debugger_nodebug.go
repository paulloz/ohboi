// +build !DEBUG

package gameboy

func debuggerStart(gb *GameBoy) {}

func debuggerStop() {}

func debuggerStep() {}

func AddBreakpoint(s string) {
}

func StepByStep(state bool) {
}
