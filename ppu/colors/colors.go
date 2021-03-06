package colors

type Color struct {
	R uint8
	G uint8
	B uint8
}

type Palette [4]Color

var Greens = Palette{
	{R: 224, G: 248, B: 208},
	{R: 136, G: 192, B: 112},
	{R: 52, G: 104, B: 86},
	{R: 8, G: 24, B: 32},
}

var Greys = Palette{
	{R: 255, G: 255, B: 255},
	{R: 170, G: 170, B: 170},
	{R: 85, G: 85, B: 85},
	{R: 0, G: 0, B: 0},
}

var SuperGameboy = Palette{
	{R: 0xf7, G: 0xe7, B: 0xc6},
	{R: 0xd6, G: 0x8e, B: 0x49},
	{R: 0xa6, G: 0x37, B: 0x25},
	{R: 0x33, G: 0x1e, B: 0x50},
}

func NewColor(r uint8, g uint8, b uint8) Color {
	return Color{R: r, G: g, B: b}
}
