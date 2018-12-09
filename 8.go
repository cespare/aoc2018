package main

import (
	"bufio"
	"log"
	"strconv"

	"github.com/kr/pretty"
)

func init() {
	addSolutions(8, (*problemContext).problem8)
}

func (ctx *problemContext) problem8() {
	var input []int64
	scanner := bufio.NewScanner(ctx.f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		n, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		input = append(input, n)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tree, rem := parseNode(input)
	if len(rem) > 0 {
		panic("input too long")
	}
	pretty.Println(tree)
	ctx.reportPart1(tree.sumMetadata())
	ctx.reportPart2(tree.value())
}

type node struct {
	children []*node
	metadata []int64
	val      int64
}

func parseNode(s []int64) (*node, []int64) {
	numMetadata := s[1]
	n := &node{
		children: make([]*node, s[0]),
		val:      -1,
	}
	s = s[2:]
	for i := range n.children {
		n.children[i], s = parseNode(s)
	}
	n.metadata = s[:numMetadata]
	return n, s[numMetadata:]
}

func (n *node) sumMetadata() int64 {
	var total int64
	for _, c := range n.children {
		total += c.sumMetadata()
	}
	for _, m := range n.metadata {
		total += m
	}
	return total
}

func (n *node) value() int64 {
	if n.val >= 0 {
		return n.val
	}
	if len(n.children) == 0 {
		n.val = n.sumMetadata()
		return n.val
	}
	n.val = 0
	for _, i := range n.metadata {
		i--
		if i >= int64(len(n.children)) {
			continue
		}
		n.val += n.children[i].value()
	}
	return n.val
}
