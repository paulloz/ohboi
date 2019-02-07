package io

type IORegister struct {
	value uint8
}

func NewIORegister(value uint8) *IORegister {
	return &IORegister{value: value}
}
