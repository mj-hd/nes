package main

const (
	_ = iota
	ctrlMainClean1
	ctrlMainClean2
	ctrlVRAMIncr
	ctrlSpriteTableBase
	ctrlBGCharacterTableBase
	ctrlSpriteSize
	ctrlType
	ctrlNMIVBlank
)

const (
	_ = iota
	maskType
	maskVisibleLeftBG
	maskVisibleLeftSprite
	maskVisibleBG
	maskVisibleSprite
	maskBlue
	maskGreen
	maskRed
)

const (
	_ = iota
	_
	_
	_
	statusStatus
	statusSpriteSize
	statusHitSprite
	statusScreen
)

type ppu struct {
	Ctrl     byte
	Mask     byte
	Status   byte
	OAM_addr byte
	Scroll   byte
	Addr     byte
	OAM_DMA  byte

	vram [0x2000]byte
	sram [0x0100]byte
	pram [0x0020]byte

	rom *rom
}

func (p *ppu) PowerOn() {
}

func (p *ppu) Reset() {
}

func (p *ppu) Tick() {

}

func (p *ppu) Get() byte {
	return p.vram[p.Addr]
}

func (p *ppu) Set(value byte) {
	p.vram[p.Addr] = value
	if p.Ctrl&ctrlVRAMIncr == 1 {
		p.Addr += 32
	} else {
		p.Addr += 1
	}
}

func (p *ppu) GetOAM() byte {
	var base uint16
	if p.Ctrl&ctrlSpriteTableBase == 1 {
		base = 0x1000
	} else {
		base = 0x0000
	}
	return p.vram[uint16(p.OAM_addr)+base]
}

func (p *ppu) SetOAM(v byte) {
	var base uint16
	if p.Ctrl&ctrlSpriteTableBase == 1 {
		base = 0x1000
	} else {
		base = 0x0000
	}
	p.vram[uint16(p.OAM_addr)+base] = v
}
