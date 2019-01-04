package main

const (
	MMC0AddressUnused     = 0x4020
	MMC0AddressBatteryRAM = 0x6000
	MMC0AddressFirstPRG   = 0x8000
	MMC0AddressSecondPRG  = 0xC000

	MMC0BankSize = 0x4000
)

type mmc0 struct {
	rom *rom

	bankAddr1 uint16
	bankAddr2 uint16
}

func NewMMC0(rom *rom) mmc {
	return &mmc0{
		rom:       rom,
		bankAddr1: 0,
		bankAddr2: uint16(len(rom.PRG) - MMC0BankSize),
	}
}

func (m *mmc0) Get(address uint16) byte {
	switch {
	case address < MMC0AddressBatteryRAM:
		return 0
	case address < MMC0AddressFirstPRG:
		return 0
	case address < MMC0AddressSecondPRG:
		return m.rom.Get(address - MMC0AddressFirstPRG + m.bankAddr1)
	}
	return m.rom.Get(address - MMC0AddressSecondPRG + m.bankAddr2)
}

func (m *mmc0) Set(address uint16, value byte) {
	m.rom.Set(address, value)
}
