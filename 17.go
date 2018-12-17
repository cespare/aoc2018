package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(17, (*problemContext).problem17)
}

func (ctx *problemContext) problem17() {
	clay := make(ivec2Set)
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if line[0] == 'x' {
			var v vertVein
			if err := hasty.Parse(line, &v, vertVeinRegexp); err != nil {
				log.Fatal(err)
			}
			v.addToSet(clay)
		}
		if line[0] == 'y' {
			var v horizVein
			if err := hasty.Parse(line, &v, horizVeinRegexp); err != nil {
				log.Fatal(err)
			}
			v.addToSet(clay)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	scan := newEarthScan(clay)
	scan.fill()
	ctx.reportPart1(scan.countWater())
	ctx.reportPart2(scan.countStillWater())
}

var (
	vertVeinRegexp  = regexp.MustCompile(`^x=(?P<X>\d+), y=(?P<YBegin>\d+)\.\.(?P<YEnd>\d+)$`)
	horizVeinRegexp = regexp.MustCompile(`^y=(?P<Y>\d+), x=(?P<XBegin>\d+)\.\.(?P<XEnd>\d+)$`)
)

type vertVein struct {
	X      int
	YBegin int
	YEnd   int
}

func (v vertVein) addToSet(s ivec2Set) {
	for y := v.YBegin; y <= v.YEnd; y++ {
		s.add(ivec2{v.X, y})
	}
}

type horizVein struct {
	Y      int
	XBegin int
	XEnd   int
}

func (v horizVein) addToSet(s ivec2Set) {
	for x := v.XBegin; x <= v.XEnd; x++ {
		s.add(ivec2{x, v.Y})
	}
}

type earthScan struct {
	clay       ivec2Set
	water      map[ivec2]byte
	xmin, xmax int
	ymin, ymax int

	debug bool
}

func newEarthScan(clay ivec2Set) *earthScan {
	s := &earthScan{
		clay:  make(map[ivec2]struct{}),
		water: make(map[ivec2]byte),
		xmin:  99999999,
		xmax:  -99999999,
		ymin:  99999999,
		ymax:  -99999999,
	}
	for v := range clay {
		s.clay.add(v)
		if v.x < s.xmin {
			s.xmin = v.x
		}
		if v.x > s.xmax {
			s.xmax = v.x
		}
		if v.y < s.ymin {
			s.ymin = v.y
		}
		if v.y > s.ymax {
			s.ymax = v.y
		}
	}
	s.xmin--
	s.xmax += 2
	s.ymax++
	return s
}

func (s *earthScan) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "x=[%d, %d); y=[%d, %d)\n", s.xmin, s.xmax, s.ymin, s.ymax)
	for y := s.ymin; y < s.ymax; y++ {
		for x := s.xmin; x < s.xmax; x++ {
			v := ivec2{x, y}
			if c, ok := s.water[v]; ok {
				if c == '~' {
					fmt.Fprint(&b, "\x1b[38;5;21m█\x1b[m")
				} else {
					fmt.Fprint(&b, "\x1b[38;5;159m█\x1b[m")
				}
			} else if s.clay.contains(v) {
				fmt.Fprint(&b, "\x1b[38;5;130m█\x1b[m")
			} else {
				b.WriteByte(' ')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (s *earthScan) fill() {
	v := ivec2{500, 0}
	for v.y < s.ymin {
		v.y++
	}
	s.fillDown(v)
}

func (s *earthScan) fillDown(v ivec2) (contained bool) {
	if s.debug {
		fmt.Printf("fillDown(%v)\n%s\n", v, s)
		time.Sleep(100 * time.Millisecond)
	}

	orig := v
	for !s.clay.contains(v) {
		if s.water[v] == '|' {
			return false
		}
		s.water[v] = '|'
		v.y++
		if v.y >= s.ymax {
			return false
		}
	}
	for ; v != orig; v.y-- {
		containedLeft := s.fillLeftRight(v, ivec2{-1, 0})
		containedRight := s.fillLeftRight(v, ivec2{1, 0})
		if !containedLeft || !containedRight {
			return false
		}
		s.makeStill(v)
	}
	return true
}

func (s *earthScan) fillLeftRight(v, d ivec2) (contained bool) {
	if s.debug {
		fmt.Printf("fillLeftRight(%v, %v)\n%s\n", v, d, s)
		time.Sleep(100 * time.Millisecond)
	}

	for !s.clay.contains(v) {
		below := v.below()
		if !s.clay.contains(below) && s.water[below] != '~' {
			if !s.fillDown(v) {
				return false
			}
		}
		s.water[v] = '|'
		v = v.add(d)
	}
	return true
}

func (s *earthScan) makeStill(v ivec2) {
	for vl := v; !s.clay.contains(vl); vl = vl.left() {
		s.water[vl] = '~'
	}
	for vr := v; !s.clay.contains(vr); vr = vr.right() {
		s.water[vr] = '~'
	}
}

func (s *earthScan) countWater() int {
	return len(s.water)
}

func (s *earthScan) countStillWater() int {
	var total int
	for _, c := range s.water {
		if c == '~' {
			total++
		}
	}
	return total
}

func (v ivec2) add(d ivec2) ivec2 { return ivec2{v.x + d.x, v.y + d.y} }
func (v ivec2) left() ivec2       { return ivec2{v.x - 1, v.y} }
func (v ivec2) right() ivec2      { return ivec2{v.x + 1, v.y} }
func (v ivec2) below() ivec2      { return ivec2{v.x, v.y + 1} }
