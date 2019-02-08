package apu

import (
	"github.com/paulloz/ohboi/io"
)

type APU struct {
	io *io.IO
}

func NewAPU(io *io.IO) *APU {
	return &APU{io: io}
}
