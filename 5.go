package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

func (ctx *problemContext) problem5() {
	p, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	p = bytes.TrimSpace(p)

	fmt.Println(reducedPolymerLen(p))

	seen := make(map[byte]struct{})
	min := len(p)
	for _, c := range p {
		c = toLower(c)
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		n := reducedPolymerLen(removePolymer(p, c))
		if n < min {
			min = n
		}
	}
	fmt.Println(min)
}

func removePolymer(p []byte, c byte) []byte {
	p0 := make([]byte, len(p))
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
	p0 := append([]byte(nil), p...)
	for {
		n := len(p0)
		p0 = reducePolymer1(p0)
		if len(p0) == n {
			return len(p0)
		}
	}
}

func reducePolymer1(p []byte) []byte {
	i := 0
	for j := 0; j < len(p); j++ {
		if j < len(p)-1 && letterPairs[p[j]] == p[j+1] {
			j++
		} else {
			p[i] = p[j]
			i++
		}
	}
	return p[:i]
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
