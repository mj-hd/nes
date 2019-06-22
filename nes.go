package main

type NES struct {
	CPU    *cpu
	PPU    *ppu
	APU    *apu
	MMC    mmc
	ROM    *rom
	PPUBus bus
	CPUBus bus

	vram [0x2000]byte
	wram [0x0800]byte
}

func NewNES(file string, renderer Renderer) (*NES, error) {
	n := &NES{}
	r := &rom{}
	err := r.Load(file)
	if err != nil {
		return nil, err
	}
	n.ROM = r
	apu := &apu{}
	n.APU = apu
	mmc := NewMMC(r.Header.MapperNum, r)
	n.MMC = mmc
	ppuBus := NewPPUBus(n.vram[:], mmc)
	dma := &dma{}
	ppu := NewPPU(ppuBus, dma, renderer)
	n.PPU = ppu
	cpuBus := NewCPUBus(n.wram[:], ppu, apu, mmc)
	dma.bus = cpuBus
	n.CPU = &cpu{
		MMC: n.MMC,
		PPU: n.PPU,
		APU: n.APU,
		bus: cpuBus,
	}
	return n, nil
}

func (n *NES) Tick() {
	n.CPU.Tick()
	n.PPU.Tick()
	n.PPU.Tick()
	n.PPU.Tick()
	//time.Sleep(1 * time.Millisecond)
}

func (n *NES) PowerOn() {
	n.CPU.PowerOn()
	n.PPU.PowerOn()
}

func (n *NES) Reset() {
	n.CPU.Reset()
	n.PPU.Reset()
}
