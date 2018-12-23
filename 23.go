package main

import (
	"bufio"
	"log"
	"regexp"
	"sort"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(23, (*problemContext).problem23)
}

func (ctx *problemContext) problem23() {
	var nanobots []nanobot
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		nanobots = append(nanobots, parseNanobot(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var max nanobot
	for _, bot := range nanobots {
		if bot.r > max.r {
			max = bot
		}
	}
	var part1 int64
	for _, bot := range nanobots {
		if max.pos.dist(bot.pos) <= max.r {
			part1++
		}
	}
	ctx.reportPart1(part1)

	neighbors := make(map[int][]int)
	for i, bot0 := range nanobots {
		for j, bot1 := range nanobots {
			if i == j || bot0.pos.dist(bot1.pos) > bot0.r+bot1.r {
				continue
			}
			neighbors[j] = append(neighbors[j], i)
		}
	}
	cliques := bronKerbosch(neighbors)
	if len(cliques) != 1 {
		panic("unhandled")
	}
	clique := cliques[0]

	var maxMin int64
	var origin vec3
	for _, i := range clique {
		bot := nanobots[i]
		d := origin.dist(bot.pos) - bot.r
		if d > maxMin {
			maxMin = d
		}
	}
	ctx.reportPart2(maxMin)
}

func bronKerbosch(g map[int][]int) [][]int {
	p := make(intSet)
	for v := range g {
		p.add(v)
	}
	state := &bkState{g: g}
	bronKerbosch1(state, make(intSet), p, make(intSet))
	return state.maxCliques
}

type bkState struct {
	g          map[int][]int
	maxCliques [][]int
}

func bronKerbosch1(state *bkState, r, p, x intSet) {
	if len(p) == 0 && len(x) == 0 {
		if len(state.maxCliques) > 0 && len(state.maxCliques[0]) > len(r) {
			return
		}
		if len(state.maxCliques) > 0 && len(state.maxCliques[0]) < len(r) {
			// Found a longer clique.
			state.maxCliques = nil
		}
		clique := make([]int, 0, len(r))
		for v := range r {
			clique = append(clique, v)
		}
		sort.Ints(clique)
		state.maxCliques = append(state.maxCliques, clique)
		return
	}
	u := -1
	if len(p) > 0 {
		for v := range p {
			u = v
			break
		}
	} else {
		for v := range x {
			u = v
			break
		}
	}
	nu := state.g[u]
	nuSet := make(intSet, len(nu))
	for _, uu := range nu {
		nuSet.add(uu)
	}
	for v := range p {
		if nuSet.contains(v) {
			continue
		}
		ns := state.g[v]
		p1 := make(intSet, len(ns))
		x1 := make(intSet, len(ns))
		for _, n := range ns {
			if p.contains(n) {
				p1.add(n)
			}
			if x.contains(n) {
				x1.add(n)
			}
		}
		r.add(v)
		bronKerbosch1(state, r, p1, x1)
		r.remove(v)
		p.remove(v)
		x.add(v)
	}
}

type vec3 struct {
	x, y, z int64
}

func (v vec3) dist(v1 vec3) int64 {
	return abs(v1.x-v.x) + abs(v1.y-v.y) + abs(v1.z-v.z)
}

type nanobot struct {
	pos vec3
	r   int64
}

var nanobotRegexp = regexp.MustCompile(`^pos=<(?P<X>[-\d]+),(?P<Y>[-\d]+),(?P<Z>[-\d]+)>, r=(?P<R>\d+)$`)

func parseNanobot(s string) nanobot {
	var n struct {
		X, Y, Z int64
		R       int64
	}
	hasty.MustParse([]byte(s), &n, nanobotRegexp)
	return nanobot{
		pos: vec3{n.X, n.Y, n.Z},
		r:   n.R,
	}
}

func (s intSet) contains(v int) bool {
	_, ok := s[v]
	return ok
}
