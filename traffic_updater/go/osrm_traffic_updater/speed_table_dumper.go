package main

import (
	"os"
	"log"
	"fmt"
	"bufio"

	"github.com/Telenav/osrm-backend/traffic_updater/go/gen-go/proxy"
)

const taskNum = 10

func dumpSpeedTable4OSRM(flows []*proxy.Flow, wayid2nodeids map[uint64][]int64, output string) {
	var begin, end int
	sink := make(chan string)

	for i := 0; i < taskNum; i++ {
		begin = end
		if end >= len(flows) {
			break
		}

		end = begin + len(flows) / taskNum
		if end > len(flows) {
			end = len(flows)
		}
		go task(flows, wayid2nodeids, begin, end, sink)
	}
}

func task(flows []*proxy.Flow, wayid2nodeids map[uint64][]int64, begin int, end int, sink chan<- string) {
	for i := begin; i < end; i++ {
		flow := flows[i]

		if nodes, ok := wayid2nodeids[ (uint64)(flow.WayId)]; ok {
			var s string
			for n := 0; (n + 1) < len(nodes); n++ {
				if flow.Speed >= 0 {
					s = fmt.Sprintf("%d,%d,%d\n", nodes[n], nodes[n+1], flow.Speed)
				} else {
					s = fmt.Sprintf("%d,%d,%d\n", nodes[n+1], nodes[n], -flow.Speed)
				}
				sink <- s
			}

		}

		if i == (len(flows) -1) {
			close(sink)
		}
	} 
}

func write(targetPath string, sink chan string) {
	outfile, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0755)
	defer outfile.Close()
	defer outfile.Sync()
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Open output file of %s failed.\n", targetPath)
		return
	}
	fmt.Printf("Open output file of %s succeed.\n", targetPath)

	w := bufio.NewWriter(outfile)
	defer w.Flush()
	for str := range sink {
		_, err := w.WriteString(str)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}