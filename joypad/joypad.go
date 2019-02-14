package joypad

import (
	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/cpu"
	"github.com/paulloz/ohboi/io"
)

const (
	M_DIRECTION = uint8(16)
	M_BUTTON    = uint8(32)

	A      = 10
	B      = 11
	START  = 12
	SELECT = 13
	UP     = 21
	RIGHT  = 22
	DOWN   = 23
	LEFT   = 24
)

type backend interface {
	ReadDirections() uint8
	ReadButtons() uint8
}

type Joypad struct {
	cpu *cpu.CPU
	io  *io.IO

	backend backend

	buttons uint8
	mode    uint8
}

func (j *Joypad) ReadInputs() uint8 {
	buttons := uint8(0x0f)

	if (j.mode & M_BUTTON) == 0 {
		buttons = (j.buttons >> 4) & 0x0f
	} else if (j.mode & M_DIRECTION) == 0 {
		buttons = j.buttons & 0x0f
	}

	return 0xc0 | j.mode | buttons
}

func (j *Joypad) WriteInputs(v uint8) {
	j.mode = v & 0x30
}

func (j *Joypad) requestInterrupt() {
	j.cpu.RequestInterrupt(cpu.I_JOYPAD)
}

func (j *Joypad) Update() {
	oldButtons := j.buttons
	j.buttons = (j.backend.ReadButtons() << 4) | j.backend.ReadDirections()

	for i := uint8(0); i < 8; i++ {
		if bits.Test(i, oldButtons) && !bits.Test(i, j.buttons) {
			j.requestInterrupt()
		}
	}
}

func NewJoypad(cpu *cpu.CPU, io_ *io.IO) *Joypad {
	jp := &Joypad{
		cpu: cpu,
		io:  io_,

		backend: NewSDL2(),

		buttons: 0xff,
		mode:    M_BUTTON,
	}

	io_.MapRegister(io.P1, jp.ReadInputs, jp.WriteInputs)

	return jp
}
