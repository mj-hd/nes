package main

func main() {

	nes := NewNES()
	err := nes.Load("./sample.nes")
	if err != nil {
		panic(err)
	}

	nes.PowerOn()
	nes.Run()

}
