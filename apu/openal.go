// +build android

package apu

import (
	"github.com/paulloz/ohboi/consts"
	"golang.org/x/mobile/exp/audio/al"
)

type openAL struct {
	buffers []al.Buffer
	source  al.Source
}

const bufferCount = 128

var queue []al.Buffer
var i int

func (o *openAL) Output(buffer [BufferSize * 2]byte) {
	for processed := o.source.BuffersProcessed(); processed > 0; processed-- {
		o.source.UnqueueBuffers(queue[0])
		queue = queue[1:]
	}

	o.buffers[i%bufferCount].BufferData(al.FormatStereo8, buffer[4:], consts.APUSampleRate)
	o.source.QueueBuffers(o.buffers[i%bufferCount])

	queue = append(queue, o.buffers[i%bufferCount])

	if o.source.BuffersQueued() == 1 && i == 0 {
		al.PlaySources(o.source)
	}

	i = (i + 1) % bufferCount

	al.PlaySources(o.source)
}

func (o *openAL) Destroy() {
	al.DeleteBuffers(o.buffers...)
	al.CloseDevice()
}

func newBackend() (*openAL, error) {
	if err := al.OpenDevice(); err != nil {
		return nil, err
	}

	buffers := al.GenBuffers(bufferCount)
	sources := al.GenSources(1)

	o := &openAL{buffers: buffers, source: sources[0]}

	return o, nil
}
