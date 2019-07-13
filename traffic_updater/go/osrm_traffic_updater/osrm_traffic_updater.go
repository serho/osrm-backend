package main

import (
	"flag"
	"fmt"
)

var flags struct {
	port          int
	ip            string
	mappingFile   string
	csvFile       string
	highPrecision bool
}

func init() {
	flag.IntVar(&flags.port, "p", 6666, "traffic proxy listening port")
	flag.StringVar(&flags.ip, "c", "127.0.0.1", "traffic proxy ip address")
	flag.StringVar(&flags.mappingFile, "m", "wayid2nodeids.csv", "OSRM way id to node ids mapping table")
	flag.StringVar(&flags.csvFile, "f", "traffic.csv", "OSRM traffic csv file")
	flag.BoolVar(&flags.highPrecision, "d", false, "use high precision speeds, i.e. decimal")
}

const TASKNUM = 128
const CACHEDOBJECTS = 1000000

func main() {
	flag.Parse()

	isFlowDoneChan := make(chan bool, 1)
	wayid2speed := make(map[uint64]int)
	go getTrafficFlow(flags.ip, flags.port, wayid2speed, isFlowDoneChan)

	var sources [TASKNUM]chan string
	for i := range sources {
		sources[i] = make(chan string, CACHEDOBJECTS)
	}
	go loadWay2NodeidsTable(flags.mappingFile, sources)

	isFlowDone := wait4PreConditions(isFlowDoneChan)
	if isFlowDone {
		dumpSpeedTable4Customize(wayid2speed, sources, flags.csvFile)
	}
}

func wait4PreConditions(flowChan <-chan bool) (bool) {
	var isFlowDone bool
	loop:
	for {
		select {
			case f := <- flowChan :
				if !f {
					fmt.Printf("[ERROR] Communication with traffic server failed.\n")
					break loop
				} else {
					isFlowDone = true
					break loop
				}
		}
	}
	return isFlowDone
}
