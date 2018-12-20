package main

import (
	"bytes"
	"crypto/sha256"
	"encoding"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"strings"
)

func init() {
	addSolutions(20, (*problemContext).problem20)
}

func (ctx *problemContext) problem20() {
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	b = bytes.TrimSpace(b)
	re := parseRegex(b)

	for ctx.loopForProfile() {
		s := newExploreState(1000)
		re.explore(s)
		ctx.reportPart1(s.longest)
		ctx.reportPart2(len(s.valid))
	}
}

type regexNode interface {
	String() string
	explore(s *exploreState)
}

type exploreState struct {
	b       []byte
	hs      [][]byte
	h       hash.Hash
	valid   map[string]struct{} // set of SHA-256 sums
	longest int

	threshold int
}

func newExploreState(threshold int) *exploreState {
	return &exploreState{
		valid:     make(map[string]struct{}, 10000),
		h:         sha256.New(),
		threshold: threshold,
	}
}

func (s *exploreState) append(c byte) {
	n := len(s.b) - 1
	if len(s.b) > 0 && dirsPaired(s.b[n], c) {
		s.b = s.b[:n]
		s.hs = s.hs[:n]
		s.h = loadHash(s.hs[n-1])
		return
	}
	s.b = append(s.b, c)
	s.h.Write([]byte{c})
	s.hs = append(s.hs, storeHash(s.h))
	if len(s.b) > s.longest {
		s.longest = len(s.b)
	}
	if len(s.b) >= s.threshold {
		sum := s.h.Sum(nil)
		s.valid[string(sum[:])] = struct{}{}
	}
}

func storeHash(h hash.Hash) []byte {
	b, err := h.(encoding.BinaryMarshaler).MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

func loadHash(b []byte) hash.Hash {
	h := sha256.New()
	if err := h.(encoding.BinaryUnmarshaler).UnmarshalBinary(b); err != nil {
		panic(err)
	}
	return h
}

type regexConcat struct {
	nodes []regexNode
}

func (rc *regexConcat) explore(s *exploreState) {
	for _, n := range rc.nodes {
		n.explore(s)
	}
}

func (rc *regexConcat) String() string {
	var b strings.Builder
	for _, n := range rc.nodes {
		b.WriteString(n.String())
	}
	return b.String()
}

type regexAlternate struct {
	nodes []regexNode
}

func (ra *regexAlternate) String() string {
	var b strings.Builder
	b.WriteByte('(')
	for i, n := range ra.nodes {
		if i > 0 {
			b.WriteByte('|')
		}
		b.WriteString(n.String())
	}
	b.WriteByte(')')
	return b.String()
}

func (ra *regexAlternate) explore(s *exploreState) {
	b := s.b
	hs := s.hs
	var shortest []byte
	var shortestHashes [][]byte
	for i, n := range ra.nodes {
		s.b = b
		s.hs = hs
		n.explore(s)
		if i == 0 || len(s.b) < len(shortest) {
			shortest = s.b
			shortestHashes = s.hs
			if len(shortest) < len(b) {
				panic("unhandled")
			}
		}
	}
	s.b = shortest
	s.hs = shortestHashes
}

type regexLiteral string

func (rl regexLiteral) String() string { return string(rl) }

func (rl regexLiteral) explore(s *exploreState) {
	for i := 0; i < len(rl); i++ {
		s.append(rl[i])
	}
}

func parseRegex(b []byte) regexNode {
	if len(b) < 2 || b[0] != '^' || b[len(b)-1] != '$' {
		panic("bad regex")
	}
	b = b[1 : len(b)-1]
	re, rem := parseRegexConcat(b)
	if len(rem) > 0 {
		panic("trailing junk")
	}
	return re
}

func parseRegexConcat(b []byte) (regexNode, []byte) {
	var nodes []regexNode
	var n regexNode
concatLoop:
	for len(b) > 0 {
		switch b[0] {
		case 'E', 'N', 'W', 'S':
			n, b = parseRegexLiteral(b)
			nodes = append(nodes, n)
		case '(':
			n, b = parseRegexAlternate(b)
			nodes = append(nodes, n)
		default:
			break concatLoop
		}
	}
	switch len(nodes) {
	case 0:
		return regexLiteral(""), b
	case 1:
		return nodes[0], b
	default:
		return &regexConcat{nodes: nodes}, b
	}
}

func parseRegexLiteral(b []byte) (regexNode, []byte) {
	var i int
	for _, c := range b {
		switch c {
		case 'E', 'N', 'W', 'S':
			i++
			continue
		}
		break
	}
	return regexLiteral(b[:i]), b[i:]
}

func parseRegexAlternate(b []byte) (regexNode, []byte) {
	if b[0] != '(' {
		panic("can't happen")
	}
	b = b[1:]
	var ra regexAlternate
	var n regexNode
	for {
		n, b = parseRegexConcat(b)
		ra.nodes = append(ra.nodes, n)
		c := b[0]
		b = b[1:]
		switch c {
		case '|':
		case ')':
			return &ra, b
		default:
			panic(fmt.Sprintf("got %c", b[0]))
		}
	}
}

func dirsPaired(d0, d1 byte) bool {
	switch d0 {
	case 'E':
		return d1 == 'W'
	case 'N':
		return d1 == 'S'
	case 'W':
		return d1 == 'E'
	case 'S':
		return d1 == 'N'
	}
	panic("blah")
}
