package main

import (
	"flag"
	"os"
	"time"

	"github.com/Telenav/osrm-backend/integration/service/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer/ranker"
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer/s2indexer"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	startTime := time.Now()

	if flags.inputFile == "" || flags.outputFolder == "" {
		glog.Fatal("Empty string for inputFile or outputFolder, please check your input.\n")
	}

	indexer := s2indexer.NewS2Indexer().Build(flags.inputFile)
	if indexer == nil {
		glog.Fatalf("Failed to build indexer, stop %s\n", os.Args[0])
	}
	indexer.Dump(flags.outputFolder)

	connectivitymap.NewConnectivityMap(flags.maxRange).
		Build(indexer, indexer, ranker.CreateRanker(ranker.SimpleRanker, nil), 1).
		Dump(flags.outputFolder)

	glog.Infof("%s totally takes %f seconds for processing.", os.Args[0], time.Since(startTime).Seconds())
}
