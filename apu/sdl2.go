package apu

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/paulloz/ohboi/consts"
)

type sdl2 struct {
	dev sdl.AudioDeviceID
}

func (s *sdl2) Output(buffer [BufferSize * 2]byte) {
	sdl.Do(func() {
		// fmt.Println(buffer)
		fmt.Printf("")
		sdl.QueueAudio(s.dev, buffer[:])
	})
}

func (s *sdl2) Destroy() {
	sdl.Do(func() {
		sdl.CloseAudioDevice(s.dev)
		sdl.PauseAudioDevice(s.dev, true)
		sdl.PauseAudio(true)
	})
}

func newSDL2(samples uint16) *sdl2 {
	var dev sdl.AudioDeviceID

	// sdl.Do(func() {
	// })
	var want, have *sdl.AudioSpec
	var err error

	want = &sdl.AudioSpec{
		Freq:     consts.APUSampleRate,
		Format:   sdl.AUDIO_U8,
		Channels: 2,
		Samples:  BufferSize,
	}

	dev, err = sdl.OpenAudioDevice("", false, want, have, 0)
	if err != nil {
		panic(err)
	}

	sdl.PauseAudio(false)
	sdl.PauseAudioDevice(dev, false)

	return &sdl2{
		dev: dev,
	}
}
