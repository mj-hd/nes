package main

import (
	"time"
)

type NES struct {
	CPU *cpu
	PPU *ppu
	APU *apu
	MMC mmc
	ROM *rom
}

func NewNES(file string) (*NES, error) {
	n := &NES{}
	r := &rom{}
	p := &ppu{
		rom: r,
	}
	a := &apu{}
	err := r.Load(file)
	if err != nil {
		return nil, err
	}
	n.ROM = r
	n.PPU = p
	n.APU = a
	n.MMC = NewMMC(r.Header.MapperNum, r)
	n.CPU = &cpu{
		MMC: n.MMC,
		PPU: n.PPU,
		APU: n.APU,
	}
	return n, nil
}

func (n *NES) Run() {
	for {
		n.PPU.Tick()
		n.CPU.Tick()
		time.Sleep(12 * time.Millisecond)
	}
}

func (n *NES) PowerOn() {
	n.CPU.PowerOn()
	n.PPU.PowerOn()
}

func (n *NES) Reset() {
	n.CPU.Reset()
	n.PPU.Reset()
}
