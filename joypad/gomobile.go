// +build android

package joypad

type gomobile struct {
}

func (g *gomobile) ReadDirections() uint8 {
	return 0
}

func (g *gomobile) ReadButtons() uint8 {
	return 0
}

func NewBackend() *gomobile {
	return &gomobile{}
}
