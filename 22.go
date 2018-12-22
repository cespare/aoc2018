package main

import (
	"bufio"
	"container/heap"
	"log"
	"math/bits"
	"strconv"
	"strings"
)

func init() {
	addSolutions(22, (*problemContext).problem22)
}

func (ctx *problemContext) problem22() {
	scanner := bufio.NewScanner(ctx.f)
	if !scanner.Scan() {
		log.Fatal("blah")
	}
	depth, err := strconv.ParseInt(strings.TrimPrefix(scanner.Text(), "depth: "), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	if !scanner.Scan() {
		log.Fatal("blah")
	}
	target, err := parseVec2(strings.TrimPrefix(scanner.Text(), "target: "))
	if err != nil {
		log.Fatal(err)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	c := newCave2(depth, target)
	ctx.reportPart1(c.risk())
	ctx.reportPart2(c.fastestRescue())
}

type cave2 struct {
	depth  int64
	target vec2
	els    map[vec2]int64
}

const caveMod = 20183

func newCave2(depth int64, target vec2) *cave2 {
	c := &cave2{
		depth:  depth,
		target: target,
		els:    make(map[vec2]int64),
	}
	return c
}

func (c *cave2) erosionLevel(v vec2) (el int64) {
	var ok bool
	el, ok = c.els[v]
	if ok {
		return el
	}
	defer func() { c.els[v] = el }()
	switch {
	case v == (vec2{0, 0}) || v == c.target:
		return c.depth % caveMod
	case v.y == 0:
		return (v.x*16807 + c.depth) % caveMod
	case v.x == 0:
		return (v.y*48271 + c.depth) % caveMod
	}
	gi := c.erosionLevel(vec2{v.x - 1, v.y}) * c.erosionLevel(vec2{v.x, v.y - 1})
	return (gi + c.depth) % caveMod
}

func (c *cave2) regionType(v vec2) uint {
	return 1 << uint(c.erosionLevel(v)%3)
}

func (c *cave2) risk() int64 {
	var r int64
	for y := int64(0); y <= c.target.y; y++ {
		for v := (vec2{0, y}); v.x <= c.target.x; v.x++ {
			r += c.erosionLevel(v) % 3
		}
	}
	return r
}

const (
	rocky  int64 = 0
	wet    int64 = 1
	narrow int64 = 2

	neither  uint = 0
	torch    uint = 1
	climbing uint = 2
)

func (c *cave2) fastestRescue() int64 {
	// A* search
	start := caveState{tool: torch}
	seen := make(map[caveState]struct{})
	frontier := newAstarQueue()
	heap.Push(frontier, caveStateCost{start, start.loc.dist(c.target)})
	gscore := map[caveState]int64{start: 0}
	for frontier.Len() > 0 {
		cur := heap.Pop(frontier).(caveStateCost)
		if cur.caveState == (caveState{c.target, torch}) {
			return cur.cost
		}
		seen[cur.caveState] = struct{}{}

		for _, n := range c.neighbors(cur.caveState) {
			g := gscore[cur.caveState] + n.cost
			if _, ok := seen[n.caveState]; ok {
				if g < gscore[n.caveState] {
					gscore[n.caveState] = g
				} else {
					continue
				}
			}
			if frontier.contains(n.caveState) {
				if g >= gscore[n.caveState] {
					continue
				}
			} else {
				f := g + n.loc.dist(c.target)
				heap.Push(frontier, caveStateCost{n.caveState, f})
			}
			gscore[n.caveState] = g
		}
	}
	panic("no path")
}

type caveState struct {
	loc  vec2
	tool uint
}

func (c *cave2) neighbors(cs caveState) []caveStateCost {
	neighbors := make([]caveStateCost, 0, 5)
	toolBit := uint(1) << cs.tool
	otherTool := uint(bits.TrailingZeros(7 &^ (toolBit | c.regionType(cs.loc))))
	neighbors = append(neighbors, caveStateCost{caveState{cs.loc, otherTool}, 7})

	left := vec2{cs.loc.x - 1, cs.loc.y}
	if left.x >= 0 && c.regionType(left)&toolBit == 0 {
		neighbors = append(neighbors, caveStateCost{caveState{left, cs.tool}, 1})
	}
	right := vec2{cs.loc.x + 1, cs.loc.y}
	if c.regionType(right)&toolBit == 0 {
		neighbors = append(neighbors, caveStateCost{caveState{right, cs.tool}, 1})
	}
	above := vec2{cs.loc.x, cs.loc.y - 1}
	if above.y >= 0 && c.regionType(above)&toolBit == 0 {
		neighbors = append(neighbors, caveStateCost{caveState{above, cs.tool}, 1})
	}
	below := vec2{cs.loc.x, cs.loc.y + 1}
	if c.regionType(below)&toolBit == 0 {
		neighbors = append(neighbors, caveStateCost{caveState{below, cs.tool}, 1})
	}
	return neighbors
}

type caveStateCost struct {
	caveState
	cost int64
}

type astarQueue struct {
	s []caveStateCost
	m map[caveState]struct{}
}

func newAstarQueue() *astarQueue {
	return &astarQueue{
		m: make(map[caveState]struct{}),
	}
}

func (q *astarQueue) Len() int           { return len(q.s) }
func (q *astarQueue) Less(i, j int) bool { return q.s[i].cost < q.s[j].cost }
func (q *astarQueue) Swap(i, j int)      { q.s[i], q.s[j] = q.s[j], q.s[i] }

func (q *astarQueue) Push(x interface{}) {
	state := x.(caveStateCost)
	q.s = append(q.s, state)
	q.m[state.caveState] = struct{}{}
}

func (q *astarQueue) Pop() interface{} {
	state := q.s[len(q.s)-1]
	q.s = q.s[:len(q.s)-1]
	delete(q.m, state.caveState)
	return state
}

func (q *astarQueue) contains(cs caveState) bool {
	_, ok := q.m[cs]
	return ok
}

func (v vec2) dist(v1 vec2) int64 {
	return abs(v1.y-v.y) + abs(v1.x-v.x)
}
