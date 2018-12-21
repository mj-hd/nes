package main

import (
	"time"
)

type NES struct {
	CPU *cpu
	PPU *ppu
	MM  *mm
	MMC mmc
	ROM *rom
}

func NewNES() *NES {
	mm := &mm{}
	rom := &rom{}
	ppu := &ppu{
		rom: rom,
	}
	apu := &apu{}
	mmc := &mmc1{
		ppu: ppu,
		apu: apu,
		mm:  mm,
		rom: rom,
	}
	return &NES{
		CPU: &cpu{
			MMC: mmc,
		},
		PPU: ppu,
		MM:  mm,
		ROM: rom,
	}
}

func (n *NES) Run() {
	for {
		n.PPU.Tick()
		n.CPU.Tick()
		time.Sleep(12 * time.Millisecond)
	}
}

func (n *NES) Load(file string) error {
	return n.ROM.Load(file)
}

func (n *NES) PowerOn() {
	n.CPU.PowerOn()
	n.PPU.PowerOn()
}

func (n *NES) Reset() {
	n.CPU.Reset()
	n.PPU.Reset()
}
