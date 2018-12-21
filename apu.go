package main

type apu struct {
	Pulse [0x4008 - 0x4000]byte
	Triangle [0x400C - 0x4008]byte
	Noise [0x4010 - 0x400C]byte
	DMC [0x4014 - 0x4010]byte
	Status byte
	FrameCounter byte
}
