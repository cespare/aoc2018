package main

import (
	"bufio"
	"log"
	"unicode/utf8"
)

func init() {
	addSolutions(2, (*problemContext).problem2)
}

func (ctx *problemContext) problem2() {
	var ids []string
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	var twoTimes, threeTimes int64
	for _, id := range ids {
		if exactly(id, 2) {
			twoTimes++
		}
		if exactly(id, 3) {
			threeTimes++
		}
	}
	ctx.reportPart1(twoTimes * threeTimes)

	seen := make(map[string]struct{})
	for _, id := range ids {
		for i, r := range id {
			w := utf8.RuneLen(r)
			s := id[:i] + "?" + id[i+w:]
			if _, ok := seen[s]; ok {
				ctx.reportPart2(id[:i] + id[i+w:])
				return
			}
			seen[s] = struct{}{}
		}
	}
}

func exactly(s string, n int) bool {
	m := make(map[rune]int)
	for _, r := range s {
		m[r]++
	}
	for _, n1 := range m {
		if n1 == n {
			return true
		}
	}
	return false
}
