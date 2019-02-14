package apu

import (
	"math"

	"github.com/paulloz/ohboi/consts"
)

type channel interface {
	Sample() uint16

	IsActive() bool
}

type basechannel struct {
	frequency float64
	duty      float64

	time float64

	active bool
}

func (c *basechannel) Sample() (sample uint16) {
	c.time += (c.frequency * (math.Pi * 2)) / consts.APUSampleRate

	if c.active {
		/*
		 * 12.5% ( _-------_-------_------- ) -> duty = 0.75
		 * 25%   ( __------__------__------ ) -> duty = 0.5
		 * 50%   ( ____----____----____---- ) -> duty = 0
		 * 75%   ( ______--______--______-- ) -> duty = -0.5
		 * We make a sin wave and send back 0xff or 0x00 depending on where we're on the curve.
		 * // TODO This needs to move where the duty definition is, as chan3 and 4 do not sample like this at all
		 */
		if math.Sin(c.time) >= c.duty {
			sample = 0xff
		}
	}

	return sample
}

func (c *basechannel) IsActive() bool {
	return c.active
}
