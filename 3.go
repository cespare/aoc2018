package main

import (
	"bufio"
	"log"
	"sort"
	"strconv"
	"strings"
)

func init() {
	addSolutions(3, (*problemContext).problem3a, (*problemContext).problem3b)
}

func (ctx *problemContext) problem3a() {
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
	ctx.reportLoad()

	ctx.reportPart1(g.count2())
	ctx.reportPart2(g.findNoOverlaps())
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

func (ctx *problemContext) problem3b() {
	var claims []*claim
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		c, err := parseClaim(scanner.Text())
		if err != nil {
			log.Fatalf("Bad claim %q: %s", scanner.Text(), err)
		}
		claims = append(claims, c)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	//var xstart, xend, ystart, yend []valID
	//for _, c := range claims {
	//        xstart = append(xstart, valID{c.x, c.id})
	//        xend = append(xend, valID{c.x + c.w, c.id})
	//}
	//for _, c := range claims {
	//        ystart = append(ystart, valID{c.y, c.id})
	//        yend = append(yend, valID{c.y + c.h, c.id})
	//}

	//ctx.reportPart1(g.count2())
	//ctx.reportPart2(g.findNoOverlaps())
}

type valID struct {
	v  int
	id int
}

func findOverlaps(starts, ends []valID, numClaims int) []intSet {
	return nil
	//sortValIDs(starts)
	//sortValIDs(ends)
	//startSets := collapseSortedValIDs(starts, numClaims)
	//endSets := collapseSortedValIDs(ends, numClaims)
	//current := make(intSet)
	//result := make([]intSet, numClaims)
	//for len(startSets) > 0 || len(endSets) > 0 {
	//        switch {
	//        case len(endSets) == 0 || startSets[0].v < endSets[0].v:
	//                for id := range startSets[0].s {

	//                }
	//        case len(startSets) == 0 || endSets[0].v < startSets[0].v:
	//        default:
	//        }
	//}
}

func collapseSortedInts(s []valID, numIDs int) []valSet {
	var result []valSet
	set := make(intSet)
	for i := 0; i < len(s); i++ {
		if i == 0 || s[i].v == s[i-1].v {
			set.add(s[i].id)
		} else {
			result = append(result, valSet{s[i].v, set})
			set = make(intSet)
		}
	}
	return result
}

func sortValIDs(s []valID) {
	sort.Slice(s, func(i, j int) bool { return s[i].v < s[j].v })
}

type intSet map[int]struct{}

func (s intSet) add(v int)    { s[v] = struct{}{} }
func (s intSet) remove(v int) { delete(s, v) }
func (s intSet) copy() intSet {
	s1 := make(intSet, len(s))
	for v := range s {
		s1[v] = struct{}{}
	}
	return s1
}

type valSet struct {
	v int
	s intSet
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
