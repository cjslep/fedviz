package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

var _ ImageBuilder = &rgbBuilder{}

type rgbBuilder struct {
	img *image.RGBA
	set [][]*color.RGBA
}

func (r *rgbBuilder) Bounds() image.Rectangle {
	return r.img.Bounds()
}

func (r *rgbBuilder) IsSet(x, y int) bool {
	return r.set[y][x] != nil
}

func (r *rgbBuilder) IsColor(x, y int, c color.Color) bool {
	rgb := c.(color.RGBA)
	if r.set[y][x] == nil {
		return false
	}
	return *r.set[y][x] == rgb
}

func (r *rgbBuilder) Set(x, y int, c color.Color) {
	rgb := c.(color.RGBA)
	r.set[y][x] = &rgb
}

func (r *rgbBuilder) Clear(x, y int) {
	r.set[y][x] = nil
}

func (r *rgbBuilder) Build() image.Image {
	for y, xSet := range r.set {
		for x, rgb := range xSet {
			if rgb != nil {
				r.img.SetRGBA(x, y, *rgb)
			}
		}
	}
	return r.img
}

var _ ImageFactory = &RGBFactory{}

type RGBFactory struct {
	created []color.RGBA
}

func (r *RGBFactory) NewImage(size int) ImageBuilder {
	rect := image.Rect(0, 0, size, size)
	img := image.NewRGBA(rect)
	set := make([][]*color.RGBA, size)
	for i := 0; i < size; i++ {
		set[i] = make([]*color.RGBA, size)
	}
	return &rgbBuilder{
		img: img,
		set: set,
	}
}

func (r *RGBFactory) Color(id int) color.Color {
	rng := rand.New(rand.NewSource(int64(id)))
	var c color.RGBA
	dupe := true
	i := 0
	for dupe {
		i++
		c = color.RGBA{0, 0, 0, 255}
		d := rng.Intn((255 / 5) * 6)
		d *= 5
		other := rng.Intn((125 / 5))
		other *= 5
		if d < 256 {
			c.R = 255
			c.G = uint8(d)
			c.B = uint8(other)
		} else if d < 255*2+1 {
			c.R = uint8(d - 256)
			c.G = 255
			c.B = uint8(other)
		} else if d < 255*3+1 {
			c.R = uint8(other)
			c.G = 255
			c.B = uint8(d - (255*2 + 1))
		} else if d < 255*4+1 {
			c.R = uint8(other)
			c.G = uint8(d - (255*3 + 1))
			c.B = 255
		} else if d < 255*5+1 {
			c.B = 255
			c.G = uint8(other)
			c.R = uint8(d - (255*4 + 1))
		} else if d < 255*6+1 {
			c.B = uint8(d - (255*5 + 1))
			c.G = uint8(other)
			c.R = 255
		}
		// Make each a multiple of 5
		c.R = uint8(math.Floor(float64(c.R)/5) * 5)
		c.G = uint8(math.Floor(float64(c.G)/5) * 5)
		c.B = uint8(math.Floor(float64(c.B)/5) * 5)
		dupe = false
		for _, created := range r.created {
			if created == c {
				dupe = true
				break
			}
		}
	}
	r.created = append(r.created, c)
	return c
}
