package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (ctx *problemContext) problem3() {
	g := newGrid(1000, 1000)

	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		c, err := parseClaim(scanner.Text())
		if err != nil {
			log.Fatalf("Bad claim %q: %s", scanner.Text(), err)
		}
		g.addClaim(c)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(g.count2())
	fmt.Println(g.findNoOverlaps())
}

type grid struct {
	cells    [][]string
	overlaps map[string]bool
	width    int
	height   int
}

func newGrid(width, height int) *grid {
	return &grid{
		cells:    make([][]string, width*height),
		overlaps: make(map[string]bool),
		width:    width,
		height:   height,
	}
}

func (g *grid) addClaim(c *claim) {
	for x := c.x; x < c.x+c.w; x++ {
		for y := c.y; y < c.y+c.h; y++ {
			i := y*g.width + x
			for _, id := range g.cells[i] {
				g.overlaps[id] = true
			}
			if !g.overlaps[c.id] {
				g.overlaps[c.id] = len(g.cells[i]) > 0
			}
			g.cells[i] = append(g.cells[i], c.id)
		}
	}
}

func (g *grid) count2() int {
	var total int
	for _, ids := range g.cells {
		if len(ids) >= 2 {
			total++
		}
	}
	return total
}

func (g *grid) findNoOverlaps() string {
	for id, overlaps := range g.overlaps {
		if !overlaps {
			return id
		}
	}
	return ""
}

type claim struct {
	id   string
	x, y int
	w, h int
}

func parseClaim(s string) (*claim, error) {
	fields := strings.Fields(s)
	xy := strings.Split(strings.TrimSuffix(fields[2], ":"), ",")
	x, err := strconv.Atoi(xy[0])
	if err != nil {
		return nil, err
	}
	y, err := strconv.Atoi(xy[1])
	if err != nil {
		return nil, err
	}
	wh := strings.Split(fields[3], "x")
	w, err := strconv.Atoi(wh[0])
	if err != nil {
		return nil, err
	}
	h, err := strconv.Atoi(wh[1])
	if err != nil {
		return nil, err
	}
	return &claim{
		id: fields[0],
		x:  x,
		y:  y,
		w:  w,
		h:  h,
	}, nil
}
