package main

type bus interface {
	Get(uint16) byte
	Set(uint16, byte)
}

const (
	PPUAddressPattern0         = 0x0000
	PPUAddressPattern1         = 0x1000
	PPUAddressNameTable0       = 0x2000
	PPUAddressAttrTable0       = 0x23C0
	PPUAddressNameTable1       = 0x2400
	PPUAddressAttrTable1       = 0x27C0
	PPUAddressNameTable2       = 0x2800
	PPUAddressAttrTable2       = 0x2BC0
	PPUAddressNameTable3       = 0x2C00
	PPUAddressAttrTable3       = 0x2FC0
	PPUAddressNameTableMirror  = 0x3000
	PPUAddressPaletteBG        = 0x3F00
	PPUAddressPaletteSprite    = 0x3F10
	PPUAddressPaletteRAMMirror = 0x3F20

	PPUAddressVRAM       = 0x2000
	PPUAddressVRAM_Limit = 0x4000
)

type ppuBus struct {
	vram []byte
	mmc  mmc
}

func NewPPUBus(vram []byte, mmc mmc) bus {
	return &ppuBus{
		vram: vram,
		mmc:  mmc,
	}
}

func (b *ppuBus) Get(addr uint16) byte {
	switch {
	case addr < PPUAddressVRAM:
		return b.mmc.Get(addr)
	case PPUAddressVRAM <= addr && addr < PPUAddressPaletteBG:
		diff := uint16((addr - PPUAddressVRAM) % 0xF00)
		return b.vram[diff]
	case PPUAddressPaletteBG <= addr:
		diff := uint16((addr - PPUAddressPaletteBG) % 0x20)
		base := uint16(PPUAddressPaletteBG - PPUAddressVRAM)
		return b.vram[base+diff]
	case PPUAddressVRAM_Limit <= addr:
		return b.mmc.Get(addr)
	default:
		return b.vram[addr-PPUAddressVRAM]
	}
}

func (b *ppuBus) Set(addr uint16, val byte) {
	switch {
	case addr < PPUAddressVRAM:
		b.mmc.Set(addr, val)
	case PPUAddressVRAM <= addr && addr < PPUAddressPaletteBG:
		diff := uint16((addr - PPUAddressVRAM) % 0xF00)
		b.vram[diff] = val
	case PPUAddressPaletteBG <= addr:
		diff := uint16((addr - PPUAddressPaletteBG) % 0x20)
		base := uint16(PPUAddressPaletteBG - PPUAddressVRAM)
		b.vram[base+diff] = val
	case PPUAddressVRAM_Limit <= addr:
		b.mmc.Set(addr, val)
	default:
		b.vram[addr-PPUAddressVRAM] = val
	}
}

const (
	AddressRAM             = 0x0000
	AddressMirror1         = 0x0800
	AddressMirror2         = 0x1000
	AddressMirror3         = 0x1800
	AddressPPUCtrl         = 0x2000
	AddressPPUMask         = 0x2001
	AddressPPUStatus       = 0x2002
	AddressOAMAddr         = 0x2003
	AddressOAMData         = 0x2004
	AddressPPUScroll       = 0x2005
	AddressPPUAddr         = 0x2006
	AddressPPUData         = 0x2007
	AddressAPUPulse1       = 0x4000
	AddressAPUPulse2       = 0x4004
	AddressAPUTriangle     = 0x4008
	AddressAPUNoise        = 0x400C
	AddressAPUDMC          = 0x4010
	AddressOAMDMA          = 0x4014
	AddressAPUStatus       = 0x4015
	AddressJoy1            = 0x4016
	AddressAPUFrameCounter = 0x4017
	AddressJoy2            = 0x4017
	AddressAPUTest         = 0x4018
)

type cpuBus struct {
	ppu  *ppu
	apu  *apu
	mmc  mmc
	wram []byte
}

func NewCPUBus(wram []byte, ppu *ppu, apu *apu, mmc mmc) bus {
	return &cpuBus{
		wram: wram,
		ppu:  ppu,
		apu:  apu,
		mmc:  mmc,
	}
}

func (b *cpuBus) Get(address uint16) byte {
	switch {
	case address < AddressMirror1:
		return b.wram[address]
	case address < AddressMirror2:
		return b.wram[address-AddressMirror1]
	case address < AddressMirror3:
		return b.wram[address-AddressMirror2]
	case address == AddressPPUCtrl:
		return b.ppu.Ctrl
	case address == AddressPPUMask:
		return b.ppu.Mask
	case address == AddressPPUStatus:
		return b.ppu.GetStatus()
	case address == AddressOAMAddr:
		return 0
	case address == AddressOAMData:
		return b.ppu.GetOAM()
	case address == AddressPPUScroll:
		return 0
	case address == AddressPPUAddr:
		return 0
	case address == AddressPPUData:
		return b.ppu.GetData()
	case address == AddressOAMDMA:
		return 0
	case address == AddressAPUStatus:
		return b.apu.Status
	case address == AddressJoy1:
		return 0
	case address == AddressAPUFrameCounter:
		return b.apu.FrameCounter
	}
	return b.mmc.Get(address)
}

func (b *cpuBus) Set(address uint16, value byte) {
	switch {
	case address < AddressMirror1:
		b.wram[address] = value
	case address < AddressMirror2:
		b.wram[address-AddressMirror1] = value
	case address < AddressMirror3:
		b.wram[address-AddressMirror2] = value
	case address == AddressPPUCtrl:
		b.ppu.Ctrl = value
	case address == AddressPPUMask:
		b.ppu.Mask = value
	case address == AddressPPUStatus:
	case address == AddressOAMAddr:
	case address == AddressOAMData:
		b.ppu.SetOAM(value)
	case address == AddressPPUScroll:
		b.ppu.SetScroll(value)
	case address == AddressPPUAddr:
		b.ppu.SetAddr(value)
	case address == AddressPPUData:
		b.ppu.SetData(value)
	case address == AddressOAMDMA:
		b.ppu.SetDMA(value)
	case address == AddressAPUStatus:
		b.apu.Status = value
	case address == AddressJoy1:
		// nop
	case address == AddressAPUFrameCounter:
		b.apu.FrameCounter = value
	default:
		b.mmc.Set(address, value)
	}
}

