// +build !android

package joypad

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/paulloz/ohboi/bits"
)

type sdl2 struct {
	mapping map[int]int
}

func (s *sdl2) ReadDirections() uint8 {
	keys := sdl.GetKeyboardState()

	x := uint8(0xff & 0xf)

	if keys[s.mapping[DOWN]] == 1 {
		x = bits.Reset(3, x)
	}
	if keys[s.mapping[UP]] == 1 {
		x = bits.Reset(2, x)
	}
	if keys[s.mapping[LEFT]] == 1 {
		x = bits.Reset(1, x)
	}
	if keys[s.mapping[RIGHT]] == 1 {
		x = bits.Reset(0, x)
	}

	return x
}

func (s *sdl2) ReadButtons() uint8 {
	keys := sdl.GetKeyboardState()

	x := uint8(0xff & 0xf)

	if keys[s.mapping[START]] == 1 {
		x = bits.Reset(3, x)
	}
	if keys[s.mapping[SELECT]] == 1 {
		x = bits.Reset(2, x)
	}
	if keys[s.mapping[B]] == 1 {
		x = bits.Reset(1, x)
	}
	if keys[s.mapping[A]] == 1 {
		x = bits.Reset(0, x)
	}

	return x
}

func NewBackend() *sdl2 {
	return &sdl2{
		mapping: map[int]int{
			A:      sdl.SCANCODE_K,
			B:      sdl.SCANCODE_J,
			START:  sdl.SCANCODE_SPACE,
			SELECT: sdl.SCANCODE_ESCAPE,
			UP:     sdl.SCANCODE_W,
			RIGHT:  sdl.SCANCODE_D,
			DOWN:   sdl.SCANCODE_S,
			LEFT:   sdl.SCANCODE_A,
		},
	}
}
