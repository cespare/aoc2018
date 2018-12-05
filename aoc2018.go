package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var funcs = []func(*problemContext){
	(*problemContext).problem1,
	(*problemContext).problem2,
	(*problemContext).problem3,
	(*problemContext).problem4,
	(*problemContext).problem5,
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal("Usage: aoc2018 <problem>")
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln("Bad problem number", os.Args[1])
	}
	if n < 1 || n > len(funcs) {
		log.Fatalln("Bad problem number:", n)
	}
	ctx, err := newProblemContext(n)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.close()
	funcs[n-1](ctx)
}

type problemContext struct {
	f         *os.File
	needClose bool
}

func newProblemContext(n int) (*problemContext, error) {
	var ctx problemContext
	if len(os.Args) > 2 && os.Args[2] == "-" {
		ctx.f = os.Stdin
	} else {
		name := fmt.Sprintf("%d.txt", n)
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		ctx.f = f
		ctx.needClose = true
	}
	return &ctx, nil
}

func (ctx *problemContext) close() {
	if ctx.needClose {
		ctx.f.Close()
	}
}
