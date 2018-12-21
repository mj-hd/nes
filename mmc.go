package main

type mmc interface {
	GetAddress(uint16) uint16
	Get(uint16) byte
	SetAddress(uint16, uint16)
	Set(uint16, byte)
}
