package main

type NES struct {
	CPU *cpu
	PPU *ppu
	APU *apu
	MMC mmc
	ROM *rom
}

func NewNES(file string, renderer Renderer) (*NES, error) {
	n := &NES{}
	r := &rom{}
	err := r.Load(file)
	if err != nil {
		return nil, err
	}
	n.ROM = r
	n.APU = &apu{}
	n.MMC = NewMMC(r.Header.MapperNum, r)
	n.PPU = NewPPU(n.MMC, renderer)
	n.CPU = &cpu{
		MMC: n.MMC,
		PPU: n.PPU,
		APU: n.APU,
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
