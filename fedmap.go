package main

import (
	"sort"
)

type FedMap struct {
	name   map[int]string
	edges  map[int][]int
	nEdges int
}

func NewFedMap() *FedMap {
	return &FedMap{
		name:  make(map[int]string, 0),
		edges: make(map[int][]int, 0),
	}
}

func (f *FedMap) AddNode(id int, name string) {
	f.name[id] = name
}

func (f *FedMap) AddEdge(from, to int) {
	fe, ok := f.edges[from]
	if ok {
		fe = append(fe, to)
	} else {
		fe = []int{to}
	}
	f.edges[from] = fe
	te, ok := f.edges[to]
	if ok {
		te = append(te, from)
	} else {
		te = []int{from}
	}
	f.edges[to] = te
	f.nEdges++
}

func (f *FedMap) Nodes() []int {
	r := make([]int, 0, len(f.name))
	for k, _ := range f.name {
		r = append(r, k)
	}
	sort.Sort(sort.IntSlice(r))
	return r
}

func (f *FedMap) Name(i int) string {
	name, ok := f.name[i]
	if !ok {
		return ""
	}
	return name
}

func (f *FedMap) Edges(i int) []int {
	edges, ok := f.edges[i]
	if !ok {
		return nil
	}
	return edges
}

func (f *FedMap) NEdges() int {
	return f.nEdges
}
