package main

const (
	MMC0AddressUnused     = 0x4020
	MMC0AddressBatteryRAM = 0x6000
	MMC0AddressPRG0       = 0x8000
	MMC0AddressPRG1       = 0xC000

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
	case address < MMC_VRAM_Limit:
		return m.rom.GetCHR(address)
	case address < MMC0AddressBatteryRAM:
		return 0
	case address < MMC0AddressPRG0:
		return 0
	case address < MMC0AddressPRG1:
		return m.rom.GetPRG(address - MMC0AddressPRG0 + m.bankAddr1)
	}
	return m.rom.GetPRG(address - MMC0AddressPRG1 + m.bankAddr2)
}

func (m *mmc0) Set(address uint16, value byte) {
	switch {
	case address < MMC_VRAM_Limit:
		m.rom.SetCHR(address, value)
	case address < MMC0AddressBatteryRAM:
	case address < MMC0AddressPRG0:
	case address < MMC0AddressPRG1:
		m.rom.SetPRG(address-MMC0AddressPRG0+m.bankAddr1, value)
	default:
		m.rom.SetPRG(address-MMC0AddressPRG1+m.bankAddr2, value)
	}
}
