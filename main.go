package main

import (
	"flag"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	infile := flag.String("input", "map.gdf", "Input GDF map to read.")
	outfile := flag.String("output", "fediverse.png", "Output filename for the visualization.")
	scaling := flag.Float64("scale", 1, "Scaling factor of image")
	flag.Parse()

	fedMap := NewFedMap()
	nodeFn := func(s string) {
		line := strings.Split(s, ",")
		if len(line) != 2 {
			log.Printf("node line does not have two values: %q", s)
			return
		}
		i64, err := strconv.ParseInt(line[0], 10, 32)
		if err != nil {
			log.Printf("node line does not have integer id: %s", err)
			return
		}
		fedMap.AddNode(int(i64), line[1])
	}
	edgeFn := func(s string) {
		line := strings.Split(s, ",")
		if len(line) != 2 {
			log.Printf("edge line does not have two values: %q", s)
			return
		}
		f64, err := strconv.ParseInt(line[0], 10, 32)
		if err != nil {
			log.Printf("edge line does not have integer 'from' id: %s", err)
			return
		}
		t64, err := strconv.ParseInt(line[1], 10, 32)
		if err != nil {
			log.Printf("edge line does not have integer 'to' id: %s", err)
			return
		}
		fedMap.AddEdge(int(f64), int(t64))
	}
	gdfScanner := NewGDFScanner(*infile, nodeFn, edgeFn)
	if err := gdfScanner.Open(); err != nil {
		log.Fatalf("gdf scanner cannot open file: %s", err)
	}
	for !gdfScanner.Done() {
		gdfScanner.Scan()
	}
	gen := &FedMapVisualizer{
		F:       fedMap,
		I:       &RGBFactory{},
		A:       NewCenterRandomWalker(0, 100),
		Scaling: *scaling,
	}
	image, _ := gen.Generate()
	out, err := os.Create(*outfile)
	if err != nil {
		log.Fatalf("error creating %q: %s", *outfile, err)
	}
	if err := png.Encode(out, image); err != nil {
		log.Fatalf("error writing to %q: %s", *outfile, err)
	}
}
