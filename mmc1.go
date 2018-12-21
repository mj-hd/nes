package main

const (
	MMC1AddressRAM             = 0x0000
	MMC1AddressReserved        = 0x07FF
	MMC1AddressPPUCtrl         = 0x2000
	MMC1AddressPPUMask         = 0x2001
	MMC1AddressPPUStatus       = 0x2002
	MMC1AddressOAMAddr         = 0x2003
	MMC1AddressOAMData         = 0x2004
	MMC1AddressPPUScroll       = 0x2005
	MMC1AddressPPUAddr         = 0x2006
	MMC1AddressPPUVRAM         = 0x2007
	MMC1AddressAPUPulse1       = 0x4000
	MMC1AddressAPUPulse2       = 0x4004
	MMC1AddressAPUTriangle     = 0x4008
	MMC1AddressAPUNoise        = 0x400C
	MMC1AddressAPUDMC          = 0x4010
	MMC1AddressOAMDMA          = 0x4014
	MMC1AddressAPUStatus       = 0x4015
	MMC1AddressJoy1            = 0x4016
	MMC1AddressAPUFrameCounter = 0x4017
	MMC1AddressJoy2            = 0x4017
	MMC1AddressRAMExtend       = 0x4020
	MMC1AddressBackup          = 0x7FFF
	MMC1AddressROM             = 0x8000
)

type mmc1 struct {
	mm  *mm
	ppu *ppu
	apu *apu
	rom *rom
}

func (m *mmc1) GetAddress(address uint16) uint16 {
	return uint16(uint16(m.Get(address))<<8 + uint16(m.Get(address+1)))
}

func (m *mmc1) Get(address uint16) byte {
	if address < MMC1AddressReserved {
		return m.mm.ram[address]
	} else if address < MMC1AddressRAMExtend {
		switch address {
		case MMC1AddressPPUCtrl:
			return m.ppu.Ctrl
		case MMC1AddressPPUMask:
			return m.ppu.Mask
		case MMC1AddressPPUStatus:
			return m.ppu.Status
		case MMC1AddressOAMAddr:
			return m.ppu.OAM_addr
		case MMC1AddressOAMData:
			return m.ppu.GetOAM()
		case MMC1AddressPPUScroll:
			return m.ppu.Scroll
		case MMC1AddressPPUAddr:
			return m.ppu.Addr
		case MMC1AddressPPUVRAM:
			return m.ppu.Get()
		case MMC1AddressOAMDMA:
			return m.ppu.OAM_DMA
		case MMC1AddressAPUStatus:
			return m.apu.Status
		case MMC1AddressJoy1:
			return byte(0)
		case MMC1AddressAPUFrameCounter:
			return m.apu.FrameCounter
		}
		return byte(0)
	} else if address < MMC1AddressBackup {
		return m.mm.ram_extend[address-MMC1AddressRAMExtend]
	} else if address < MMC1AddressROM {
		return m.mm.backup[address-MMC1AddressBackup]
	} else if address <= 0xFFFF {
		return m.rom.Get(address - MMC1AddressROM)
	} else {
		panic("RAM Range Exception")
	}
}

func (m *mmc1) SetAddress(address uint16, value uint16) {
	m.Set(address, uint8(value&0x0F))
	m.Set(address+1, uint8(value>>8))
}

func (m *mmc1) Set(address uint16, value byte) {
	if address < MMC1AddressReserved {
		m.mm.ram[address] = value
	} else if address < MMC1AddressRAMExtend {
		switch address {
		case MMC1AddressPPUCtrl:
			m.ppu.Ctrl = value
		case MMC1AddressPPUMask:
			m.ppu.Mask = value
		case MMC1AddressPPUStatus:
			m.ppu.Status = value
		case MMC1AddressOAMAddr:
			m.ppu.OAM_addr = value
		case MMC1AddressOAMData:
			m.ppu.SetOAM(value)
		case MMC1AddressPPUScroll:
			m.ppu.Scroll = value
		case MMC1AddressPPUAddr:
			m.ppu.Addr = value
		case MMC1AddressPPUVRAM:
			m.ppu.Set(value)
		case MMC1AddressOAMDMA:
			m.ppu.OAM_DMA = value
		case MMC1AddressAPUStatus:
			m.apu.Status = value
		case MMC1AddressJoy1:
			// ???
		case MMC1AddressAPUFrameCounter:
			m.apu.FrameCounter = value
		}
		// ???
	} else if address < MMC1AddressBackup {
		m.mm.ram_extend[address-MMC1AddressRAMExtend] = value
	} else if address < MMC1AddressROM {
		m.mm.backup[address-MMC1AddressBackup] = value
	} else if address <= 0xFFFF {
		m.rom.Set(address-MMC1AddressROM, value)
	} else {
		panic("RAM Range Exception")
	}
}
