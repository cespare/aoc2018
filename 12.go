package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

func init() {
	addSolutions(12, (*problemContext).problem12)
}

func (ctx *problemContext) problem12() {
	scanner := bufio.NewScanner(ctx.f)
	if !scanner.Scan() {
		panic("bad input")
	}
	initial := strings.TrimPrefix(scanner.Text(), "initial state: ")
	rules := make(map[string]struct{})
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}
		rule, c, err := parsePlantRule(s)
		if err != nil {
			log.Fatal(err)
		}
		if c == '#' {
			rules[rule] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	state := newPlantState(initial)
	for i := 0; i < 20; i++ {
		state = state.step(rules)
	}
	ctx.reportPart1(state.score())

	state = newPlantState(initial)
	gen := int64(1)
	var delta int64
	for ; ; gen++ {
		prev := state
		state = state.step(rules)
		if state.plants == prev.plants {
			delta = state.offset - prev.offset
			break
		}
	}
	// Skip a few.
	state.offset += delta * (50e9 - gen)
	ctx.reportPart2(state.score())
}

type plantState struct {
	plants string
	offset int64
}

func newPlantState(initial string) *plantState {
	s := &plantState{
		plants: "....." + initial + ".....",
		offset: -5,
	}
	s.trim()
	return s
}

// trim chops off leading and trailing '.' until there are exactly 5
// at each end. (It assumes there are already at least 5.)
func (s *plantState) trim() {
	i := strings.LastIndexByte(s.plants, '#')
	if i == -1 {
		panic("empty")
	}
	s.plants = s.plants[:i+1+5]
	i = strings.IndexByte(s.plants, '#')
	s.plants = s.plants[i-5:]
	s.offset += int64(i - 5)
}

func (s *plantState) step(rules map[string]struct{}) *plantState {
	var b strings.Builder
	b.WriteString(".....")
	offset := s.offset - 3
	for i := 2; i < len(s.plants)-2; i++ {
		if _, ok := rules[s.plants[i-2:i+3]]; ok {
			b.WriteByte('#')
		} else {
			b.WriteByte('.')
		}
	}
	b.WriteString(".....")
	s1 := &plantState{
		plants: b.String(),
		offset: offset,
	}
	s1.trim()
	return s1
}

func (s *plantState) score() int64 {
	var score int64
	for i := 0; i < len(s.plants); i++ {
		if s.plants[i] == '#' {
			score += int64(i) + s.offset
		}
	}
	return score
}

func (s *plantState) string(offset int) string {
	if s.offset < int64(offset) {
		panic("need smaller offset")
	}
	var b strings.Builder
	for i := int64(offset); i < s.offset; i++ {
		b.WriteByte('.')
	}
	b.WriteString(s.plants)
	return b.String()
}

func parsePlantRule(s string) (string, byte, error) {
	parts := strings.SplitN(s, " => ", 2)
	if len(parts) != 2 || len(parts[0]) != 5 || len(parts[1]) != 1 {
		return "", 0, fmt.Errorf("bad rule: %q", s)
	}
	return parts[0], parts[1][0], nil
}
