package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(10, (*problemContext).problem10)
}

func (ctx *problemContext) problem10() {
	var s sky
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		pv, err := parsePositionVelocity(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, pv)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	lastArea := s.area()
	secs := 1
	for ; ; secs++ {
		s.step()
		area := s.area()
		if area > lastArea {
			s.unstep()
			secs--
			break
		}
		lastArea = area
	}
	ctx.reportPart1(s)
	ctx.reportPart2(secs)
}

type sky []positionVelocity

func (s sky) step() {
	for i, pv := range s {
		pv.p = pv.p.add(pv.v)
		s[i] = pv
	}
}

func (s sky) unstep() {
	for i, pv := range s {
		pv.p = pv.p.add(pv.v.neg())
		s[i] = pv
	}
}

func (s sky) bounds() (vec2, vec2) {
	var xmin, xmax, ymin, ymax int64
	for i, pv := range s {
		x, y := pv.p.x, pv.p.y
		if i == 0 || x < xmin {
			xmin = x
		}
		if i == 0 || x > xmax {
			xmax = x
		}
		if i == 0 || y < ymin {
			ymin = y
		}
		if i == 0 || y > ymax {
			ymax = y
		}
	}
	return vec2{xmin, ymin}, vec2{xmax, ymax}
}

func (s sky) area() int64 {
	min, max := s.bounds()
	return (max.x - min.x) * (max.y - min.y)
}

func (s sky) String() string {
	set := make(map[vec2]struct{})
	for _, pv := range s {
		set[pv.p] = struct{}{}
	}
	min, max := s.bounds()
	var b strings.Builder
	for y := min.y - 2; y <= max.y+2; y++ {
		for x := min.x - 2; x <= max.x+2; x++ {
			if _, ok := set[vec2{x, y}]; ok {
				fmt.Fprint(&b, "#")
			} else {
				fmt.Fprint(&b, " ")
			}
		}
		fmt.Fprintln(&b)
	}
	return b.String()
}

type positionVelocity struct {
	p vec2
	v vec2
}

func parsePositionVelocity(s string) (pv positionVelocity, err error) {
	parts := strings.SplitN(s, " velocity=", 2)
	if len(parts) != 2 {
		return pv, fmt.Errorf("bad position+velocity: %q", s)
	}
	pv.p, err = parseVec2(strings.TrimPrefix(parts[0], "position="))
	if err != nil {
		return pv, err
	}
	pv.v, err = parseVec2(parts[1])
	return pv, err
}

type vec2 struct {
	x, y int64
}

func (v vec2) add(v1 vec2) vec2 {
	return vec2{v.x + v1.x, v.y + v1.y}
}

func (v vec2) neg() vec2 {
	return vec2{-v.x, -v.y}
}

func parseVec2(s string) (vec2, error) {
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")
	pair := strings.SplitN(s, ",", 2)
	if len(pair) != 2 {
		return vec2{}, fmt.Errorf("bad vec: %q", s)
	}
	x, err := strconv.ParseInt(strings.TrimSpace(pair[0]), 10, 64)
	if err != nil {
		return vec2{}, err
	}
	y, err := strconv.ParseInt(strings.TrimSpace(pair[1]), 10, 64)
	if err != nil {
		return vec2{}, err
	}
	return vec2{x, y}, nil
}
