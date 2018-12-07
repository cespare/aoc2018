package main

import (
	"bytes"
	"io/ioutil"
	"log"
)

func init() {
	addSolutions(5, (*problemContext).problem5)
}

func (ctx *problemContext) problem5() {
	p, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()
	p = bytes.TrimSpace(p)

	p0 := append([]byte(nil), p...)
	ctx.reportPart1(reducedPolymerLen(p0))

	seen := make(map[byte]struct{})
	min := len(p)
	for _, c := range p {
		c = toLower(c)
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		n := reducedPolymerLen(removePolymer(p, c, p0))
		if n < min {
			min = n
		}
	}
	ctx.reportPart2(min)
}

func removePolymer(p []byte, c byte, p0 []byte) []byte {
	p0 = p0[:len(p)]
	i := 0
	for _, c1 := range p {
		if toLower(c1) != c {
			p0[i] = c1
			i++
		}
	}
	return p0[:i]
}

func reducedPolymerLen(p []byte) int {
	i, j := 0, 1
	for {
		for i >= 0 && j < len(p) && letterPairs[p[i]] == p[j] {
			i--
			j++
		}
		if j == len(p) {
			return i + 1
		}
		i++
		p[i] = p[j]
		j++
	}
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c - 'A' + 'a'
	}
	return c
}

var letterPairs [256]byte

func init() {
	for c := byte('a'); c <= 'z'; c++ {
		letterPairs[c] = c - 'a' + 'A'
	}
	for c := byte('A'); c <= 'Z'; c++ {
		letterPairs[c] = c - 'A' + 'a'
	}
}
