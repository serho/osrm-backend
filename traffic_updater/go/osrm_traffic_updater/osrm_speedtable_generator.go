package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/golang/snappy"
)

// todo:
//       Write data into more compressed format(parquet)
//       Statistic to avoid unmatched element
//       Multiple go routine for convert()
func generateSpeedTable(wayid2speed map[uint64]int, way2nodeidsPath string, target string) {
	startTime := time.Now()

	// format is: wayid, nodeid, nodeid, nodeid...
	source := make(chan string)
	// format is fromid, toid, speed
	sink := make(chan string)

	go load(way2nodeidsPath, source)
	go convert(source, sink, wayid2speed)
	write(target, sink)

	endTime := time.Now()
	fmt.Printf("Processing time for generate speed table takes %f seconds\n", endTime.Sub(startTime).Seconds())
}

func load(mappingPath string, source chan<- string) {
	defer close(source)

	f, err := os.Open(mappingPath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Open idsmapping file of %v failed.\n", mappingPath)
		return
	}
	fmt.Printf("Open idsmapping file of %s succeed.\n", mappingPath)

	scanner := bufio.NewScanner(snappy.NewReader(f))
	for scanner.Scan() {
		source <- (scanner.Text())
	}
}


// data format
// wayid1, n1, (n2 - n1), (n3 - n2)...
// (wayid2 - wayid1), (n10 - n1), (n11 - n10), (n12 - n11) ...
func convert(source <-chan string, sink chan<- string, wayid2speed map[uint64]int) {
	var err error
	defer close(sink)

	var preWayid, preNodeid int64
	for str := range source {
		elements := strings.Split(str, ",")
		if len(elements) < 3 {
			fmt.Printf("Invalid string %s in wayid2nodeids mapping file\n", str)
			continue
		}

		var deltaWayid, wayid int64
		if deltaWayid, err = strconv.ParseInt(elements[0], 10, 64); err != nil {
			fmt.Printf("#Error during decoding wayid, row = %v\n", elements)
			continue
		}
		wayid = preWayid + deltaWayid
		preWayid = wayid

		var firstNodeId int64
		if speed, ok := wayid2speed[(uint64)(wayid)]; ok {
			var nodes []string = elements[1:]
			for i := 0; (i + 1) < len(nodes); i++ {
				var deltaN1, deltaN2 int64
				if deltaN1, err = strconv.ParseInt(nodes[i], 10, 64); err != nil {
					fmt.Printf("#Error during decoding nodeid, row = %v\n", elements)
					continue
				}
				n1 := preNodeid + deltaN1
				preNodeid = n1
				if i == 0 {
					firstNodeId = n1
				}

				if deltaN2, err = strconv.ParseInt(nodes[i+1], 10, 64); err != nil {
					fmt.Printf("#Error during decoding nodeid, row = %v\n", elements)
					continue
				}
				n2 := preNodeid + deltaN2

				var s string
				if speed >= 0 {
					s = fmt.Sprintf("%d,%d,%d\n", n1, n2, speed)
				} else {
					s = fmt.Sprintf("%d,%d,%d\n", n2, n1, -speed)
				}

				sink <- s
			}
			preNodeid = firstNodeId
		} else {
			var delta int64
			if delta, err = strconv.ParseInt(elements[1], 10, 64); err != nil {
				fmt.Printf("#Error during decoding nodeid, row = %v\n", elements)
				continue
			}
			preNodeid = preNodeid + delta
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
