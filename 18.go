package main

import (
	"bufio"
	"log"
	"strings"
)

func init() {
	addSolutions(18, (*problemContext).problem18)
}

func (ctx *problemContext) problem18() {
	var lines [][]byte
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		lines = append(lines, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	ls := newLandscape(lines)
	for i := 0; i < 10; i++ {
		ls.step()
	}
	ctx.reportPart1(ls.resourceValue())

	ls = newLandscape(lines)
	for {
		loop := ls.step()
		if loop >= 0 {
			rem := (1e9 - ls.gen) % loop
			for i := 0; i < rem; i++ {
				ls.step()
			}
			ctx.reportPart2(ls.resourceValue())
			break
		}
	}
}

type landscape struct {
	acres  [][]byte
	states map[string]genValue
	gen    int
	w, h   int
}

type genValue struct {
	gen int
	val int
}

func newLandscape(acres [][]byte) *landscape {
	ls := &landscape{
		acres:  copyByteMatrix(acres),
		states: make(map[string]genValue),
		w:      len(acres[0]),
		h:      len(acres),
	}
	ls.states[ls.String()] = genValue{gen: 0, val: ls.resourceValue()}
	return ls
}

func copyByteMatrix(m [][]byte) [][]byte {
	m1 := make([][]byte, len(m))
	for i, line := range m {
		m1[i] = append([]byte(nil), line...)
	}
	return m1
}

func (ls *landscape) step() (loop int) {
	nextAcres := copyByteMatrix(ls.acres)
	for y := 0; y < ls.h; y++ {
		for x := 0; x < ls.w; x++ {
			c := ls.acres[y][x]
			switch c {
			case '.':
				if ls.countNeighbors(x, y, '|') >= 3 {
					nextAcres[y][x] = '|'
				}
			case '|':
				if ls.countNeighbors(x, y, '#') >= 3 {
					nextAcres[y][x] = '#'
				}
			case '#':
				if ls.countNeighbors(x, y, '#') < 1 ||
					ls.countNeighbors(x, y, '|') < 1 {
					nextAcres[y][x] = '.'
				}
			}
		}
	}
	ls.acres = nextAcres
	ls.gen++
	s := ls.String()
	if gv, ok := ls.states[s]; ok {
		return ls.gen - gv.gen
	}
	ls.states[s] = genValue{gen: ls.gen, val: ls.resourceValue()}
	return -1
}

func (ls *landscape) String() string {
	var b strings.Builder
	for _, line := range ls.acres {
		b.Write(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func (ls *landscape) countNeighbors(x, y int, c byte) int {
	var total int
	for _, n := range ls.neighbors(x, y) {
		if n == c {
			total++
		}
	}
	return total
}

func (ls *landscape) neighbors(x, y int) []byte {
	xmin := max(x-1, 0)
	xmax := min(x+1, ls.w-1)
	ymin := max(y-1, 0)
	ymax := min(y+1, ls.h-1)
	var ns []byte
	for yy := ymin; yy <= ymax; yy++ {
		for xx := xmin; xx <= xmax; xx++ {
			if xx == x && yy == y {
				continue
			}
			ns = append(ns, ls.acres[yy][xx])
		}
	}
	return ns
}

func (ls *landscape) resourceValue() int {
	var wood, lumber int
	for y := 0; y < ls.h; y++ {
		for x := 0; x < ls.w; x++ {
			switch ls.acres[y][x] {
			case '|':
				wood++
			case '#':
				lumber++
			}
		}
	}
	return wood * lumber
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
