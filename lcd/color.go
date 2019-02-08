package lcd

type color struct {
	r uint8
	g uint8
	b uint8
}

func newColor(r uint8, g uint8, b uint8) color {
	return color{r: r, g: g, b: b}
}
