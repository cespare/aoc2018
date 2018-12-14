package main

import (
	"bufio"
	"fmt"
	"log"
	"sort"
	"strings"
)

func init() {
	addSolutions(13, (*problemContext).problem13)
}

func (ctx *problemContext) problem13() {
	scanner := bufio.NewScanner(ctx.f)
	var grid [][]byte
	for scanner.Scan() {
		line := []byte(scanner.Text())
		if len(grid) > 0 && len(line) != len(grid[0]) {
			log.Fatalln("not all lines are the same length")
		}
		grid = append(grid, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tracks := newCartTracks(grid)
	foundCrash := false
	for {
		//fmt.Println(tracks)
		v, crashed := tracks.step()
		if crashed && !foundCrash {
			foundCrash = true
			ctx.reportPart1(fmt.Sprintf("%d,%d", v.x, v.y))
			return
		}
		if len(tracks.carts) == 1 {
			var v vec2
			for loc := range tracks.carts {
				v = loc
				break
			}
			ctx.reportPart2(fmt.Sprintf("%d,%d", v.x, v.y))
			return
		}
		if len(tracks.carts) == 0 {
			panic("empty")
		}
	}
}

type cartTracks struct {
	grid  [][]byte
	w, h  int
	carts map[vec2]cart
}

type cart struct {
	loc   vec2
	turns int
	c     byte
}

func newCartTracks(grid [][]byte) *cartTracks {
	t := &cartTracks{
		grid:  grid,
		w:     len(grid[0]),
		h:     len(grid),
		carts: make(map[vec2]cart),
	}
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			c := t.grid[y][x]
			switch c {
			case '>', '<':
				t.grid[y][x] = '-'
			case '^', 'v':
				t.grid[y][x] = '|'
			default:
				continue
			}
			v := vec2{int64(x), int64(y)}
			t.carts[v] = cart{
				loc: v,
				c:   c,
			}
		}
	}
	return t
}

func (t *cartTracks) at(v vec2) byte {
	return t.grid[v.y][v.x]
}

func (t *cartTracks) step() (crash vec2, crashed bool) {
	var carts []cart
	for _, c := range t.carts {
		carts = append(carts, c)
	}
	sort.Slice(carts, func(i, j int) bool {
		v0, v1 := carts[i].loc, carts[j].loc
		if v0.y < v1.y {
			return true
		}
		if v0.y > v1.y {
			return false
		}
		return v0.x < v1.x
	})

	var b strings.Builder
	for _, c := range carts {
		fmt.Fprintf(&b, "(%d,%d) ", c.loc.x, c.loc.y)
	}
	fmt.Println(b.String())

	var result vec2
	for _, c := range carts {
		if _, ok := t.carts[c.loc]; !ok {
			// Already deleted.
			continue
		}
		c1 := t.move(c)
		delete(t.carts, c.loc)
		if _, ok := t.carts[c1.loc]; ok {
			if !crashed {
				result = c1.loc
				crashed = true
			}
			delete(t.carts, c1.loc)
		} else {
			t.carts[c1.loc] = c1
		}
	}
	return result, crashed
}

func (t *cartTracks) String() string {
	grid := make([][]byte, len(t.grid))
	for i, line := range t.grid {
		grid[i] = append([]byte(nil), line...)
	}
	for loc, c := range t.carts {
		grid[loc.y][loc.x] = c.c
	}
	var b strings.Builder
	for _, line := range grid {
		b.Write(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func (t *cartTracks) move(c cart) cart {
	c1 := c
	switch c.c {
	case '<':
		c1.loc.x--
	case '^':
		c1.loc.y--
	case '>':
		c1.loc.x++
	case 'v':
		c1.loc.y++
	}

	switch g := t.at(c1.loc); g {
	case '+':
		switch c.turns % 3 {
		case 0:
			switch c.c {
			case '<':
				c1.c = 'v'
			case '^':
				c1.c = '<'
			case '>':
				c1.c = '^'
			case 'v':
				c1.c = '>'
			}
		case 2:
			switch c.c {
			case '<':
				c1.c = '^'
			case '^':
				c1.c = '>'
			case '>':
				c1.c = 'v'
			case 'v':
				c1.c = '<'
			}
		}
		c1.turns++
	case '\\':
		switch c.c {
		case '<':
			c1.c = '^'
		case '^':
			c1.c = '<'
		case '>':
			c1.c = 'v'
		case 'v':
			c1.c = '>'
		}
	case '/':
		switch c.c {
		case '<':
			c1.c = 'v'
		case '^':
			c1.c = '>'
		case '>':
			c1.c = '^'
		case 'v':
			c1.c = '<'
		}
	case '|', '-':
	default:
		panic(fmt.Sprintf("unexpected character: %c", g))
	}
	return c1
}
