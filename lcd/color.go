package lcd

type Color struct {
	R uint8
	G uint8
	B uint8
}

var Greys = [4]Color{
	{R: 52, G: 104, B: 86},
	{R: 136, G: 192, B: 112},
	{R: 224, G: 248, B: 208},
	{R: 8, G: 24, B: 32},
}

func NewColor(r uint8, g uint8, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

func (lcd *LCD) getPalette(ioAddr uint8) [4]Color {
	palette := lcd.io.Read(ioAddr)
	colorPalette := [4]Color{}
	for i := uint8(0); i < 8; i += 2 {
		shade := (palette >> i) & 3
		colorPalette[i/2] = Greys[shade]
	}
	return colorPalette
}
