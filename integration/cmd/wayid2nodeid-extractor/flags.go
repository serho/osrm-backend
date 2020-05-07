package main

import (
	"flag"
	"fmt"

	"github.com/serho/osrm-backend/integration/util/mapsource"
)

var flags struct {
	input     string
	output    string
	mapSource string
}

const (
	pbfSourceUniDB = "unidb"
	pbfSourceOSM   = "osm"
)

func init() {
	flag.StringVar(&flags.input, "i", "", "Input pbf file.")
	flag.StringVar(&flags.output, "o", "", "Output csv file")
	flag.StringVar(&flags.mapSource, "mapsource", mapsource.OSM, fmt.Sprintf("pbf map data source, can be '%s' or '%s'.", mapsource.UniDB, mapsource.OSM))
}
