package lcd

type Color struct {
	R uint8
	G uint8
	B uint8
}

var Greys = [4]Color{
	{R: 255, G: 255, B: 255},
	{R: 170, G: 170, B: 170},
	{R: 85, G: 85, B: 85},
	{R: 0, G: 0, B: 0},
}

func NewColor(r uint8, g uint8, b uint8) Color {
	return Color{R: r, G: g, B: b}
}
