package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
)

func init() {
	addSolutions(7, (*problemContext).problem7)
}

func (ctx *problemContext) problem7() {
	steps, orders, err := ctx.loadInput7()
	if err != nil {
		log.Fatal(err)
	}
	ctx.reportLoad()

	parents := makeParentSets(orders, len(steps))

	var result1 []int
	getNext := func() int {
		next := -1
		for step, p := range parents {
			if p == nil {
				continue
			}
			if len(p) == 0 && (next < 0 || step < next) {
				next = step
			}
		}
		if next < 0 {
			return -1
		}
		parents[next] = nil
		return next
	}
	for i := 0; i < len(steps); i++ {
		next := getNext()
		for _, p := range parents {
			if p != nil {
				p.remove(next)
			}
		}
		result1 = append(result1, next)
	}
	var b strings.Builder
	for _, step := range result1 {
		b.WriteString(steps[step])
	}
	ctx.reportPart1(b.String())

	parents = makeParentSets(orders, len(steps))
	tick := 0
	const (
		numWorkers = 5
		overhead   = 60
	)
	type workerStatus struct {
		step      int
		remaining int
	}
	workers := make([]*workerStatus, numWorkers)
	for i := range workers {
		workers[i] = &workerStatus{step: -1}
	}
	for ; ; tick++ {
		for _, ws := range workers {
			if ws.remaining == 0 && ws.step >= 0 {
				for _, p := range parents {
					if p != nil {
						p.remove(ws.step)
					}
				}
			}
		}
		allIdle := true
		for _, ws := range workers {
			if ws.remaining == 0 {
				next := getNext()
				if next < 0 {
					continue
				}
				ws.step = next
				ws.remaining = next + 1 + overhead
			}
			allIdle = false
			ws.remaining--
		}
		if allIdle {
			break
		}
	}
	ctx.reportPart2(tick)
}

func makeParentSets(orders []stepOrder, numSteps int) []intSet {
	parents := make([]intSet, numSteps)
	for i := range parents {
		parents[i] = make(intSet)
	}
	for _, order := range orders {
		parents[order.second].add(order.first)
	}
	return parents
}

func (ctx *problemContext) loadInput7() (steps []string, orders []stepOrder, err error) {
	stepIdx := make(map[string]int)
	addStep := func(s string) int {
		n, ok := stepIdx[s]
		if ok {
			return n
		}
		n = len(steps)
		steps = append(steps, s)
		stepIdx[s] = n
		return n
	}
	scanner := bufio.NewScanner(ctx.f)
	for scanner.Scan() {
		first, second, err := parseStepOrder(scanner.Text())
		if err != nil {
			return nil, nil, err
		}
		var order stepOrder
		order.first = addStep(first)
		order.second = addStep(second)
		orders = append(orders, order)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	sort.Strings(steps)
	reorder := make([]int, len(steps))
	for i, step := range steps {
		reorder[stepIdx[step]] = i
	}
	for i, o := range orders {
		orders[i].first = reorder[o.first]
		orders[i].second = reorder[o.second]
	}
	return steps, orders, nil
}

type stepOrder struct {
	first  int
	second int
}

var stepRegex = regexp.MustCompile(`^Step (\w+) must be finished before step (\w+) can begin.$`)

func parseStepOrder(s string) (first, second string, err error) {
	matches := stepRegex.FindAllStringSubmatch(s, 1)
	if len(matches) != 1 {
		return "", "", fmt.Errorf("bad input")
	}
	return matches[0][1], matches[0][2], nil
}
