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
	numPlayers, err := strconv.ParseInt(matches[0][1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	end, err := strconv.ParseInt(matches[0][2], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	g := newMarbleGame(numPlayers)
	for i := int64(0); i < end; i++ {
		g.step()
		//ctx.l.Println(g.String())
	}
	ctx.reportPart1(g.highScore())

	g = newMarbleGame(numPlayers)
	for i := int64(0); i < end*100; i++ {
		g.step()
	}
	ctx.reportPart2(g.highScore())
}

type marbleGame struct {
	curPlayer int
	points    []int64
	last      int64
	current   *marble
}

func newMarbleGame(numPlayers int64) *marbleGame {
	m := &marble{id: 0}
	m.left = m
	m.right = m
	return &marbleGame{
		points:  make([]int64, numPlayers),
		current: m,
	}
}

func (g *marbleGame) step() bool {
	g.last++
	id := g.last
	done := false
	if id%23 == 0 {
		g.points[g.curPlayer] += id
		for i := 0; i < 7; i++ {
			g.current = g.current.left
		}
		g.points[g.curPlayer] += g.current.id
		g.current = g.current.remove()
	} else {
		g.current = g.current.right.insertRight(id)
	}
	g.curPlayer = (g.curPlayer + 1) % len(g.points)
	return !done
}

func (g *marbleGame) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%3d:", g.last)
	start := g.current.id
	cur := g.current
	for {
		fmt.Fprintf(&b, " %d", cur.id)
		cur = cur.right
		if cur.id == start {
			break
		}
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
	id          int64
	left, right *marble
}

func (m *marble) insertRight(id int64) *marble {
	cur := &marble{
		id:    id,
		left:  m,
		right: m.right,
	}
	m.right = cur
	cur.right.left = cur
	return cur
}

func (m *marble) remove() *marble {
	m.left.right = m.right
	m.right.left = m.left
	cur := m.right
	m.left, m.right = nil, nil
	return cur
}
