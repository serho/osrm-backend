package main

import (
	"os"
	"log"
	"fmt"
	"bufio"
	"sync"
	"time"

	"github.com/Telenav/osrm-backend/traffic_updater/go/gen-go/proxy"
)

const taskNum = 2
var dumpFinishedWg sync.WaitGroup
var tasksWg sync.WaitGroup

func dumpSpeedTable4Customize(flows []*proxy.Flow, wayid2nodeids map[uint64][]int64, outputPath string) {
	startTime := time.Now()

	if len(flows) == 0 || len(wayid2nodeids) == 0 {
		fmt.Printf("dumpSpeedTable4Customize failed due to empty input.  len(flows) = %d, len(wayid2nodeids) = %d\n", 
		len(flows), len(wayid2nodeids))
		return
	}

	sink := make(chan string)
	startMatchTasks(flows, wayid2nodeids, sink)
	startDump(outputPath, sink)
	wait4AllTasksFinished(sink)

	endTime := time.Now()
	fmt.Printf("Processing time for dumpSpeedTable4Customize takes %f seconds\n", endTime.Sub(startTime).Seconds())
}

func startMatchTasks(flows []*proxy.Flow, wayid2nodeids map[uint64][]int64, sink chan<- string) {
	tasksWg.Add(taskNum)
	var begin, end int
	for i := 0; i < taskNum; i++ {
		begin = end
		if begin >= len(flows) {
			break
		}
		end = begin + len(flows) / taskNum
		if end > len(flows) {
			end = len(flows)
		}
		go task(flows, wayid2nodeids, begin, end, sink)
	}
}

func startDump(outputPath string, sink <-chan string) {
	dumpFinishedWg.Add(1)
	go write(outputPath, sink)
}

func wait4AllTasksFinished(sink chan string) {
	tasksWg.Wait()
	close(sink)
	dumpFinishedWg.Wait()
}

func task(flows []*proxy.Flow, wayid2nodeids map[uint64][]int64, begin int, end int, sink chan<- string) {
	for i := begin; i < end; i++ {
		flow := flows[i]

		if nodes, ok := wayid2nodeids[ (uint64)(flow.WayId)]; ok {
			for n := 0; (n + 1) < len(nodes); n++ {
				sink <- generateSingleRecord(nodes[n], nodes[n+1], flow.Speed)
			}
		}
	}
	tasksWg.Done() 
}

// format
// if speed >= 0, means traffic for forward, generate: from, to, speed
// if speed < 0, means traffic for backward, generate: to, from, abs(speed)
// To be confirm: When speed = 0, do we need to ban both directions?
func generateSingleRecord(from, to int64, speed float64) (string){
	if (speed >= 0) {
		return fmt.Sprintf("%d,%d,%d\n", from, to, (int)(speed))
	} else {
		return fmt.Sprintf("%d,%d,%d\n", to, from, -(int)(speed))
	}
}

func write(targetPath string, sink <-chan string) {
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

	dumpFinishedWg.Done()
}
