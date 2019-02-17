package apu

import (
	"github.com/paulloz/ohboi/bits"
	"github.com/paulloz/ohboi/consts"
	"github.com/paulloz/ohboi/io"
)

type backend interface {
	Output([BufferSize * 2]byte)
	Destroy()
}

var (
	Backend string
)

const (
	BufferSize = 1024
)

type APU struct {
	io *io.IO

	backend backend

	cycles uint32

	buffer         [BufferSize * 2]byte
	bufferIterator int

	channels []channel
	active   bool

	volume struct {
		left  float64
		right float64
	}

	nr51 uint8
}

func (apu *APU) Update(cycles uint32) {
	apu.cycles += cycles

	if apu.cycles >= consts.CPUCyclesPerAPUSample {
		// SAMPLE TIME!

		apu.cycles -= consts.CPUCyclesPerAPUSample

		var valueLeft uint16
		var valueRight uint16

		for _, channel := range apu.channels {
			value := channel.Sample()
			if channel.IsActiveLeft() {
				valueLeft += value
			}
			if channel.IsActiveRight() {
				valueRight += value
			}
		}

		valueLeft /= uint16(len(apu.channels))
		valueRight /= uint16(len(apu.channels))

		apu.buffer[apu.bufferIterator] = byte(float64(valueLeft) * apu.volume.left)         // left
		apu.buffer[(apu.bufferIterator + 1)] = byte(float64(valueRight) * apu.volume.right) // right

		apu.bufferIterator += 2
		if apu.bufferIterator >= (BufferSize * 2) {
			// Buffer is full, send it and start over

			apu.backend.Output(apu.buffer)
			apu.bufferIterator = 0
			apu.buffer = [BufferSize * 2]byte{}
		}
	}
}

func (apu *APU) Destroy() {
	apu.backend.Destroy()
}

func (apu *APU) ReadNR50() uint8 {
	// TODO
	return 0xf
}

func (apu *APU) WriteNR50(val uint8) {
	// volume is defined on 3 bits for a value between 0 and 7
	apu.volume.left = float64(((val >> 4) & 0x07)) / 7
	apu.volume.right = float64((val & 0x07)) / 7
}

func (apu *APU) ReadNR51() uint8 {
	return apu.nr51
}

func (apu *APU) WriteNR51(val uint8) {
	apu.nr51 = val

	n := uint8(len(apu.channels))
	for i := uint8(0); i < n; i++ {
		apu.channels[i].SetActive(bits.Test(i+4, apu.nr51), bits.Test(i, apu.nr51))
	}
}

func (apu *APU) ReadNR52() uint8 {
	val := bits.FromBool(apu.active) << 7
	for i := uint8(0); i < 4; i++ {
		val |= bits.FromBool(apu.channels[i].IsActive()) << i
	}
	return val
}

func (apu *APU) WriteNR52(val uint8) {
	// only bit 7 is writable
	apu.active = bits.Test(7, val)
}

func NewAPU(io_ *io.IO) *APU {
	var backend backend
	switch Backend {
	case "sdl2":
		backend = newSDL2(BufferSize)
	default:
		backend = &dummy{}
	}

	// TODO chan1 should probably extend chan2 as it's basically behaving
	//		the same way with other shenanigans on top
	chan1 := newChannel2()
	chan2 := newChannel2()

	apu := &APU{
		io: io_,

		backend: backend,

		channels: []channel{
			chan1,
			chan2,
		},
	}

	io_.MapRegister(io.NR50, apu.ReadNR50, apu.WriteNR50)
	io_.MapRegister(io.NR51, apu.ReadNR51, apu.WriteNR51)
	io_.MapRegister(io.NR52, apu.ReadNR52, apu.WriteNR52)

	// chan1 ports
	io_.MapRegister(io.NR13, chan1.ReadNR23, chan1.WriteNR23)
	io_.MapRegister(io.NR14, chan1.ReadNR24, chan1.WriteNR24)

	// chan2 ports
	io_.MapRegister(io.NR23, chan2.ReadNR23, chan2.WriteNR23)
	io_.MapRegister(io.NR24, chan2.ReadNR24, chan2.WriteNR24)

	io_.Write(io.NR50, 0x77)

	return apu
}
