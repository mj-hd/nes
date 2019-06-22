package main

type dma struct {
	bus bus
}

func (d *dma) Transfer(addr uint16, target []byte) {
	for i := uint16(0); int(i) < len(target); i++ {
		target[i] = d.bus.Get(addr + i)
	}
}
