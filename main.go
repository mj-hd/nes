package main

import (
	"os"
)

func main() {

	nes, err := NewNES(os.Args[1])
	if err != nil {
		panic(err)
	}

	//	for i := uint16(0); i < uint16(0xFFFF); i++ {
	//		if i%16 == 0 {
	//			fmt.Printf("\n %04x: ", i)
	//		}
	//		fmt.Printf("%02x ", nes.CPU.get(i))
	//	}

	nes.PowerOn()
	nes.Run()

}
