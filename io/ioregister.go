package io

type IORegister struct {
	value uint8
}

func newIORegister(value uint8) *IORegister {
	return &IORegister{value: value}
}
