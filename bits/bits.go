package bits

func Test(b uint8, value uint8) bool {
	return (value >> b) == 1
}

func Set(b uint8, value uint8) uint8 {
	return value | (1 << b)
}

func Reset(b uint8, value uint8) uint8 {
	return value & (0xF ^ (1 << b))
}
