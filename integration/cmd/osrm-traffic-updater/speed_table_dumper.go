package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tasksWg sync.WaitGroup
var dumpFinishedWg sync.WaitGroup

func dumpSpeedTable4Customize(wayid2speed map[int64][3]int, sources [TASKNUM]chan string,
	outputPath string, ds *dumperStatistic) {
	startTime := time.Now()

	if len(wayid2speed) == 0 {
		return
	}

	sink := make(chan string)
	startTasks(wayid2speed, sources, sink, ds)
	startDump(outputPath, sink)
	wait4AllTasksFinished(sink, ds)

	endTime := time.Now()
	fmt.Printf("Processing time for dumpSpeedTable4Customize takes %f seconds\n", endTime.Sub(startTime).Seconds())
}

func startTasks(wayid2speed map[int64][3]int, sources [TASKNUM]chan string,
	sink chan<- string, ds *dumperStatistic) {
	tasksWg.Add(TASKNUM)
	for i := 0; i < TASKNUM; i++ {
		go task(wayid2speed, sources[i], sink, ds)
	}
}

func startDump(outputPath string, sink <-chan string) {
	dumpFinishedWg.Add(1)
	go write(outputPath, sink)
}

func wait4AllTasksFinished(sink chan string, ds *dumperStatistic) {
	tasksWg.Wait()
	close(sink)
	ds.Close()
	dumpFinishedWg.Wait()
}

func task(wayid2speed map[int64][3]int, source <-chan string, sink chan<- string, ds *dumperStatistic) {
	var wayCnt, nodeCnt, fwdRecordCnt, bwdRecordCnt, wayMatched, nodeMatched, fwdTrafficMatched, bwdTrafficMatched uint64
	var err error
	for str := range source {
		elements := strings.Split(str, ",")
		wayCnt += 1
		nodeCnt += (uint64)(len(elements) - 1)
		if len(elements) < 3 {
			continue
		}

		var wayid uint64
		if wayid, err = strconv.ParseUint(elements[0], 10, 64); err != nil {
			fmt.Printf("#Error during decoding wayid, row = %v\n", elements)
			continue
		}

		paramsFwd, okFwd := wayid2speed[(int64)(wayid)]
		paramsBwd, okBwd := wayid2speed[(int64)(-wayid)]

		speedFwd := paramsFwd[0]
		speedBwd := paramsBwd[0]

		offset := 1
		limit := len(elements)

		if paramsFwd[1] > 0 {
			offset = paramsFwd[1]
		}

		if paramsFwd[2] > 0 {
			limit = offset + paramsFwd[2]
		}

		if okFwd || okBwd {
			var nodes []string = elements[offset:limit]

			wayMatched += 1
			nodeMatched += (uint64)(len(nodes))
			if okFwd {
				fwdTrafficMatched += 1
			}
			if okBwd {
				bwdTrafficMatched += 1
			}

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
				if okFwd {
					fwdRecordCnt += 1
					sink <- generateSingleRecord(n1, n2, speedFwd, true)
				}
				if okBwd {
					bwdRecordCnt += 1
					sink <- generateSingleRecord(n1, n2, speedBwd, false)
				}

			}
		}
	}

	ds.Update(wayCnt, nodeCnt, fwdRecordCnt, bwdRecordCnt, wayMatched, nodeMatched, fwdTrafficMatched, bwdTrafficMatched)
	tasksWg.Done()
}

// format
// if dir = true, means traffic for forward, generate: from, to, speed
// if dir = false, means this speed is for backward flow, generate: to, from, speed
func generateSingleRecord(from, to uint64, speed int, dir bool) string {
	if dir {
		return fmt.Sprintf("%d,%d,%d\n", from, to, speed)
	} else {
		return fmt.Sprintf("%d,%d,%d\n", to, from, speed)
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
