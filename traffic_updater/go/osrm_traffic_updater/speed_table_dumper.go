package main

import (
	"os"
	"log"
	"fmt"
	"bufio"
	"sync"
	"time"
	"strings"
	"strconv"
)

var tasksWg sync.WaitGroup
var dumpFinishedWg sync.WaitGroup

func dumpSpeedTable4Customize(wayid2speed map[uint64]int, sources [TASKNUM]chan string, outputPath string) {
	startTime := time.Now()

	if len(wayid2speed) == 0 {
		return
	}

	sink := make(chan string)
	startMatchTasks(wayid2speed, sources, sink)
	startDump(outputPath, sink)
	wait4AllTasksFinished(sink)

	endTime := time.Now()
	fmt.Printf("Processing time for dumpSpeedTable4Customize takes %f seconds\n", endTime.Sub(startTime).Seconds())
}

func startMatchTasks(wayid2speed map[uint64]int, sources [TASKNUM]chan string, sink chan<- string) {
	tasksWg.Add(TASKNUM)
	for i := 0; i < TASKNUM; i++ {
		go task(wayid2speed, sources[i], sink)
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

func task(wayid2speed map[uint64]int, source <-chan string, sink chan<- string) {
	var err error
	for str := range source {
		elements := strings.Split(str, ",")
		if len(elements) < 3 {
			continue
		}

		var wayid uint64
		if wayid, err = strconv.ParseUint(elements[0], 10, 64); err != nil {
			fmt.Printf("#Error during decoding wayid, row = %v\n", elements)
			continue
		}

		if speed, ok:= wayid2speed[wayid]; ok {
			var nodes []string = elements[1:]
			for i := 0; (i + 1) < len(nodes); i++ {
				var n1, n2 uint64
				if n1, err = strconv.ParseUint(nodes[i], 10, 64); err != nil {
					fmt.Printf("#Error during decoding nodeid, row = %v\n", elements)
					continue
				}
				if n2, err = strconv.ParseUint(nodes[i+1], 10, 64); err != nil {
					fmt.Printf("#Error during decoding nodeid, row = %v\n", elements)
					continue
				}
				sink <- generateSingleRecord(n1, n2, speed)
			}
		}
	}
	tasksWg.Done() 
}

// format
// if speed >= 0, means traffic for forward, generate: from, to, speed
// if speed < 0, means traffic for backward, generate: to, from, abs(speed)
// To be confirm: When speed = 0, do we need to ban both directions?
func generateSingleRecord(from, to uint64, speed int) (string){
	if (speed >= 0) {
		return fmt.Sprintf("%d,%d,%d\n", from, to, speed)
	} else {
		return fmt.Sprintf("%d,%d,%d\n", to, from, -speed)
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
