package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

var funcs = []func(*problemContext){
	(*problemContext).problem1,
	(*problemContext).problem2,
	(*problemContext).problem3,
	(*problemContext).problem4,
	(*problemContext).problem5,
	(*problemContext).problem6,
}

func main() {
	log.SetFlags(0)

	cpuProfile := flag.String("cpuprofile", "", "write CPU profile to `file`")
	printTiming := flag.Bool("t", false, "print elapsed time")
	flag.Parse()

	if *printTiming && *cpuProfile != "" {
		log.Fatal("-t and -cpuprofile are incompatible")
	}
	if flag.NArg() < 1 {
		log.Fatal("Usage: aoc2018 [flags] <problem>")
	}
	n, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatalln("Bad problem number", flag.Arg(0))
	}
	if n < 1 || n > len(funcs) {
		log.Fatalln("Bad problem number:", n)
	}
	fn := funcs[n-1]
	ctx, err := newProblemContext(n)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.close()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln("Error writing CPU profile:", err)
			}
		}()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalln("Error starting CPU profile:", err)
		}
		defer pprof.StopCPUProfile()

		ctx.l = log.New(ioutil.Discard, "", 0)

		end := time.Now().Add(5 * time.Second)
		for time.Until(end) > 0 {
			fn(ctx)
		}
	}

	start := time.Now()
	fn(ctx)
	if *printTiming {
		ctx.l.Println("Elapsed:", time.Since(start))
	}
}

type problemContext struct {
	f         *os.File
	needClose bool
	l         *log.Logger
}

func newProblemContext(n int) (*problemContext, error) {
	ctx := &problemContext{
		l: log.New(os.Stdout, "", 0),
	}
	if flag.NArg() > 1 && flag.Arg(1) == "-" {
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
	return ctx, nil
}

func (ctx *problemContext) close() {
	if ctx.needClose {
		ctx.f.Close()
	}
}
