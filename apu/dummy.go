package apu

type dummy struct{}

func (d *dummy) Output([BufferSize * 2]byte) {}
func (d *dummy) Destroy()                    {}
