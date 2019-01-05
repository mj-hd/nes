package main

const MMC_VRAM_Limit = 0x2000

type mmc interface {
	Get(uint16) byte
	Set(uint16, byte)
}

func NewMMC(mapper_num int, rom *rom) mmc {
	switch mapper_num {
	case 1:
		return NewMMC1(rom)
	}
	return NewMMC0(rom)
}
