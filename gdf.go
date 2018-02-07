package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type GDFScanner struct {
	filename        string
	nodeFn          func(string)
	edgeFn          func(string)
	f               *os.File
	s               *bufio.Scanner
	scanning        bool
	foundNodeHeader bool
	foundEdgeHeader bool
}

func NewGDFScanner(filename string, nodeFn, edgeFn func(string)) *GDFScanner {
	return &GDFScanner{
		filename: filename,
		nodeFn:   nodeFn,
		edgeFn:   edgeFn,
	}
}

func (g *GDFScanner) IsOpen() bool {
	return g.s != nil
}

func (g *GDFScanner) Close() error {
	err := g.f.Close()
	g.s = nil
	g.f = nil
	return err
}

func (g *GDFScanner) Open() error {
	if g.IsOpen() {
		if err := g.Close(); err != nil {
			return err
		}
	}
	f, err := os.Open(g.filename)
	if err != nil {
		return err
	}
	g.f = f
	g.s = bufio.NewScanner(f)
	g.scanning = true
	return nil
}

func (g *GDFScanner) Scan() error {
	g.scanning = g.s.Scan()
	if !g.scanning {
		return g.s.Err()
	}
	if !g.foundNodeHeader {
		if strings.HasPrefix(g.s.Text(), "nodedef>") {
			g.foundNodeHeader = true
		} else {
			return fmt.Errorf("expected gdf nodedef header")
		}
	} else if !g.foundEdgeHeader {
		if strings.HasPrefix(g.s.Text(), "edgedef>") {
			g.foundEdgeHeader = true
		} else {
			g.nodeFn(g.s.Text())
		}
	} else {
		g.edgeFn(g.s.Text())
	}
	return g.s.Err()
}

func (g *GDFScanner) Done() bool {
	return !g.scanning
}
