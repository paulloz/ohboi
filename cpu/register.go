package cpu

// Register ...
type Register struct {
	hilo uint16
}

// HiLo ...
func (r *Register) HiLo() uint16 {
	return r.hilo
}

// Hi ...
func (r *Register) Hi() uint8 {
	return uint8(r.hilo >> 8)
}

// Lo ...
func (r *Register) Lo() uint8 {
	return uint8(r.hilo & 0xFF)
}

// Set ...
func (r *Register) Set(value uint16) {
	r.hilo = value
}

// SetHi ...
func (r *Register) SetHi(value uint8) {
	r.hilo = (r.hilo & 0xFF) | (uint16(value) << 8)
}

// SetLo ...
func (r *Register) SetLo(value uint8) {
	r.hilo = (r.hilo & 0xFF00) | uint16(value)
}

func NewRegister(value uint16) Register {
	return Register{hilo: value}
}
