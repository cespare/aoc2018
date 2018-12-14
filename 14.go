package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func init() {
	addSolutions(14, (*problemContext).problem14)
}

func (ctx *problemContext) problem14() {
	b, err := ioutil.ReadAll(ctx.f)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.TrimSpace(string(b))
	searchDigits := stringDigits(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	r := newRecipeBoard()
	for {
		r.addDigits()
		r.moveElf(0)
		r.moveElf(1)
		if score, ok := r.getScore(n); ok {
			ctx.reportPart1(score)
			break
		}
	}

	for ctx.loopForProfile() {
		r = newRecipeBoard()
		for i := 0; ; i++ {
			r.addDigits()
			r.moveElf(0)
			r.moveElf(1)
			if i < 3 {
				continue
			}
			if result, ok := r.checkRecipes(searchDigits); ok {
				ctx.reportPart2(result)
				break
			}
		}
	}
}

func stringDigits(s string) []uint8 {
	digits := make([]uint8, len(s))
	for i := range digits {
		digits[i] = s[i] - '0'
	}
	return digits
}

type recipeBoard struct {
	digits []uint8
	elves  [2]int
}

func newRecipeBoard() *recipeBoard {
	return &recipeBoard{
		digits: []uint8{3, 7},
		elves:  [2]int{0, 1},
	}
}

func (r *recipeBoard) addDigits() {
	sum := r.digits[r.elves[0]] + r.digits[r.elves[1]]
	if sum >= 10 {
		r.digits = append(r.digits, 1, sum-10)
	} else {
		r.digits = append(r.digits, sum)
	}
}

func (r *recipeBoard) moveElf(i int) {
	j := r.elves[i]
	next := j + int(r.digits[j]) + 1
	for next >= len(r.digits) {
		next -= len(r.digits)
	}
	r.elves[i] = next
}

func (r *recipeBoard) getScore(n int) (int, bool) {
	if len(r.digits) < n+10 {
		return 0, false
	}
	var score int
	for i := n; i < n+10; i++ {
		score *= 10
		score += int(r.digits[i])
	}
	return score, true
}

func (r *recipeBoard) checkRecipes(digits []uint8) (int, bool) {
	start := len(r.digits) - len(digits) - 1
	if bytes.Equal(r.digits[start:start+len(digits)], digits) {
		return start, true
	}
	start++
	if bytes.Equal(r.digits[start:start+len(digits)], digits) {
		return start, true
	}
	return 0, false
}
