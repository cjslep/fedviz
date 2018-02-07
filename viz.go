package main

import (
	"image"
	"image/color"
	"math"
)

type Imager interface {
	Bounds() image.Rectangle
	IsSet(x, y int) bool
	IsColor(x, y int, c color.Color) bool
}

type ImageBuilder interface {
	Imager
	Set(x, y int, c color.Color)
	Build() image.Image
	Clear(x, y int)
}

type ImageFactory interface {
	NewImage(size int) ImageBuilder
	Color(id int) color.Color
}

type Algorithm interface {
	Start(Imager, color.Color) (x, y int)
	Next(Imager, color.Color) (x, y int, ok bool)
}

type FedMapVisualizer struct {
	F       *FedMap
	I       ImageFactory
	A       Algorithm
	Scaling float64
}

type point struct {
	x int
	y int
}

func (f *FedMapVisualizer) Generate() (image.Image, map[string]color.Color) {
	nodes := f.F.Nodes()
	builder := f.I.NewImage(int(f.Scaling * 0.59 * math.Sqrt(float64(f.F.NEdges()))))
	legend := make(map[string]color.Color, len(nodes))
	for _, node := range nodes {
		c := f.I.Color(node)
		legend[f.F.Name(node)] = c
		n := int(f.Scaling * math.Sqrt(float64(len(f.F.Edges(node)))))
		if n == 0 {
			continue
		}
		retry := true
		for retry {
			retry = false
			var myPoints []point
			x, y := f.A.Start(builder, c)
			builder.Set(x, y, c)
			myPoints = append(myPoints, point{x, y})
			for i := 1; i < n; i++ {
				x, y, ok := f.A.Next(builder, c)
				if !ok {
					retry = true
				} else {
					builder.Set(x, y, c)
					myPoints = append(myPoints, point{x, y})
				}
			}
			if retry {
				for _, p := range myPoints {
					builder.Clear(p.x, p.y)
				}
			}
		}
	}
	return builder.Build(), legend
}
