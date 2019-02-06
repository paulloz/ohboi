package cpu

type Register struct {
	hilo uint16
}

func (r *Register) HiLo() uint16 {
	return r.hilo
}

func (r *Register) Hi() uint8 {
	return uint8(r.hilo >> 8)
}

func (r *Register) Lo() uint8 {
	return uint8(r.hilo & 0xFF)
}

func (r *Register) Set(value uint16) {
	r.hilo = value
}

func (r *Register) SetHi(value uint8) {
	r.hilo = (r.hilo & 0xFF) | (uint16(value) << 8)
}

func (r *Register) SetLo(value uint8) {
	r.hilo = (r.hilo & 0xFF00) | uint16(value)
}

func NewRegister(value uint16) Register {
	return Register{hilo: value}
}
