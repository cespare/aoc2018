package main

import (
	"bufio"
	"log"
	"sort"
	"strconv"
	"strings"
)

func (ctx *problemContext) problem3() {
	var ids []string
	var colSet, rowSet intervalSet
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		c, err := parseClaim(scanner.Text())
		if err != nil {
			log.Fatalf("Bad claim %q: %s", scanner.Text(), err)
		}
		idx := len(ids)
		ids = append(ids, c.id)
		colSet.addClaim(idx, c.x, c.w)
		rowSet.addClaim(idx, c.y, c.h)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	colSet.sort()
	rowSet.sort()
}

type intervalSet struct {
	starts pointSet
	ends   pointSet
}

func (iset *intervalSet) addClaim(idx, start, delta int) {
	iset.starts = append(iset.starts, [2]int{idx, start})
	iset.ends = append(iset.ends, [2]int{idx, start + delta})
}

func (iset *intervalSet) sort() {
	sort.Sort(iset.starts)
	sort.Sort(iset.ends)
}

type pointSet [][2]int

func (s pointSet) Len() int           { return len(s) }
func (s pointSet) Less(i, j int) bool { return s[i][1] < s[j][1] }
func (s pointSet) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

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
