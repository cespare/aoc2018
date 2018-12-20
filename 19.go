package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(19, (*problemContext).problem19)
}

func (ctx *problemContext) problem19() {
	var program []instruction2
	scanner := bufio.NewScanner(ctx.f)
	if !scanner.Scan() {
		panic("empty")
	}
	ipReg, err := strconv.Atoi(strings.TrimPrefix(scanner.Text(), "#ip "))
	if err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		inst, err := parseInstruction2(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		program = append(program, inst)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	cpu := &cpu2{
		ipReg: ipReg,
		ins:   program,
	}
	for !cpu.step() {
	}
	ctx.reportPart1(cpu.regs[0])
	ctx.reportPart2(do19part2())
}

type instruction2 struct {
	Op      string
	A, B, C int
}

var instruction2Regexp = regexp.MustCompile(`^(?P<Op>\w+)\s+(?P<A>\d+)\s+(?P<B>\d+)\s+(?P<C>\d+)`)

func parseInstruction2(s string) (instruction2, error) {
	var inst instruction2
	err := hasty.Parse([]byte(s), &inst, instruction2Regexp)
	return inst, err
}

type cpu2 struct {
	ipReg int
	ip    int
	regs  [6]int
	ins   []instruction2

	debug bool
}

func (c *cpu2) step() (halt bool) {
	c.regs[c.ipReg] = c.ip
	if c.ip < 0 || c.ip >= len(c.ins) {
		return true
	}
	c.exec(c.ins[c.ip])
	c.ip = c.regs[c.ipReg]
	c.ip++
	if c.ip < 0 || c.ip >= len(c.ins) {
		return true
	}
	return false
}

func (c *cpu2) exec(in instruction2) {
	if c.debug {
		fmt.Printf("ip=%2d %s %2d %2d %2d %10d\n", c.ip, in.Op, in.A, in.B, in.C, c.regs)
	}
	switch in.Op {
	case "addr":
		c.regs[in.C] = c.regs[in.A] + c.regs[in.B]
	case "addi":
		c.regs[in.C] = c.regs[in.A] + in.B
	case "mulr":
		c.regs[in.C] = c.regs[in.A] * c.regs[in.B]
	case "muli":
		c.regs[in.C] = c.regs[in.A] * in.B
	case "banr":
		c.regs[in.C] = c.regs[in.A] & c.regs[in.B]
	case "bani":
		c.regs[in.C] = c.regs[in.A] & in.B
	case "borr":
		c.regs[in.C] = c.regs[in.A] | c.regs[in.B]
	case "bori":
		c.regs[in.C] = c.regs[in.A] | in.B
	case "setr":
		c.regs[in.C] = c.regs[in.A]
	case "seti":
		c.regs[in.C] = in.A
	case "gtir":
		c.regs[in.C] = boolToInt(in.A > c.regs[in.B])
	case "gtri":
		c.regs[in.C] = boolToInt(c.regs[in.A] > in.B)
	case "gtrr":
		c.regs[in.C] = boolToInt(c.regs[in.A] > c.regs[in.B])
	case "eqir":
		c.regs[in.C] = boolToInt(in.A == c.regs[in.B])
	case "eqri":
		c.regs[in.C] = boolToInt(c.regs[in.A] == in.B)
	case "eqrr":
		c.regs[in.C] = boolToInt(c.regs[in.A] == c.regs[in.B])
	}
}

func do19part2() int {
	// Simplified version of the asm.
	x := 0
	for y := 1; y <= 10551373; y++ {
		if 10551373%y == 0 {
			x += y
			continue
		}
	}
	return x
}
