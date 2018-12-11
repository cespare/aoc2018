package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	addSolutions(9, (*problemContext).problem9)
}

var problem9Regexp = regexp.MustCompile(`^(\d+) players; last marble is worth (\d+) points`)

func (ctx *problemContext) problem9() {
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	matches := problem9Regexp.FindAllStringSubmatch(string(b), 1)
	if len(matches) != 1 {
		log.Fatal("Failed to parse input")
	}
	numPlayers, err := strconv.Atoi(matches[0][1])
	if err != nil {
		log.Fatal(err)
	}
	end, err := strconv.Atoi(matches[0][2])
	if err != nil {
		log.Fatal(err)
	}

	g := newMarbleGame(numPlayers)
	for i := 0; i < end; i++ {
		g.step()
	}
	ctx.reportPart1(g.highScore())

	for ctx.loopForProfile() {
		g = newMarbleGame(numPlayers)
		for i := 0; i < end*100; i++ {
			g.step()
		}
	}
	ctx.reportPart2(g.highScore())
}

type marbleGame struct {
	curPlayer int
	points    []int64
	last      int32
	current   marble
}

func newMarbleGame(numPlayers int) *marbleGame {
	chunk := &marbleChunk{n: 1}
	chunk.ids[0] = 0
	chunk.left = chunk
	chunk.right = chunk
	return &marbleGame{
		points:  make([]int64, numPlayers),
		current: marble{chunk, 0},
	}
}

func (g *marbleGame) step() bool {
	g.last++
	id := g.last
	done := false
	if id%23 == 0 {
		g.points[g.curPlayer] += int64(id)
		g.current = g.current.left(7)
		g.points[g.curPlayer] += int64(g.current.get())
		g.current = g.current.remove()
	} else {
		g.current = g.current.right().insertRight(id)
	}
	g.curPlayer++
	if g.curPlayer == len(g.points) {
		g.curPlayer = 0
	}
	return !done
}

func (g *marbleGame) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%3d:", g.last)
	chunk := g.current.chunk
	for i := 0; ; i++ {
		for j := 0; j < int(chunk.n); j++ {
			if i == 0 && j == g.current.i {
				fmt.Fprintf(&b, " [%d]", chunk.ids[j])
			} else {
				fmt.Fprintf(&b, " %d", chunk.ids[j])
			}
		}
		chunk = chunk.right
		if chunk == g.current.chunk {
			break
		}
		fmt.Fprint(&b, " ...")
	}
	return b.String()
}

func (g *marbleGame) highScore() int64 {
	var max int64
	for _, p := range g.points {
		if p > max {
			max = p
		}
	}
	return max
}

type marble struct {
	chunk *marbleChunk
	i     int
}

type marbleChunk struct {
	left, right *marbleChunk
	ids         [11]int32
	n           int32
}

func (m marble) get() int32 {
	return m.chunk.ids[m.i]
}

func (m marble) insertRight(id int32) marble {
	m1 := marble{m.chunk, m.i + 1}
	if m.chunk.tryInsert(m.i+1, id) {
		return m1
	}
	if m.i < int(len(m.chunk.ids))-1 {
		tmp := m.chunk.ids[m.chunk.n-1]
		copy(m.chunk.ids[m.i+2:], m.chunk.ids[m.i+1:])
		m.chunk.ids[m.i+1] = id
		id = tmp
	} else {
		m1 = marble{m.chunk.right, 0}
	}
	if m.chunk.right.tryInsert(0, id) {
		return m1
	}
	next := &marbleChunk{
		left:  m.chunk,
		right: m.chunk.right,
		n:     1,
	}
	if m1.i == 0 {
		m1.chunk = next
	}
	m.chunk.right.left = next
	m.chunk.right = next
	next.ids[0] = id
	return m1
}

func (c *marbleChunk) tryInsert(i int, id int32) bool {
	if int(c.n) >= len(c.ids) {
		return false
	}
	copy(c.ids[i+1:], c.ids[i:])
	c.ids[i] = id
	c.n++
	return true
}

func (m marble) remove() marble {
	copy(m.chunk.ids[m.i:], m.chunk.ids[m.i+1:])
	m.chunk.n--
	if m.chunk.n > 0 {
		if m.i == int(m.chunk.n) {
			return marble{m.chunk.right, 0}
		}
		return m
	}
	m.chunk.left.right = m.chunk.right
	m.chunk.right.left = m.chunk.left
	return marble{m.chunk.right, 0}
}

func (m marble) right() marble {
	if m.i < int(m.chunk.n)-1 {
		return marble{m.chunk, m.i + 1}
	}
	return marble{m.chunk.right, 0}
}

func (m marble) left(n int) marble {
	if m.i >= n {
		return marble{m.chunk, m.i - n}
	}
	chunk := m.chunk.left
	result := marble{chunk, int(chunk.n) - 1}.left(n - m.i - 1)
	return result
}
