package main

import (
	"flag"
	"fmt"
	"github.com/Telenav/osrm-backend/traffic_updater/go/gen-go/proxy"
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


func main() {
	flag.Parse()

	isFlowDoneChan := make(chan bool)
	var flows []*proxy.Flow
	go getTrafficFlow(flags.ip, flags.port, flows, isFlowDoneChan)

	isLoadTableDoneChan := make(chan bool)
	wayid2Nodes := make(map[uint64][]int64)
	generateSpeedTable(flags.mappingFile, wayid2Nodes, isLoadTableDoneChan)
	//generateSpeedTable(wayid2speed, flags.mappingFile, flags.csvFile)

	var isFlowDone, isLoadTableDone bool
	controlChan := make(chan string, 2)
	for {
		select {
			case f := <- isFlowDoneChan :
				if !f {
					fmt.Printf("[ERROR] Communication with traffic server failed.\n")
					break
				} else {
					controlChan <- "flowIsDone"
				}
			case t := <- isLoadTableDoneChan :
				if !t {
					fmt.Printf("[ERROR] Load way to node mapping table failed.\n")
					break
				} else {
					controlChan <- "TableIsDone"
				}
			case r := <- controlChan : 
				switch r {
				case "flowIsDone":
					isFlowDone = true
				case "TableIsDone":
					isLoadTableDone = true
				}
				if isFlowDone && isLoadTableDone {
					dumpSpeedTable4OSRM(flows, wayid2Nodes, flags.csvFile)
				}
		}
	}

}
