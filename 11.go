package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"strconv"
)

func init() {
	addSolutions(11, (*problemContext).problem11)
}

func (ctx *problemContext) problem11() {
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	s := string(bytes.TrimSpace(b))
	serial, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	powerLevelSerial = serial
	for y := int64(1); y <= 300; y++ {
		for x := int64(1); x <= 300; x++ {
			cell := vec2{x, y}
			powerLevelCache[powerIndex(cell, 1)] = powerLevel(cell)
		}
	}
	ctx.reportLoad()

	var max int64
	var best vec2
	for y := int64(1); y <= 300-3+1; y++ {
		for x := int64(1); x <= 300-3+1; x++ {
			cell := vec2{x, y}
			if p := powerLevelSquare(cell, 3); p > max {
				max = p
				best = cell
			}
		}
	}
	ctx.reportPart1(best)

	var bestSize int64
	for ctx.loopForProfile() {
		for size := int64(1); size <= 300; size++ {
			for y := int64(1); y <= 300-size+1; y++ {
				for x := int64(1); x <= 300-size+1; x++ {
					cell := vec2{x, y}
					if p := powerLevelSquare(cell, size); p > max {
						max = p
						best = cell
						bestSize = size
					}
				}
			}
		}
	}
	ctx.reportPart2(best, bestSize)
}

var powerLevelSerial int64

var (
	powerLevelCache [300 * 300 * 300]int64
)

func init() {
	for i := range powerLevelCache {
		powerLevelCache[i] = math.MinInt64
	}
}

func powerIndex(cell vec2, size int64) int64 {
	return (size-1)*300*300 + (cell.y-1)*300 + cell.x - 1
}

func powerLevel(cell vec2) int64 {
	rackID := cell.x + 10
	p := rackID * cell.y
	p += powerLevelSerial
	p *= rackID
	p = (p / 100) % 10
	p -= 5
	return p
}

type cellAndSize struct {
	cell vec2
	size int64
}

func powerLevelSquare(cell vec2, size int64) int64 {
	if size == 1 {
		return powerLevelCache[powerIndex(cell, 1)]
	}
	p := powerLevelCache[powerIndex(cell, size)]
	if p != math.MinInt64 {
		return p
	}
	if size%2 == 0 {
		// clockwise
		quad0 := cell
		quad1 := vec2{cell.x + size/2, cell.y}
		quad2 := vec2{cell.x + size/2, cell.y + size/2}
		quad3 := vec2{cell.x, cell.y + size/2}
		p = powerLevelSquare(quad0, size/2) +
			powerLevelSquare(quad1, size/2) +
			powerLevelSquare(quad2, size/2) +
			powerLevelSquare(quad3, size/2)
	} else {
		p = powerLevelSquare(vec2{cell.x + 1, cell.y + 1}, size-1)
		for x := cell.x; x < cell.x+size; x++ {
			p += powerLevelCache[powerIndex(vec2{x, cell.y}, 1)]
		}
		for y := cell.y + 1; y < cell.y+size; y++ {
			p += powerLevelCache[powerIndex(vec2{cell.x, y}, 1)]
		}
	}
	powerLevelCache[powerIndex(cell, size)] = p
	return p
}
