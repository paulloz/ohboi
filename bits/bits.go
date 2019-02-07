package bits

func Test(b uint8, value uint8) bool {
	return ((value >> b) & 1) == 1
}

func Set(b uint8, value uint8) uint8 {
	return value | (1 << b)
}

func Reset(b uint8, value uint8) uint8 {
	return value &^ (1 << b)
}

func HalfCarryCheck(a uint8, b uint8) bool {
	return (((a & 0xf) + (b & 0xf)) & 0x10) == 0x10
}

func FromBool(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
