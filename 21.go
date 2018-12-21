package main

func init() {
	addSolutions(21, (*problemContext).problem21)
}

func (ctx *problemContext) problem21() {
	first, last := do21()
	ctx.reportPart1(first)
	ctx.reportPart2(last)
}

func do21() (first, last int64) {
	seen := make(map[int64]struct{})
	var r2, r3, r4, r5 int64

	r5 = 0
l06:
	r2 = r5 | 65536
	r5 = 7571367
l08:
	r4 = r2 & 255
	r5 += r4
	r5 &= 16777215
	r5 *= 65899
	r5 &= 16777215

	// r4 = 256 > r2 ? 1 : 0
	// jmp +r4
	// jmp 16
	// jmp 27
	if 256 > r2 {
		r4 = 1
		goto l28
	} else {
		r4 = 0
		goto l17
	}

l17:
	r4 = 0
l18:
	r3 = r4 + 1
	r3 *= 256

	// r3 = r3 > r2 ? 1 : 0
	// jmp +r3
	// jmp 23
	// jmp 25
	if r3 > r2 {
		r3 = 1
		goto l26
	} else {
		r3 = 0
		goto l24
	}

l24:
	r4++
	// jmp 17
	goto l18
l26:
	r2 = r4
	// jmp 07
	goto l08

l28:
	// r4 = r5 == r0 ? 1 : 0
	// jmp +r4
	// jmp 05
	if len(seen) == 0 {
		first = r5
	}
	if _, ok := seen[r5]; ok {
		return
	}
	last = r5
	seen[r5] = struct{}{}
	goto l06
}
