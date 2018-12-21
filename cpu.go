package main

import "fmt"

const CPUFrequency = 1789773
const CPUStackStart = 0x0100

// interrupt types
const (
	_ = iota
	interruptNone
	interruptNMI
	interruptIRQ
	interruptBRK
)

// addressing modes
const (
	_ = iota
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeAccumulator
	modeImmediate
	modeImplied
	modeIndexedIndirect
	modeIndirect
	modeIndirectIndexed
	modeRelative
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
)

// instruction_modes indicates the addressing mode for each instruction
var instruction_modes = [256]int{
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	1, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 8, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
}

// instruction_sizes indicates the size of each instruction in bytes
var instruction_sizes = [256]int{
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
}

// instruction_cycles indicates the number of cycles used by each instruction,
// not including conditional cycles
var instruction_cycles = [256]int{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

// instruction_page_cycles indicates the number of cycles used by each
// instruction when a page is crossed
var instruction_page_cycles = [256]int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
}

// instructionNames indicates the name of each instruction
var instruction_names = [256]string{
	"BRK", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"PHP", "ORA", "ASL", "ANC", "NOP", "ORA", "ASL", "SLO",
	"BPL", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"CLC", "ORA", "NOP", "SLO", "NOP", "ORA", "ASL", "SLO",
	"JSR", "AND", "KIL", "RLA", "BIT", "AND", "ROL", "RLA",
	"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA",
	"BMI", "AND", "KIL", "RLA", "NOP", "AND", "ROL", "RLA",
	"SEC", "AND", "NOP", "RLA", "NOP", "AND", "ROL", "RLA",
	"RTI", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
	"BVC", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"CLI", "EOR", "NOP", "SRE", "NOP", "EOR", "LSR", "SRE",
	"RTS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"PLA", "ADC", "ROR", "ARR", "JMP", "ADC", "ROR", "RRA",
	"BVS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA",
	"NOP", "STA", "NOP", "SAX", "STY", "STA", "STX", "SAX",
	"DEY", "NOP", "TXA", "XAA", "STY", "STA", "STX", "SAX",
	"BCC", "STA", "KIL", "AHX", "STY", "STA", "STX", "SAX",
	"TYA", "STA", "TXS", "TAS", "SHY", "STA", "SHX", "AHX",
	"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"TAY", "LDA", "TAX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"BCS", "LDA", "KIL", "LAX", "LDY", "LDA", "LDX", "LAX",
	"CLV", "LDA", "TSX", "LAS", "LDY", "LDA", "LDX", "LAX",
	"CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
	"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP",
	"BNE", "CMP", "KIL", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CLD", "CMP", "NOP", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CPX", "SBC", "NOP", "ISC", "CPX", "SBC", "INC", "ISC",
	"INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
	"BEQ", "SBC", "KIL", "ISC", "NOP", "SBC", "INC", "ISC",
	"SED", "SBC", "NOP", "ISC", "NOP", "SBC", "INC", "ISC",
}

// status flags
const (
	_ = iota
	FlagC
	FlagZ
	FlagI
	FlagD
	FlagB
	FlagR
	FlagV
	FlagN
)

var instructions = [256]func(*cpu, uint16, int){
	brk, ora, kil, slo, nop, ora, asl, slo,
	php, ora, asl, anc, nop, ora, asl, slo,
	bpl, ora, kil, slo, nop, ora, asl, slo,
	clc, ora, nop, slo, nop, ora, asl, slo,
	jsr, and, kil, rla, bit, and, rol, rla,
	plp, and, rol, anc, bit, and, rol, rla,
	bmi, and, kil, rla, nop, and, rol, rla,
	sec, and, nop, rla, nop, and, rol, rla,
	rti, eor, kil, sre, nop, eor, lsr, sre,
	pha, eor, lsr, alr, jmp, eor, lsr, sre,
	bvc, eor, kil, sre, nop, eor, lsr, sre,
	cli, eor, nop, sre, nop, eor, lsr, sre,
	rts, adc, kil, rra, nop, adc, ror, rra,
	pla, adc, ror, arr, jmp, adc, ror, rra,
	bvs, adc, kil, rra, nop, adc, ror, rra,
	sei, adc, nop, rra, nop, adc, ror, rra,
	nop, sta, nop, sax, sty, sta, stx, sax,
	dey, nop, txa, xaa, sty, sta, stx, sax,
	bcc, sta, kil, ahx, sty, sta, stx, sax,
	tya, sta, txs, tas, shy, sta, shx, ahx,
	ldy, lda, ldx, lax, ldy, lda, ldx, lax,
	tay, lda, tax, lax, ldy, lda, ldx, lax,
	bcs, lda, kil, lax, ldy, lda, ldx, lax,
	clv, lda, tsx, las, ldy, lda, ldx, lax,
	cpy, cmp, nop, dcp, cpy, cmp, dec, dcp,
	iny, cmp, dex, axs, cpy, cmp, dec, dcp,
	bne, cmp, kil, dcp, nop, cmp, dec, dcp,
	cld, cmp, nop, dcp, nop, cmp, dec, dcp,
	cpx, sbc, nop, isc, cpx, sbc, inc, isc,
	inx, sbc, nop, sbc, cpx, sbc, inc, isc,
	beq, sbc, kil, isc, nop, sbc, inc, isc,
	sed, sbc, nop, isc, nop, sbc, inc, isc,
}

type cpu struct {
	A           byte
	X           byte
	Y           byte
	S           byte
	P           byte
	PC          uint16
	MMC         mmc
	wait_cycles int
	interrupt   int
}

func (c *cpu) Tick() {
	if c.wait_cycles > 0 {
		c.wait_cycles--
		return
	}

	switch c.interrupt {
	case interruptNMI:
		c.pushAddress(c.PC)
		c.push(c.P)
		c.setStatusFlag(FlagI, true)
		c.setStatusFlag(FlagB, false)
		c.PC = c.MMC.GetAddress(0xFFFA)
	case interruptIRQ:
		if !c.getStatusFlagBool(FlagI) {
			c.pushAddress(c.PC)
			c.push(c.P)
			c.setStatusFlag(FlagI, true)
			c.setStatusFlag(FlagB, false)
			c.PC = c.MMC.GetAddress(0xFFFE)
		}
	case interruptBRK:
		if !c.getStatusFlagBool(FlagI) {
			c.setStatusFlag(FlagB, true)
			c.PC += 1
			c.pushAddress(c.PC)
			c.push(c.P)
			c.setStatusFlag(FlagI, true)
			c.PC = c.MMC.GetAddress(0xFFFE)
		}
	}
	if c.interrupt != interruptNone {
		c.wait_cycles += 7
	}
	c.interrupt = interruptNone

	opecode := c.MMC.Get(c.PC)
	c.PC += 1

	mode := instruction_modes[opecode]

	var address uint16
	page_crossed := false
	switch mode {
	case modeAbsolute:
		address = c.MMC.GetAddress(c.PC)
	case modeAbsoluteX:
		address = c.MMC.GetAddress(c.PC) + uint16(c.X)
		page_crossed = pages_differ(address-uint16(c.X), address)
	case modeAbsoluteY:
		address = c.MMC.GetAddress(c.PC) + uint16(c.Y)
		page_crossed = pages_differ(address-uint16(c.Y), address)
	case modeAccumulator:
		address = 0
	case modeImmediate:
		address = c.PC
	case modeImplied:
		address = 0
	case modeIndexedIndirect:
		address = c.MMC.GetAddress(uint16(c.MMC.Get(c.PC) + c.X))
	case modeIndirect:
		address = c.MMC.GetAddress(c.MMC.GetAddress(c.PC))
	case modeIndirectIndexed:
		address = c.MMC.GetAddress(uint16(c.MMC.Get(c.PC))) + uint16(c.Y)
		page_crossed = pages_differ(address-uint16(c.Y), address)
	case modeRelative:
		offset := uint16(c.MMC.Get(c.PC))
		if offset < 0x80 {
			address = c.PC + 2 + offset
		} else {
			address = c.PC + 2 + offset - 0x100
		}
	case modeZeroPage:
		address = uint16(c.MMC.Get(c.PC))
	case modeZeroPageX:
		address = uint16(c.MMC.Get(c.PC) + c.X)
	case modeZeroPageY:
		address = uint16(c.MMC.Get(c.PC) + c.Y)
	}

	c.PC += uint16(instruction_sizes[opecode] - 1)
	c.wait_cycles += instruction_cycles[opecode]
	if page_crossed {
		c.wait_cycles += instruction_page_cycles[opecode]
	}

	instructions[opecode](c, address, mode)
	fmt.Println("PC:", c.PC, " OP:", instruction_names[opecode])
}

func (c *cpu) PowerOn() {
	c.A = 0
	c.X = 0
	c.Y = 0
	c.S = 0xFD
	c.P = FlagZ | FlagR
	c.PC = c.MMC.GetAddress(0xFFFC)
	//c.PC = 0x8000
	c.wait_cycles = 0
	// TODO: do not relay on mmc1
	c.MMC.Set(MMC1AddressAPUFrameCounter, 0x00)
	c.MMC.Set(MMC1AddressAPUStatus, 0x00)
	c.MMC.Set(MMC1AddressAPUPulse1, 0x00)
}

func (c *cpu) Reset() {
	c.S -= 0x03
	c.P |= FlagI
	c.PC = c.MMC.GetAddress(0xFFFC)
	// TODO: do not relay on mmc1
	c.MMC.Set(MMC1AddressAPUStatus, 0x00)
}

func (c *cpu) push(value uint8) {
	c.MMC.Set(CPUStackStart+uint16(c.S), value)
	c.S--
}

func (c *cpu) pushAddress(value uint16) {
	c.MMC.SetAddress(CPUStackStart+uint16(c.S)-1, value)
	c.S -= 2
}

func (c *cpu) pop() uint8 {
	c.S++
	return c.MMC.Get(CPUStackStart + uint16(c.S))
}

func (c *cpu) popAddress() uint16 {
	c.S++
	v := c.MMC.GetAddress(CPUStackStart + uint16(c.S))
	c.S++
	return v
}

func (c *cpu) setStatusFlag(flag int, value bool) {
	if value {
		c.P |= byte(flag)
	} else {
		c.P &= ^byte(flag)
	}
}

func (c *cpu) getStatusFlagBool(flag int) bool {
	return c.P&byte(flag) != 0
}

func (c *cpu) getStatusFlagByte(flag int) byte {
	if c.getStatusFlagBool(flag) {
		return 1
	} else {
		return 0
	}
}

func (c *cpu) setZN(value byte) {
	c.setStatusFlag(FlagN, is_negative(value))
	c.setStatusFlag(FlagZ, value == 0)
}

func (c *cpu) add_branch_wait(addr uint16) {
	c.PC++
	if pages_differ(c.PC, addr) {
		c.PC++
	}
}

func pages_differ(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

func is_negative(a byte) bool {
	return a&0x80 != 0
}

// TODO: メンバにした方がいい
func adc(c *cpu, addr uint16, mode int) {
	v1 := c.A
	v2 := c.MMC.Get(addr)
	v3 := c.getStatusFlagByte(FlagC)
	c.A = v1 + v2 + v3

	c.setZN(c.A)

	if !is_negative(v1^v2) && is_negative(v1^c.A) {
		c.setStatusFlag(FlagV, true)
	} else {
		c.setStatusFlag(FlagV, false)
	}

	if int(v1)+int(v2)+int(v3) > 0xFF {
		c.setStatusFlag(FlagC, true)
	} else {
		c.setStatusFlag(FlagC, false)
	}
}
func and(c *cpu, addr uint16, mode int) {
	c.A &= c.MMC.Get(addr)
	c.setZN(c.A)
}
func asl(c *cpu, addr uint16, mode int) {
	var v uint8
	if mode == modeAccumulator {
		v = c.A
	} else {
		v = c.MMC.Get(addr)
	}

	c.setStatusFlag(FlagC, is_negative(v))

	v <<= 1

	if mode == modeAccumulator {
		c.A = v
	} else {
		c.MMC.Set(addr, v)
	}

	c.setZN(v)
}
func bcc(c *cpu, addr uint16, mode int) {
	if !c.getStatusFlagBool(FlagC) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func bcs(c *cpu, addr uint16, mode int) {
	if c.getStatusFlagBool(FlagC) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func beq(c *cpu, addr uint16, mode int) {
	if c.getStatusFlagBool(FlagZ) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func bit(c *cpu, addr uint16, mode int) {
	v := c.MMC.Get(addr)
	c.setStatusFlag(FlagV, v>>6&1 == 1)
	c.setStatusFlag(FlagN, is_negative(v))
	c.setStatusFlag(FlagZ, v&c.A == 0)
}
func bmi(c *cpu, addr uint16, mode int) {
	if c.getStatusFlagBool(FlagN) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func bne(c *cpu, addr uint16, mode int) {
	if !c.getStatusFlagBool(FlagZ) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func bpl(c *cpu, addr uint16, mode int) {
	if !c.getStatusFlagBool(FlagN) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func brk(c *cpu, addr uint16, mode int) {
	c.interrupt = interruptBRK
}
func bvc(c *cpu, addr uint16, mode int) {
	if !c.getStatusFlagBool(FlagV) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func bvs(c *cpu, addr uint16, mode int) {
	if c.getStatusFlagBool(FlagV) {
		c.PC = addr
		c.add_branch_wait(addr)
	}
}
func clc(c *cpu, addr uint16, mode int) {
	c.setStatusFlag(FlagC, false)
}
func cld(c *cpu, addr uint16, mode int) {
	// not implemented
}
func cli(c *cpu, addr uint16, mode int) {
	c.setStatusFlag(FlagI, false)
}
func clv(c *cpu, addr uint16, mode int) {
	c.setStatusFlag(FlagV, false)
}
func cmp(c *cpu, addr uint16, mode int) {
	v1 := c.A
	v2 := c.MMC.Get(addr)
	v3 := v1 - v2
	c.setZN(v3)
	c.setStatusFlag(FlagC, !is_negative(v3))
}
func cpx(c *cpu, addr uint16, mode int) {
	v1 := c.X
	v2 := c.MMC.Get(addr)
	v3 := v1 - v2
	c.setZN(v3)
	c.setStatusFlag(FlagC, !is_negative(v3))
}
func cpy(c *cpu, addr uint16, mode int) {
	v1 := c.Y
	v2 := c.MMC.Get(addr)
	v3 := v1 - v2
	c.setZN(v3)
	c.setStatusFlag(FlagC, !is_negative(v3))
}
func dec(c *cpu, addr uint16, mode int) {
	v := c.MMC.Get(addr) - 1
	c.setZN(v)
	c.MMC.Set(addr, v)
}
func dex(c *cpu, addr uint16, mode int) {
	v := c.X - 1
	c.setZN(v)
	c.X = v
}
func dey(c *cpu, addr uint16, mode int) {
	v := c.Y - 1
	c.setZN(v)
	c.Y = v
}
func eor(c *cpu, addr uint16, mode int) {
	v1 := c.A
	v2 := c.MMC.Get(addr)
	v3 := v1 ^ v2
	c.setZN(v3)
	c.A = v3
}
func inc(c *cpu, addr uint16, mode int) {
	v := c.MMC.Get(addr) + 1
	c.setZN(v)
	c.MMC.Set(addr, v)
}
func inx(c *cpu, addr uint16, mode int) {
	v := c.X + 1
	c.setZN(v)
	c.X = v
}
func iny(c *cpu, addr uint16, mode int) {
	v := c.Y + 1
	c.setZN(v)
	c.Y = v
}
func jmp(c *cpu, addr uint16, mode int) {
	c.PC = addr
}
func jsr(c *cpu, addr uint16, mode int) {
	c.pushAddress(c.PC - 1)
	c.PC = addr
}
func lda(c *cpu, addr uint16, mode int) {
	v := c.MMC.Get(addr)
	c.setZN(v)
	c.A = v
}
func ldx(c *cpu, addr uint16, mode int) {
	v := c.MMC.Get(addr)
	c.setZN(v)
	c.X = v
}
func ldy(c *cpu, addr uint16, mode int) {
	v := c.MMC.Get(addr)
	c.setZN(v)
	c.Y = v
}
func lsr(c *cpu, addr uint16, mode int) {
	var v uint8
	if mode == modeAccumulator {
		v = c.A
	} else {
		v = c.MMC.Get(addr)
	}

	c.setStatusFlag(FlagC, v&1 == 1)
	v >>= 1

	if mode == modeAccumulator {
		c.A = v
	} else {
		c.MMC.Set(addr, v)
	}
	c.setZN(v)
}
func nop(c *cpu, addr uint16, mode int) {
	// nopping
}
func ora(c *cpu, addr uint16, mode int) {
	v1 := c.A
	v2 := c.MMC.Get(addr)
	v3 := v1 | v2
	c.setZN(v3)
	c.A = v3
}
func pha(c *cpu, addr uint16, mode int) {
	c.push(c.A)
}
func php(c *cpu, addr uint16, mode int) {
	c.push(c.P)
}
func pla(c *cpu, addr uint16, mode int) {
	v := c.pop()
	c.setZN(v)
	c.A = v
}
func plp(c *cpu, addr uint16, mode int) {
	c.P = c.pop()
}
func rol(c *cpu, addr uint16, mode int) {
	var v uint8
	if mode == modeAccumulator {
		v = c.A
	} else {
		v = c.MMC.Get(addr)
	}

	carry := v & 0x80 >> 7
	v <<= 1
	v |= carry

	c.setStatusFlag(FlagC, carry == 1)
	c.setZN(v)

	if mode == modeAccumulator {
		c.A = v
	} else {
		c.MMC.Set(addr, v)
	}
}
func ror(c *cpu, addr uint16, mode int) {
	var v uint8
	if mode == modeAccumulator {
		v = c.A
	} else {
		v = c.MMC.Get(addr)
	}

	carry := v & 1
	v >>= 1
	v |= carry << 7

	c.setStatusFlag(FlagC, carry == 1)
	c.setZN(v)

	if mode == modeAccumulator {
		c.A = v
	} else {
		c.MMC.Set(addr, v)
	}
}
func rti(c *cpu, addr uint16, mode int) {
	c.P = c.pop()
	c.PC = c.popAddress()
}
func rts(c *cpu, addr uint16, mode int) {
	c.PC = c.popAddress() + 1
}
func sbc(c *cpu, addr uint16, mode int) {
	v1 := c.A
	v2 := c.MMC.Get(addr)
	v3 := 1 - c.getStatusFlagByte(FlagC)
	c.A = v1 - v2 - v3

	c.setZN(c.A)

	if !is_negative(v1^v2) && is_negative(v1^c.A) {
		c.setStatusFlag(FlagV, true)
	} else {
		c.setStatusFlag(FlagV, false)
	}

	if int(v1)-int(v2)-int(v3) < 0x00 {
		c.setStatusFlag(FlagC, true)
	} else {
		c.setStatusFlag(FlagC, false)
	}
}
func sec(c *cpu, addr uint16, mode int) {
	c.setStatusFlag(FlagC, true)
}
func sed(c *cpu, addr uint16, mode int) {
	// not implemented
}
func sei(c *cpu, addr uint16, mode int) {
	c.setStatusFlag(FlagI, true)
}
func sta(c *cpu, addr uint16, mode int) {
	c.MMC.Set(addr, c.A)
}
func stx(c *cpu, addr uint16, mode int) {
	c.MMC.Set(addr, c.X)
}
func sty(c *cpu, addr uint16, mode int) {
	c.MMC.Set(addr, c.Y)
}
func tax(c *cpu, addr uint16, mode int) {
	c.X = c.A
	c.setZN(c.A)
}
func tay(c *cpu, addr uint16, mode int) {
	c.Y = c.A
	c.setZN(c.A)
}
func tsx(c *cpu, addr uint16, mode int) {
	c.X = c.S
	c.setZN(c.S)
}
func txa(c *cpu, addr uint16, mode int) {
	c.A = c.X
	c.setZN(c.X)
}
func txs(c *cpu, addr uint16, mode int) {
	c.S = c.X
	c.setZN(c.X)
}
func tya(c *cpu, addr uint16, mode int) {
	c.A = c.Y
	c.setZN(c.Y)
}
func ahx(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func alr(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func anc(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func arr(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func axs(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func dcp(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func isc(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func kil(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func las(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func lax(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func rla(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func rra(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func sax(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func shx(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func shy(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func slo(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func sre(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func tas(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
func xaa(c *cpu, addr uint16, mode int) { panic("illegal opecode") }
