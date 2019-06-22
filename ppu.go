package main

import (
	"image/color"
	"log"
)

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
	_
	statusSpriteOverflow
	statusSpriteZeroHit
	statusVBlank
)

const (
	PPUWidth         = 340
	PPUHeight        = 261
	PPUVisibleWidth  = 256
	PPUVisibleHeight = 240
)

var colors = []color.RGBA{
	{0x80, 0x80, 0x80, 0xFF}, {0x00, 0x3D, 0xA6, 0xFF}, {0x00, 0x12, 0xB0, 0xFF}, {0x44, 0x00, 0x96, 0xFF},
	{0xA1, 0x00, 0x5E, 0xFF}, {0xC7, 0x00, 0x28, 0xFF}, {0xBA, 0x06, 0x00, 0xFF}, {0x8C, 0x17, 0x00, 0xFF},
	{0x5C, 0x2F, 0x00, 0xFF}, {0x10, 0x45, 0x00, 0xFF}, {0x05, 0x4A, 0x00, 0xFF}, {0x00, 0x47, 0x2E, 0xFF},
	{0x00, 0x41, 0x66, 0xFF}, {0x00, 0x00, 0x00, 0xFF}, {0x05, 0x05, 0x05, 0xFF}, {0x05, 0x05, 0x05, 0xFF},
	{0xC7, 0xC7, 0xC7, 0xFF}, {0x00, 0x77, 0xFF, 0xFF}, {0x21, 0x55, 0xFF, 0xFF}, {0x82, 0x37, 0xFA, 0xFF},
	{0xEB, 0x2F, 0xB5, 0xFF}, {0xFF, 0x29, 0x50, 0xFF}, {0xFF, 0x22, 0x00, 0xFF}, {0xD6, 0x32, 0x00, 0xFF},
	{0xC4, 0x62, 0x00, 0xFF}, {0x35, 0x80, 0x00, 0xFF}, {0x05, 0x8F, 0x00, 0xFF}, {0x00, 0x8A, 0x55, 0xFF},
	{0x00, 0x99, 0xCC, 0xFF}, {0x21, 0x21, 0x21, 0xFF}, {0x09, 0x09, 0x09, 0xFF}, {0x09, 0x09, 0x09, 0xFF},
	{0xFF, 0xFF, 0xFF, 0xFF}, {0x0F, 0xD7, 0xFF, 0xFF}, {0x69, 0xA2, 0xFF, 0xFF}, {0xD4, 0x80, 0xFF, 0xFF},
	{0xFF, 0x45, 0xF3, 0xFF}, {0xFF, 0x61, 0x8B, 0xFF}, {0xFF, 0x88, 0x33, 0xFF}, {0xFF, 0x9C, 0x12, 0xFF},
	{0xFA, 0xBC, 0x20, 0xFF}, {0x9F, 0xE3, 0x0E, 0xFF}, {0x2B, 0xF0, 0x35, 0xFF}, {0x0C, 0xF0, 0xA4, 0xFF},
	{0x05, 0xFB, 0xFF, 0xFF}, {0x5E, 0x5E, 0x5E, 0xFF}, {0x0D, 0x0D, 0x0D, 0xFF}, {0x0D, 0x0D, 0x0D, 0xFF},
	{0xFF, 0xFF, 0xFF, 0xFF}, {0xA6, 0xFC, 0xFF, 0xFF}, {0xB3, 0xEC, 0xFF, 0xFF}, {0xDA, 0xAB, 0xEB, 0xFF},
	{0xFF, 0xA8, 0xF9, 0xFF}, {0xFF, 0xAB, 0xB3, 0xFF}, {0xFF, 0xD2, 0xB0, 0xFF}, {0xFF, 0xEF, 0xA6, 0xFF},
	{0xFF, 0xF7, 0x9C, 0xFF}, {0xD7, 0xE8, 0x95, 0xFF}, {0xA6, 0xED, 0xAF, 0xFF}, {0xA2, 0xF2, 0xDA, 0xFF},
	{0x99, 0xFF, 0xFC, 0xFF}, {0xDD, 0xDD, 0xDD, 0xFF}, {0x11, 0x11, 0x11, 0xFF}, {0x11, 0x11, 0x11, 0xFF},
}

type ppu struct {
	Ctrl          byte
	Mask          byte
	buffer        []byte
	vblank        bool
	spriteZeroHit bool

	Cycle int
	Line  int

	oam [0x100]byte
	bus bus
	dma *dma

	renderer Renderer
}

func NewPPU(bus bus, dma *dma, renderer Renderer) *ppu {
	return &ppu{
		buffer:   make([]byte, 0, 2),
		bus:      bus,
		dma:      dma,
		renderer: renderer,
	}
}

func (p *ppu) PowerOn() {
}

func (p *ppu) Reset() {
}

func (p *ppu) Tick() {
	p.Cycle++
	if p.Cycle > PPUWidth {
		p.Cycle = 0
		p.Line++
		if p.Line > PPUHeight {
			p.vblank = false
			p.Line = 0
			p.renderer.Render()
		} else if p.Line > PPUVisibleHeight {
			p.vblank = true
			// TODO: fire nmi interruption
		}
	}
	if p.Cycle%8 == 0 {
		if p.Line%8 == 0 {
			p.drawBG(p.Cycle, p.Line)
		}
	}
	if p.Cycle < 256-8 && p.Line < 240-8 {
		p.drawSprites(p.Cycle, p.Line)
	}
	//log.Printf("PPU Cycle:%d Line:%d\n", p.Cycle, p.Line)
}

func (p *ppu) isVisible() bool {
	return p.Cycle < PPUVisibleWidth && p.Line < PPUVisibleHeight
}

func (p *ppu) drawBG(cycle, line int) {
	if !p.isVisible() {
		return
	}
	x := cycle
	y := line
	tileX, tileY := x/8, y/8
	attrX, attrY := x/16, y/16
	tileCharNum := p.getTileCharNumber(tileX, tileY)
	paletteNum := p.getPaletteNumber(attrX, attrY, tileX, tileY)
	colors := p.getColors(paletteNum, PPUAddressPaletteBG)
	char := p.getChar(int(tileCharNum))
	pixels := p.toPixels(char, colors)
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			p.renderer.SetPixel(tileX*8+j, tileY*8+i, pixels[i][j])
		}
	}
}

func (p *ppu) drawSprites(cycle, line int) {
	if !p.isVisible() {
		return
	}
	x := cycle
	y := line
	for i := 0; i < len(p.oam); i += 4 {
		spriteY := int(p.oam[i])
		spriteX := int(p.oam[i+3])
		if x == 0 || y == 0 || x != spriteX || y != spriteY {
			continue
		}
		log.Printf("draw sprite: (%d, %d)", x, y)
		spriteTile := p.oam[i+1]
		spriteAttr := p.oam[i+2] & 0x3
		p.drawSprite(x, y, spriteTile, spriteAttr)
	}
}

func (p *ppu) drawSprite(x, y int, tile, attr byte) {
	colors := p.getColors(int(attr), PPUAddressPaletteSprite)
	char := p.getChar(int(tile))
	pixels := p.toPixels(char, colors)
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			p.renderer.SetPixel(x+j, y+i, pixels[i][j])
		}
	}
}

func (p *ppu) toPixels(char []byte, colors []color.RGBA) [8][8]color.RGBA {
	var res [8][8]color.RGBA
	for y := 0; y < 8; y++ {
		charLow := char[y]
		charHi := char[y+8]
		for x := 0; x < 8; x++ {
			colorNum := (((charHi >> uint(7-x)) & 1) << 1) | ((charLow >> uint(7-x)) & 1)
			res[y][x] = colors[colorNum]
		}
	}
	return res
}

func (p *ppu) getColors(paletteNumber int, baseAddr uint16) []color.RGBA {
	res := make([]color.RGBA, 0, 4)
	for i := 0; i < 4; i++ {
		colorNum := p.bus.Get(baseAddr + uint16(paletteNumber*4+i))
		res = append(res, colors[int(colorNum)])
	}
	return res
}

func (p *ppu) getTileCharNumber(tileX, tileY int) byte {
	tileAddr := uint16((tileY*32 + tileX) + PPUAddressNameTable0)
	return p.bus.Get(tileAddr)
}

func (p *ppu) getPaletteNumber(attrX, attrY, tileX, tileY int) int {
	attrAddr := uint16((attrY*16 + attrX) + PPUAddressAttrTable0)
	attr := p.bus.Get(attrAddr)
	index := 0
	if tileX%2 == 1 {
		index |= 1
	}
	if tileY%2 == 1 {
		index |= 1 << 1
	}
	attr >>= uint(index * 2)
	attr &= 0x03
	return int(attr)
}

func (p *ppu) getChar(charNum int) []byte {
	charAddr := uint16(charNum * 16)
	char := make([]byte, 0, 16)
	for i := uint16(0); i < 16; i++ {
		char = append(char, p.bus.Get(charAddr+i))
	}
	return char
}

func (p *ppu) getAddr() uint16 {
	addr := uint16(p.buffer[1])
	addr |= uint16(p.buffer[0]) << 8
	return addr
}

func (p *ppu) setAddr(addr uint16) {
	p.buffer[0] = byte(addr >> 8)
	p.buffer[1] = byte(addr & 0xFF)
}

func (p *ppu) GetData() byte {
	return p.bus.Get(p.getAddr())
}

func (p *ppu) SetData(value byte) {
	log.Printf("set data(%x)", value)
	addr := p.getAddr()
	p.bus.Set(addr, value)
	if p.Ctrl&ctrlVRAMIncr == 1 {
		addr += 32
	} else {
		addr += 1
	}
	p.setAddr(addr)
}

func (p *ppu) GetStatus() byte {
	val := byte(0)
	if p.vblank {
		val |= 1 << statusVBlank
	}
	if p.spriteZeroHit {
		val |= 1 << statusSpriteZeroHit
	}
	p.buffer = p.buffer[:]
	p.vblank = false
	return val
}

func (p *ppu) GetOAM() byte {
	return 0
}

func (p *ppu) SetOAM(v byte) {
}

func (p *ppu) SetScroll(v byte) {
	if len(p.buffer) == 2 {
		p.buffer = p.buffer[:0]
	}
	p.buffer = append(p.buffer, v)
}

func (p *ppu) SetAddr(v byte) {
	log.Printf("set addr(%x), %v", v, p.buffer)
	if len(p.buffer) >= 2 {
		p.buffer = p.buffer[:0]
	}
	p.buffer = append(p.buffer, v)
}

func (p *ppu) SetDMA(v byte) {
	cpuAddr := uint16(v) << 8
	log.Printf("dma(%x): %x", v, cpuAddr)
	p.dma.Transfer(cpuAddr, p.oam[:])
	log.Printf("copied %v", p.oam)
}
