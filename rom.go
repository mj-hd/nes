package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
)

type rom struct {
	Header      Header
	Trainer     [512]byte
	PRG         []byte
	CHR         []byte
	PC10InstROM []byte
	PC10PROM    []byte
}

type Header struct {
	MapperNum int
}

// NES 2.0
func (r *rom) Load(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	buf := bytes.NewReader(data)

	log.Println("start loading...")

	n, _ := buf.ReadByte()
	e, _ := buf.ReadByte()
	s, _ := buf.ReadByte()
	f, _ := buf.ReadByte()

	if !(0x4E == n &&
		0x45 == e &&
		0x53 == s &&
		0x1A == f) {
		return errors.New("unknown format")
	}

	prg_size, _ := buf.ReadByte()
	chr_size, _ := buf.ReadByte()

	log.Printf("PRG: %d, CHR: %d\n", prg_size, chr_size)

	flag1, _ := buf.ReadByte()
	flag2, _ := buf.ReadByte()

	prg_ram_size, _ := buf.ReadByte()
	_ = prg_ram_size

	mapper_num := flag1 & 0xF0 >> 4
	four_screen := flag1 & 0x08 >> 3
	trainer := flag1 & 0x04 >> 2
	battery := flag1 & 0x02 >> 1
	mirroring := flag1 & 0x01
	_ = four_screen
	_ = battery
	_ = mirroring

	mapper_num += flag2 & 0xF0
	pc10 := flag2 & 0x02 >> 1
	vs := flag2 & 0x01
	_ = vs

	log.Printf("MAPPER: %d\n", mapper_num)
	r.Header.MapperNum = int(mapper_num)

	if flag2&0x0C == 0x08 {
		flag3, _ := buf.ReadByte()
		flag4, _ := buf.ReadByte()

		pal := flag3 & 0x08 >> 3
		bus_conflicts := flag4 & 0x80 >> 7
		prg_ram := flag4 & 0x40 >> 6
		tv_system := flag4 & 0x30 >> 4
		_ = pal
		_ = bus_conflicts
		_ = prg_ram
		_ = tv_system
	}

	buf.Seek(16, io.SeekStart)

	if trainer == 1 {
		io.ReadFull(buf, r.Trainer[:])
	}

	r.PRG = make([]byte, int64(prg_size)*16*1024)
	r.CHR = make([]byte, int64(chr_size)*8*1024)

	io.ReadFull(buf, r.PRG)
	io.ReadFull(buf, r.CHR)

	if pc10 == 1 {
		io.ReadFull(buf, r.PC10InstROM)
	}

	if buf.Len() >= 16 {
		io.ReadFull(buf, r.PC10PROM)
	}

	return nil
}

func (r *rom) GetCHR(addr uint16) byte {
	if addr >= uint16(len(r.CHR)) {
		return 0
	}
	return r.CHR[addr]
}

func (r *rom) SetCHR(addr uint16, value byte) {
	if addr >= uint16(len(r.CHR)) {
		return
	}
	r.CHR[addr] = value
}

func (r *rom) GetPRG(addr uint16) byte {
	if addr >= uint16(len(r.PRG)) {
		return 0
	}
	return r.PRG[addr]
}

func (r *rom) SetPRG(addr uint16, value byte) {
	if addr >= uint16(len(r.PRG)) {
		return
	}
	r.PRG[addr] = value
}
