package main

const (
	MMC1AddressUnused      = 0x4020
	MMC1AddressOptionalPRG = 0x6000
	MMC1AddressPRG0        = 0x8000
	MMC1AddressPRG1        = 0xC000
)

type mmc1 struct {
	rom *rom
}

func NewMMC1(rom *rom) mmc {
	return &mmc1{
		rom: rom,
	}
}

func (m *mmc1) Get(address uint16) byte {
	// TODO: implement here
	return 0
}

func (m *mmc1) Set(address uint16, value byte) {
	// TODO: implement here
}

func (m *mmc1) GetVRAM(address uint16) byte {
	return 0
}

func (m *mmc1) SetVRAM(address uint16, value byte) {

}
