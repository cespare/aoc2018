package main

import (
	"bufio"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(25, (*problemContext).problem25)
}

func (ctx *problemContext) problem25() {
	var points []vec4
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		points = append(points, parseVec4(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	g := newPointGraph()
	for i, v0 := range points {
		g.addVertex(i)
		for j := i + 1; j < len(points); j++ {
			v1 := points[j]
			if v0.dist(v1) <= 3 {
				g.addEdge(i, j)
			}
		}
	}

	components := g.connectedComponents()
	ctx.reportPart1(len(components))
}

type vec4 [4]int64

func (v vec4) dist(v1 vec4) int64 {
	var d int64
	for i, vv := range v {
		d += abs(vv - v1[i])
	}
	return d
}

func parseVec4(s string) vec4 {
	parts := strings.SplitN(s, ",", 4)
	if len(parts) != 4 {
		panic("bad vec4")
	}
	var v vec4
	for i, part := range parts {
		n, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			panic(err)
		}
		v[i] = n
	}
	return v
}

type pointGraph struct {
	vertices  intSet
	neighbors map[int][]int
}

func newPointGraph() *pointGraph {
	return &pointGraph{
		vertices:  make(intSet),
		neighbors: make(map[int][]int),
	}
}

func (g *pointGraph) addVertex(i int) {
	g.vertices.add(i)
}

func (g *pointGraph) addEdge(i, j int) {
	g.neighbors[i] = append(g.neighbors[i], j)
	g.neighbors[j] = append(g.neighbors[j], i)
}

func (g pointGraph) connectedComponents() []intSet {
	rem := g.vertices.copy()
	var extractComponent func(int, intSet)
	extractComponent = func(v int, comp intSet) {
		comp.add(v)
		rem.remove(v)
		for _, n := range g.neighbors[v] {
			if !comp.contains(n) {
				extractComponent(n, comp)
			}
		}
	}
	var components []intSet
	for len(rem) > 0 {
		comp := make(intSet)
		var v int
		for vv := range rem {
			v = vv
			break
		}
		extractComponent(v, comp)
		components = append(components, comp)
	}
	return components
}
