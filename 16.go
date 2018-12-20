package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(16, (*problemContext).problem16)
}

func (ctx *problemContext) problem16() {
	var samples []*cpuSample
	var program []instruction
	var sampleLines []string
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if len(program) > 0 ||
			(len(sampleLines) == 0 && !strings.HasPrefix(line, "Before:")) {
			in, err := parseInstruction(line)
			if err != nil {
				log.Fatal(err)
			}
			program = append(program, in)
			continue
		}
		sampleLines = append(sampleLines, line)
		if len(sampleLines) == 3 {
			lines := [3]string{sampleLines[0], sampleLines[1], sampleLines[2]}
			sample, err := parseCPUSample(lines)
			if err != nil {
				log.Fatal(err)
			}
			samples = append(samples, sample)
			sampleLines = nil
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if len(sampleLines) > 0 {
		log.Fatal("Incomplete CPU sample")
	}

	opcodeMapping := make(map[uint8]map[uint8]struct{})
	for i := uint8(0); i < uint8(len(instructionFuncs)); i++ {
		m := make(map[uint8]struct{})
		for j := uint8(0); j < uint8(len(instructionFuncs)); j++ {
			m[j] = struct{}{}
		}
		opcodeMapping[i] = m
	}

	var part1 int
	for _, sample := range samples {
		matching := sample.before.matchingOpcodes(sample.in, sample.after)
		if len(matching) >= 3 {
			part1++
		}
		old := opcodeMapping[sample.in.op]
		if len(old) == 1 {
			continue // done
		}
		m := make(map[uint8]struct{})
		for _, i := range matching {
			if _, ok := old[i]; ok {
				m[i] = struct{}{}
			}
		}
		opcodeMapping[sample.in.op] = m
	}
	ctx.reportPart1(part1)

	for len(opcodeMapping) > 0 {
		for op, m := range opcodeMapping {
			if len(m) > 1 {
				continue
			}
			var i uint8
			for ii := range m {
				i = ii
			}
			opcodeToFunc[op] = instructionFuncs[i]
			delete(opcodeMapping, op)
			for _, m1 := range opcodeMapping {
				delete(m1, i)
			}
		}
	}

	regs := registers{0, 0, 0, 0}
	for _, in := range program {
		regs = regs.eval(in)
	}
	ctx.reportPart2(regs[0])
}

type instruction struct {
	op uint8
	a  int
	b  int
	c  int
}

func parseInstruction(s string) (instruction, error) {
	var in instruction
	parts := strings.Fields(s)
	if len(parts) != 4 {
		return in, fmt.Errorf("bad instruction: %q", s)
	}
	n, err := strconv.ParseUint(parts[0], 10, 8)
	if err != nil {
		return in, err
	}
	in.op = uint8(n)
	in.a, err = strconv.Atoi(parts[1])
	if err != nil {
		return in, err
	}
	in.b, err = strconv.Atoi(parts[2])
	if err != nil {
		return in, err
	}
	in.c, err = strconv.Atoi(parts[3])
	if err != nil {
		return in, err
	}
	return in, nil
}

type registers [4]int

func parseRegisters(s string) (registers, error) {
	s = strings.TrimSpace(s)
	var regs registers
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	parts := strings.Split(s, ", ")
	if len(parts) != 4 {
		return regs, fmt.Errorf("bad registers: %q", s)
	}
	for i, rs := range parts {
		var err error
		regs[i], err = strconv.Atoi(rs)
		if err != nil {
			return regs, err
		}
	}
	return regs, nil
}

type cpuSample struct {
	in            instruction
	before, after registers
}

func parseCPUSample(ss [3]string) (*cpuSample, error) {
	var sample cpuSample
	var err error
	sample.before, err = parseRegisters(strings.TrimPrefix(ss[0], "Before: "))
	if err != nil {
		return nil, err
	}
	sample.in, err = parseInstruction(ss[1])
	if err != nil {
		return nil, err
	}
	sample.after, err = parseRegisters(strings.TrimPrefix(ss[2], "After: "))
	if err != nil {
		return nil, err
	}
	return &sample, nil
}

func (r registers) eval(in instruction) registers {
	return opcodeToFunc[in.op](r, in.a, in.b, in.c)
}

func (r registers) matchingOpcodes(in instruction, result registers) []uint8 {
	var matching []uint8
	for i, fn := range instructionFuncs {
		if fn(r, in.a, in.b, in.c) == result {
			matching = append(matching, uint8(i))
		}
	}
	return matching
}

var instructionFuncs = []func(registers, int, int, int) registers{
	registers.addr,
	registers.addi,
	registers.mulr,
	registers.muli,
	registers.banr,
	registers.bani,
	registers.borr,
	registers.bori,
	registers.setr,
	registers.seti,
	registers.gtir,
	registers.gtri,
	registers.gtrr,
	registers.eqir,
	registers.eqri,
	registers.eqrr,
}

var opcodeToFunc = make([]func(registers, int, int, int) registers, len(instructionFuncs))

func (r registers) addr(a, b, c int) registers {
	r[c] = r[a] + r[b]
	return r
}

func (r registers) addi(a, b, c int) registers {
	r[c] = r[a] + b
	return r
}

func (r registers) mulr(a, b, c int) registers {
	r[c] = r[a] * r[b]
	return r
}

func (r registers) muli(a, b, c int) registers {
	r[c] = r[a] * b
	return r
}

func (r registers) banr(a, b, c int) registers {
	r[c] = r[a] & r[b]
	return r
}

func (r registers) bani(a, b, c int) registers {
	r[c] = r[a] & b
	return r
}

func (r registers) borr(a, b, c int) registers {
	r[c] = r[a] | r[b]
	return r
}

func (r registers) bori(a, b, c int) registers {
	r[c] = r[a] | b
	return r
}

func (r registers) setr(a, _, c int) registers {
	r[c] = r[a]
	return r
}

func (r registers) seti(a, _, c int) registers {
	r[c] = a
	return r
}

func (r registers) gtir(a, b, c int) registers {
	r[c] = boolToInt(a > r[b])
	return r
}

func (r registers) gtri(a, b, c int) registers {
	r[c] = boolToInt(r[a] > b)
	return r
}

func (r registers) gtrr(a, b, c int) registers {
	r[c] = boolToInt(r[a] > r[b])
	return r
}

func (r registers) eqir(a, b, c int) registers {
	r[c] = boolToInt(a == r[b])
	return r
}

func (r registers) eqri(a, b, c int) registers {
	r[c] = boolToInt(r[a] == b)
	return r
}

func (r registers) eqrr(a, b, c int) registers {
	r[c] = boolToInt(r[a] == r[b])
	return r
}

func boolToInt(t bool) int {
	if t {
		return 1
	}
	return 0
}
