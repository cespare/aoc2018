package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
)

func init() {
	addSolutions(15, (*problemContext).problem15)
}

func (ctx *problemContext) problem15() {
	scanner := bufio.NewScanner(ctx.f)
	var lines [][]byte
	for scanner.Scan() {
		lines = append(lines, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	c, err := newCave(lines, 3)
	if err != nil {
		log.Fatal(err)
	}
	for rounds := int64(0); ; rounds++ {
		if _, done := c.step(); done {
			ctx.reportPart1(rounds * c.sumHP())
			break
		}
		//fmt.Println(c)
		//time.Sleep(250 * time.Millisecond)
	}

	for elfAttack := 4; ; elfAttack++ {
		c, err = newCave(lines, elfAttack)
		if err != nil {
			log.Fatal(err)
		}
		for rounds := int64(0); ; rounds++ {
			elfDied, done := c.step()
			if elfDied {
				break
			}
			if done {
				ctx.reportPart2(rounds * c.sumHP())
				return
			}
		}
	}
}

type cave struct {
	board [][]byte
	w, h  int
	units map[ivec2]*caveUnit
}

func newCave(lines [][]byte, elfAttack int) (*cave, error) {
	c := &cave{
		units: make(map[ivec2]*caveUnit),
		w:     len(lines[0]),
		h:     len(lines),
	}
	for y, line := range lines {
		if len(line) != c.w {
			return nil, errors.New("bad board")
		}
		c.board = append(c.board, append([]byte(nil), line...))
		for x, sq := range line {
			unit := &caveUnit{
				attack: 3,
				hp:     200,
			}
			switch sq {
			case '#', '.':
				continue
			case 'E':
				unit.typ = elf
				unit.attack = elfAttack
			case 'G':
				unit.typ = goblin
			default:
				return nil, fmt.Errorf("unexpected square: %c", sq)
			}
			unit.loc = ivec2{x, y}
			c.units[unit.loc] = unit
		}
	}
	return c, nil
}

func (c *cave) at(loc ivec2) byte {
	return c.board[loc.y][loc.x]
}

func (c *cave) occupied(loc ivec2) bool {
	return c.at(loc) != '.'
}

func (c *cave) set(loc ivec2, t byte) {
	c.board[loc.y][loc.x] = t
}

func (c *cave) step() (elfDied, done bool) {
	for _, unit := range c.orderUnits() {
		if unit.hp <= 0 {
			// Killed by a different unit earlier this round.
			continue
		}
		e, d := c.stepUnit(unit)
		if e {
			elfDied = true
		}
		if d {
			return elfDied, true
		}
	}
	return elfDied, false
}

// orderUnits retrieves the units as a slice sorted in turn (reading) order.
func (c *cave) orderUnits() []*caveUnit {
	var units []*caveUnit
	for _, unit := range c.units {
		units = append(units, unit)
	}
	sort.Slice(units, func(i, j int) bool {
		return units[i].loc.readingLess(units[j].loc)
	})
	return units
}

func (c *cave) stepUnit(unit *caveUnit) (elfDied, done bool) {
	targets := c.listTargets(unit)
	if len(targets) == 0 {
		return false, true
	}

	// Check if we're already adjacent to a target.
	needMove := true
	inRange := make(ivec2Set)
	for _, target := range targets {
		if unit.loc.adjacent(target.loc) {
			needMove = false
			break
		}
		for _, loc := range target.loc.neighbors() {
			if !c.occupied(loc) {
				inRange.add(loc)
			}
		}
	}

	if needMove {
		if len(inRange) == 0 {
			// No open squares in range of a target. End turn.
			return false, false
		}
		newLoc, ok := c.move(unit.loc, inRange)
		if !ok {
			// Can't move; nothing to do.
			return false, false
		}
		// Update the board.
		c.set(unit.loc, '.')
		c.set(newLoc, unit.typ.byte())
		// Update the unit map.
		delete(c.units, unit.loc)
		c.units[newLoc] = unit
		unit.loc = newLoc
	}

	// Done moving (if necessary). Attack!
	enemy := c.chooseAttack(unit)
	if enemy == nil {
		// No enemies in range.
		return false, false
	}
	enemy.hp -= unit.attack
	if enemy.hp <= 0 {
		if enemy.typ == elf {
			elfDied = true
		}
		c.set(enemy.loc, '.')
		delete(c.units, enemy.loc)
	}
	return elfDied, false
}

func (c *cave) listTargets(unit *caveUnit) []*caveUnit {
	var targets []*caveUnit
	for _, other := range c.units {
		if unit.typ != other.typ {
			targets = append(targets, other)
		}
	}
	return targets
}

func (c *cave) move(loc ivec2, goal ivec2Set) (ivec2, bool) {
	// Do a BFS from loc to any of the goal squares,
	// keeping track of all shortest paths.

	// First, a quick sanity check.
	if _, ok := goal[loc]; ok {
		panic("search starts at goal")
	}

	type exploredLoc struct {
		d       int
		parents []ivec2
	}

	explored := map[ivec2]*exploredLoc{
		loc: &exploredLoc{d: 0},
	}
	frontier := make(ivec2Set)
	frontier.add(loc)
	found := make(ivec2Set)
	for distance := 1; len(frontier) > 0 && len(found) == 0; distance++ {
		nextFrontier := make(ivec2Set)
		for loc1 := range frontier {
			for _, loc2 := range loc1.neighbors() {
				if c.occupied(loc2) {
					continue
				}
				if frontier.contains(loc2) {
					continue
				}
				el, ok := explored[loc2]
				if !ok {
					el = &exploredLoc{d: distance}
					explored[loc2] = el
					nextFrontier.add(loc2)
				}
				if el.d < distance {
					// Previously explored (part of last
					// round's frontier).
					continue
				}
				el.parents = append(el.parents, loc1)
				if goal.contains(loc2) {
					found.add(loc2)
				}
			}
		}
		frontier = nextFrontier
	}

	if len(found) == 0 {
		return ivec2{}, false
	}

	// We found one or more shortest paths to our goals.
	// found is the set of ends those paths.
	// We pick the one that's first in reading order.
	dest := ivec2{-1, -1}
	for end := range found {
		if dest.x < 0 || end.readingLess(dest) {
			dest = end
		}
	}

	// Backtrack from dest through all shortest paths back to loc,
	// then pick the best first step down any of those paths.
	best := ivec2{-1, -1}
	var backtrack func(v ivec2)
	backtrack = func(v ivec2) {
		el, ok := explored[v]
		if !ok {
			panic("location to backtrack not in explored")
		}
		for _, parent := range el.parents {
			if parent == loc {
				if best.x < 0 || v.readingLess(best) {
					best = v
				}
			} else {
				backtrack(parent)
			}
		}
	}
	backtrack(dest)
	return best, true
}

func (c *cave) chooseAttack(unit *caveUnit) *caveUnit {
	var enemies []*caveUnit
	for _, loc := range unit.loc.neighbors() {
		if u1, ok := c.units[loc]; ok {
			if u1.typ != unit.typ {
				enemies = append(enemies, u1)
			}
		}
	}
	if len(enemies) == 0 {
		return nil
	}
	sort.Slice(enemies, func(i, j int) bool {
		e0, e1 := enemies[i], enemies[j]
		if e0.hp == e1.hp {
			return e0.loc.readingLess(e1.loc)
		}
		return e0.hp < e1.hp
	})
	return enemies[0]
}

func (c *cave) sumHP() int64 {
	var total int64
	for _, unit := range c.units {
		total += int64(unit.hp)
	}
	return total
}

func (c *cave) String() string {
	var b strings.Builder
	for y := 0; y < c.h; y++ {
		var units []*caveUnit
		for x := 0; x < c.w; x++ {
			loc := ivec2{x, y}
			b.WriteByte(c.at(loc))
			if unit, ok := c.units[loc]; ok {
				units = append(units, unit)
			}
		}
		b.WriteByte(' ')
		for i, unit := range units {
			fmt.Fprintf(&b, "%c(%d)", unit.typ.byte(), unit.hp)
			if i < len(units)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

type unitType int

const (
	elf unitType = iota
	goblin
)

func (t unitType) byte() byte {
	switch t {
	case elf:
		return 'E'
	case goblin:
		return 'G'
	default:
		panic("bad unitType")
	}
}

type caveUnit struct {
	typ    unitType
	loc    ivec2
	attack int
	hp     int
}

func (u *caveUnit) String() string {
	return fmt.Sprintf("%c(%d)@(%d,%d)", u.typ.byte(), u.hp, u.loc.x, u.loc.y)
}

type ivec2 struct {
	x, y int
}

func (v ivec2) adjacent(v1 ivec2) bool {
	switch {
	case v.x == v1.x:
		return v.y == v1.y-1 || v.y == v1.y+1
	case v.y == v1.y:
		return v.x == v1.x-1 || v.x == v1.x+1
	}
	return false
}

func (v ivec2) neighbors() []ivec2 {
	return []ivec2{
		{v.x, v.y - 1},
		{v.x, v.y + 1},
		{v.x - 1, v.y},
		{v.x + 1, v.y},
	}
}

// readingLess returns true if v comes before v1 in reading order.
func (v ivec2) readingLess(v1 ivec2) bool {
	if v.y == v1.y {
		return v.x < v1.x
	}
	return v.y < v1.y
}

type ivec2Set map[ivec2]struct{}

func (s ivec2Set) add(v ivec2) { s[v] = struct{}{} }

func (s ivec2Set) contains(v ivec2) bool {
	_, ok := s[v]
	return ok
}
