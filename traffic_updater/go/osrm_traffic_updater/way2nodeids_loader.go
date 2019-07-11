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

func loadWay2NodeidsTable(filepath string, way2nodeids map[uint64][]int64, c chan <- bool) {
	startTime := time.Now()

	source := make(chan string)
	go load(filepath, source)
	convert(source, way2nodeids, c)

	endTime := time.Now()
	fmt.Printf("Processing time for loadWay2NodeidsTable takes %f seconds\n", endTime.Sub(startTime).Seconds())
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


// input data format
// wayid1, n1, (n2 - n1), (n3 - n2)...
// (wayid2 - wayid1), (n10 - n1), (n11 - n10), (n12 - n11) ...
func convert(source <-chan string, way2nodeids map[uint64][]int64, c chan<- bool) {
	var err error

	var preWayid, preNodeid int64
	for str := range source {
		elements := strings.Split(str, ",")
		if len(elements) < 3 {
			fmt.Printf("Invalid string %s in wayid2nodeids mapping file.\n", str)
			continue
		}

		var deltaWayid, wayid int64
		if deltaWayid, err = strconv.ParseInt(elements[0], 10, 64); err != nil {
			fmt.Printf("#Error during decoding wayid, row = %v\n", elements)
			continue
		}
		wayid = preWayid + deltaWayid
		preWayid = wayid

		var nodes []string = elements[1:]
		var firstNodeId int64
		for i := 0; i < len(nodes); i++ {
			var delta int64
			if delta, err = strconv.ParseInt(nodes[i], 10, 64); err != nil {
				fmt.Printf("#Error during decoding nodeid, row = %v\n", elements)
				continue
			}
			n := preNodeid + delta
			preNodeid = n
			if i == 0 {
				firstNodeId = n
			}

			way2nodeids[(uint64)(wayid)] = append(way2nodeids[(uint64)(wayid)], n)
		}
		preNodeid = firstNodeId
	}

	c <- true
}


