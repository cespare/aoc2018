package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(24, (*problemContext).problem24)
}

func (ctx *problemContext) problem24() {
	battle := new(immuneBattle)
	scanner := bufio.NewScanner(ctx.f)
	wantScan(scanner, "Immune System:")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}
		g := parseImmuneGroup(text)
		battle.immune = append(battle.immune, g)
		g.id = len(battle.immune)
	}
	wantScan(scanner, "Infection:")
	for scanner.Scan() {
		g := parseImmuneGroup(scanner.Text())
		battle.infect = append(battle.infect, g)
		g.id = len(battle.infect)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	b1 := battle.copy()
	for b1.step() {
	}
	_, units := b1.result()
	ctx.reportPart1(units)

	for boost := int64(1); ; boost++ {
		b2 := battle.copy()
		b2.boost(boost)
		for b2.step() {
		}
		immuneWon, units := b2.result()
		if !immuneWon {
			continue
		}
		ctx.reportPart2(units)
		break
	}
}

type immuneBattle struct {
	immune []*immuneGroup
	infect []*immuneGroup
}

type immuneGroup struct {
	id         int
	units      int64
	hp         int64
	attack     int64
	attackType string
	weaknesses []string
	immunities []string
	initiative int64
}

func (b *immuneBattle) copy() *immuneBattle {
	b1 := &immuneBattle{
		immune: make([]*immuneGroup, len(b.immune)),
		infect: make([]*immuneGroup, len(b.infect)),
	}
	for i, g := range b.immune {
		b1.immune[i] = g.copy()
	}
	for i, g := range b.infect {
		b1.infect[i] = g.copy()
	}
	return b1
}

func (g *immuneGroup) copy() *immuneGroup {
	g1 := *g
	// Shallow copy of lists is ok; they don't change.
	return &g1
}

func (b *immuneBattle) boost(n int64) {
	for _, g := range b.immune {
		g.attack += n
	}
}

func (b *immuneBattle) String() string {
	var sb strings.Builder
	printArmy := func(army []*immuneGroup) {
		for _, g := range army {
			fmt.Fprintf(&sb, "Group %d contains %d units\n", g.id, g.units)
		}
	}
	fmt.Fprintln(&sb, "Immune System:")
	printArmy(b.immune)
	fmt.Fprintln(&sb, "Infection:")
	printArmy(b.infect)
	return sb.String()
}

func (b *immuneBattle) result() (immuneWon bool, units int64) {
	winners := b.infect
	if len(winners) == 0 {
		immuneWon = true
		winners = b.immune
	}
	var total int64
	for _, g := range winners {
		total += g.units
	}
	return immuneWon, total
}

func (b *immuneBattle) step() bool {
	targets := make(map[*immuneGroup]*immuneGroup)
	for _, kv := range selectBattleTargets(b.immune, b.infect) {
		targets[kv[0]] = kv[1]
	}
	for _, kv := range selectBattleTargets(b.infect, b.immune) {
		targets[kv[0]] = kv[1]
	}
	var allGroups []*immuneGroup
	allGroups = append(allGroups, b.immune...)
	allGroups = append(allGroups, b.infect...)
	sort.Slice(allGroups, func(i, j int) bool {
		init0, init1 := allGroups[i].initiative, allGroups[j].initiative
		if init0 == init1 {
			panic("equal initiative")
		}
		return init0 > init1
	})
	anyKilled := false
	for _, g := range allGroups {
		enemy, ok := targets[g]
		if !ok {
			continue
		}
		if enemy.immuneTo(g.attackType) {
			continue
		}
		dmg := g.effectivePower()
		if enemy.weakTo(g.attackType) {
			dmg *= 2
		}
		numKilled := dmg / enemy.hp
		if numKilled > enemy.units {
			enemy.units = 0
		} else {
			enemy.units -= numKilled
		}
		if numKilled > 0 {
			anyKilled = true
		}
	}
	if !anyKilled {
		// stalemate
		return false
	}

	cleanUp := func(army []*immuneGroup) []*immuneGroup {
		var i int
		for _, g := range army {
			if g.units > 0 {
				army[i] = g
				i++
			}
		}
		return army[:i]
	}
	b.immune = cleanUp(b.immune)
	b.infect = cleanUp(b.infect)
	return len(b.immune) > 0 && len(b.infect) > 0
}

func selectBattleTargets(allies, enemies []*immuneGroup) [][2]*immuneGroup {
	alliesOrdered := append([]*immuneGroup(nil), allies...)
	sort.Slice(alliesOrdered, func(i, j int) bool {
		g0, g1 := alliesOrdered[i], alliesOrdered[j]
		ep0, ep1 := g0.effectivePower(), g1.effectivePower()
		if ep0 != ep1 {
			return ep0 > ep1
		}
		return g0.initiative > g1.initiative
	})
	enemiesLeft := append([]*immuneGroup(nil), enemies...)
	var targets [][2]*immuneGroup
	for _, ally := range alliesOrdered {
		var target *immuneGroup
		index := -1
		for i, enemy := range enemiesLeft {
			if ally.preferTarget(target, enemy) {
				target = enemy
				index = i
			}
		}
		if target != nil {
			targets = append(targets, [2]*immuneGroup{ally, target})
			enemiesLeft = append(enemiesLeft[:index], enemiesLeft[index+1:]...)
			if len(enemiesLeft) == 0 {
				break
			}
		}
	}
	return targets
}

// preferTarget says whether g would prefer to swap from t0 to t1.
func (g *immuneGroup) preferTarget(t0, t1 *immuneGroup) bool {
	if t1.immuneTo(g.attackType) {
		return false
	}
	if t0 == nil {
		return true
	}
	t0Weak, t1Weak := t0.weakTo(g.attackType), t1.weakTo(g.attackType)
	if t0Weak != t1Weak {
		return t1Weak
	}
	ep0, ep1 := t0.effectivePower(), t1.effectivePower()
	if ep0 != ep1 {
		return ep1 > ep0
	}
	return t1.initiative > t0.initiative
}

func (g *immuneGroup) immuneTo(typ string) bool { return stringsContain(g.immunities, typ) }
func (g *immuneGroup) weakTo(typ string) bool   { return stringsContain(g.weaknesses, typ) }

func stringsContain(slice []string, s string) bool {
	for _, s1 := range slice {
		if s1 == s {
			return true
		}
	}
	return false
}

func (g *immuneGroup) effectivePower() int64 {
	return g.units * g.attack
}

var immuneRegexp = regexp.MustCompile(`^(?P<Units>\d+) units each with (?P<HP>\d+) hit points (?P<ImmuneWeak>\(.*\) )?with an attack that does (?P<Attack>\d+) (?P<AttackType>\w+) damage at initiative (?P<Initiative>\d+)$`)

func parseImmuneGroup(s string) *immuneGroup {
	var v struct {
		Units      int64
		HP         int64
		ImmuneWeak string
		Attack     int64
		AttackType string
		Initiative int64
	}
	hasty.MustParse([]byte(s), &v, immuneRegexp)
	g := &immuneGroup{
		units:      v.Units,
		hp:         v.HP,
		attack:     v.Attack,
		attackType: v.AttackType,
		initiative: v.Initiative,
	}
	if v.ImmuneWeak == "" {
		return g
	}
	v.ImmuneWeak = strings.TrimPrefix(v.ImmuneWeak, "(")
	v.ImmuneWeak = strings.TrimSuffix(v.ImmuneWeak, ") ")
	for _, part := range strings.SplitN(v.ImmuneWeak, "; ", 2) {
		if l, ok := trimPrefix(part, "weak to "); ok {
			g.weaknesses = strings.Split(l, ", ")
		} else if l, ok := trimPrefix(part, "immune to "); ok {
			g.immunities = strings.Split(l, ", ")
		} else {
			panic("couldn't parse weak/immune")
		}
	}
	return g
}

func trimPrefix(s, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}
	return s, false
}

func wantScan(s *bufio.Scanner, text string) {
	if !s.Scan() {
		panic("short read")
	}
	if s.Text() != text {
		panic(fmt.Sprintf("want %q; got %q", text, s.Text()))
	}
}
