package config

import (
	"github.com/paulloz/ohboi/ppu/colors"
)

type audioconfig struct {
	Enabled bool
}

type emulationconfig struct {
	SkipBoot bool
}

type videoconfig struct {
	ColorTheme colors.Palette
	Scale      float64
}

type config struct {
	Audio     audioconfig
	Emulation emulationconfig
	Video     videoconfig
}

var (
	instance *config
)

func Get() *config {
	if instance == nil {
		instance = &config{}
	}

	return instance
}
