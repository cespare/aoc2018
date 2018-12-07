package main

import (
	"bufio"
	"log"
	"strconv"
)

func init() {
	addSolutions(1, (*problemContext).problem1)
}

func (ctx *problemContext) problem1() {
	var deltas []int64
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		n, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatalln("Bad number:", scanner.Text())
		}
		deltas = append(deltas, n)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	var total int64
	for _, d := range deltas {
		total += d
	}
	ctx.reportPart1(total)

	var freq int64
	seen := make(map[int64]struct{})
	for i := 0; ; i = (i + 1) % len(deltas) {
		if _, ok := seen[freq]; ok {
			ctx.reportPart2(freq)
			return
		}
		seen[freq] = struct{}{}
		freq += deltas[i]
	}
}
