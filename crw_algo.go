package main

import (
	"image/color"
	"math/rand"
)

var _ Algorithm = &centerRandomWalker{}

type centerRandomWalker struct {
	r           *rand.Rand
	x0          int
	y0          int
	maxAttempts int
}

func NewCenterRandomWalker(seed int64, maxAttempts int) Algorithm {
	return &centerRandomWalker{
		r:           rand.New(rand.NewSource(seed)),
		maxAttempts: maxAttempts,
	}
}

func (c *centerRandomWalker) Start(i Imager, col color.Color) (int, int) {
	rect := i.Bounds()
	xS := (rect.Max.X - rect.Min.X) / 2
	yS := (rect.Max.Y - rect.Min.Y) / 2
	found := false
	for !found {
		dir := c.r.Intn(8)
		t, u, v := c.dirTransforms(dir)
		x := xS
		y := yS
		for x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y && i.IsSet(x, y) {
			next := c.r.Intn(3)
			var tf func(int, int) (int, int)
			if next == 0 {
				tf = t
			} else if next == 1 {
				tf = u
			} else {
				tf = v
			}
			x, y = tf(x, y)
		}
		if x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y {
			if !i.IsSet(x, y) {
				found = true
				c.x0 = x
				c.y0 = y
			}
		}
	}
	return c.x0, c.y0
}

func (c *centerRandomWalker) Next(i Imager, col color.Color) (int, int, bool) {
	rect := i.Bounds()
	found := false
	var xS int
	var yS int
	attempts := 0 // Holes may form that are not big enough
	for !found && attempts < c.maxAttempts {
		dir := c.r.Intn(8)
		t, u, v := c.dirTransforms(dir)
		x := c.x0
		y := c.y0
		for x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y && i.IsColor(x, y, col) {
			next := c.r.Intn(3)
			var tf func(int, int) (int, int)
			if next == 0 {
				tf = t
			} else if next == 1 {
				tf = u
			} else {
				tf = v
			}
			x, y = tf(x, y)
		}
		if x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y {
			if !i.IsSet(x, y) {
				found = true
				xS = x
				yS = y
			}
		}
		attempts++
	}
	if !found {
		return 0, 0, false
	}
	return xS, yS, true
}

func (c *centerRandomWalker) dirTransforms(n int) (t, u, v func(x, y int) (int, int)) {
	xp := func(x, y int) (int, int) {
		return x + 1, y
	}
	xn := func(x, y int) (int, int) {
		return x - 1, y
	}
	yp := func(x, y int) (int, int) {
		return x, y + 1
	}
	yn := func(x, y int) (int, int) {
		return x, y - 1
	}
	xpyp := func(x, y int) (int, int) {
		return x + 1, y + 1
	}
	xpyn := func(x, y int) (int, int) {
		return x + 1, y - 1
	}
	xnyn := func(x, y int) (int, int) {
		return x - 1, y - 1
	}
	xnyp := func(x, y int) (int, int) {
		return x - 1, y + 1
	}
	if n == 0 {
		return xp, xpyp, xpyn
	} else if n == 1 {
		return xpyp, yp, xp
	} else if n == 2 {
		return yp, xnyp, xpyp
	} else if n == 3 {
		return xnyp, xn, yp
	} else if n == 4 {
		return xn, xnyn, xnyp
	} else if n == 5 {
		return xnyn, yn, xn
	} else if n == 6 {
		return yn, xpyn, xnyn
	} else {
		return xpyn, xp, yn
	}
}
