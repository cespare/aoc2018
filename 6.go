package main

import (
	"bufio"
	"log"
	"strconv"
	"strings"
)

func (ctx *problemContext) problem6() {
	var points []point
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		p, err := parsePoint(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		points = append(points, p)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	xmin := points[0].x
	xmax := xmin
	ymin := points[0].y
	ymax := ymin

	for i := range points {
		if points[i].x < xmin {
			xmin = points[i].x
		}
		if points[i].x > xmax {
			xmax = points[i].x
		}
		if points[i].y < ymin {
			ymin = points[i].y
		}
		if points[i].y > ymax {
			ymax = points[i].y
		}
	}
	xmin -= 10
	xmax += 10
	ymin -= 10
	ymax += 10

	width := xmax - xmin
	height := ymax - ymin

	for i := range points {
		points[i].x -= xmin
		points[i].y -= ymin
	}

	grid := make([]int, width*height)
	for y := int64(0); y < height; y++ {
		for x := int64(0); x < width; x++ {
			min := width + height
			for i, p := range points {
				d := p.dist(point{x, y})
				if d == min {
					grid[y*width+x] = -1
				} else if d < min {
					min = d
					grid[y*width+x] = i
				}
			}
		}
	}

	areas := make([]int, len(points))
	for y := int64(0); y < height; y++ {
		for x := int64(0); x < width; x++ {
			i := grid[y*width+x]
			switch {
			case i == -1:
				continue
			case x == 0 || x == width-1 || y == 0 || y == height-1:
				areas[i] = -1
			case areas[i] >= 0:
				areas[i]++
			}
		}
	}
	max := 0
	for _, a := range areas {
		if a > max {
			max = a
		}
	}
	ctx.l.Println(max)

	for i := range grid {
		grid[i] = 0
	}
	n := 0
	const limit = 10e3
	for y := int64(0); y < height; y++ {
		for x := int64(0); x < width; x++ {
			var d int64
			for _, p := range points {
				d += p.dist(point{x, y})
				if d >= limit {
					break
				}
			}
			if d < limit {
				n++
			}
		}
	}
	ctx.l.Println(n)
}

type point struct {
	x, y int64
}

func (p point) dist(p0 point) int64 {
	dx := p0.x - p.x
	if dx < 0 {
		dx = -dx
	}
	dy := p0.y - p.y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func parsePoint(s string) (point, error) {
	parts := strings.Split(s, ", ")
	var p point
	var err error
	p.x, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return p, err
	}
	p.y, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return p, nil
	}
	return p, nil
}
