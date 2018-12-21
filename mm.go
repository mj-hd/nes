package main

type mm struct {
	ram [0x07FF - 0x0000]byte
	ram_extend [0x5FFF - 0x4020]byte
	backup [0x7FFF - 0x6000]byte
	rom_low [0xBFFF - 0x8000]byte
	rom_high [0xFFFF - 0xC000]byte
}

func (r *mm) Clear() {
	// clear all memory
}
