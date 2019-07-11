package main

import (
	"testing"
	"os"
	"log"
	"fmt"
	"strconv"
	"io"
	"encoding/csv"
	"io/ioutil"
	"strings"
	"reflect"

	"github.com/Telenav/osrm-backend/traffic_updater/go/gen-go/proxy"
)

func TestSpeedTableDumper(t *testing.T) {
	way2nodeids := make(map[uint64][]int64)
	generateMockWay2nodeids2(way2nodeids)

	var flows []*proxy.Flow
	flows = loadMockTraffic("./testdata/mock-traffic.csv", flows)

	dumpSpeedTable4Customize(flows, way2nodeids, "./testdata/target.csv")

	compareFileContentUnstable("./testdata/target.csv", "./testdata/expect.csv", t)
}

func generateMockWay2nodeids2(way2nodeids map[uint64][]int64) {
	way2nodeids[24418325] = []int64{84760891102, 19496208102}
	way2nodeids[24418332] = []int64{84762609102,244183320001101,84762607102}
	way2nodeids[24418343] = []int64{84760849102,84760850102}
	way2nodeids[24418344] = []int64{84760846102,84760858102}
}

func loadMockTraffic(trafficPath string, flows []*proxy.Flow) []*proxy.Flow {
	// load mock traffic file
	mockfile, err := os.Open(trafficPath)
	defer mockfile.Close()
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Open pbf file of %v failed.\n", trafficPath)
		return nil
	}
	fmt.Printf("Open pbf file of %s succeed.\n", trafficPath)

	csvr := csv.NewReader(mockfile)
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				fmt.Printf("Error during decoding mock traffic, err = %v\n", err)
				return nil
			}
		}

		var wayid uint64
		var speed int64
		if wayid, err = strconv.ParseUint(row[0], 10, 64); err != nil {
			fmt.Printf("#Error during decoding wayid, row = %v\n", row)
		}
		if speed, err = strconv.ParseInt(row[1], 10, 32); err != nil {
			fmt.Printf("#Error during decoding speed, row = %v\n", row)
		}

		var flow proxy.Flow
		flow.WayId = (int64)(wayid)
		flow.Speed = (float64)(speed)
		flows = append(flows, &flow)
	}
	return flows
}

type tNodePair struct {
	f, t uint64
}

func loadSpeedCsv(f string, m map[tNodePair]int) {
		// load mock traffic file
		mockfile, err := os.Open(f)
		defer mockfile.Close()
		if err != nil {
			log.Fatal(err)
			fmt.Printf("Open file of %v failed.\n", f)
			return
		}
		fmt.Printf("Open file of %s succeed.\n", f)

		csvr := csv.NewReader(mockfile)
		for {
			row, err := csvr.Read()
			if err != nil {
				if err == io.EOF {
					err = nil
					break
				} else {
					fmt.Printf("Error during decoding file %s, err = %v\n", f, err)
					return
				}
			}

			var from, to uint64
			var speed int
			if from, err = strconv.ParseUint(row[0], 10, 64); err != nil {
				fmt.Printf("#Error during decoding, row = %v\n", row)
				return
			}
			if to, err = strconv.ParseUint(row[1], 10, 64); err != nil {
				fmt.Printf("#Error during decoding, row = %v\n", row)
				return
			}
			if speed, err = strconv.Atoi(row[2]); err != nil {
				fmt.Printf("#Error during decoding, row = %v\n", row)
				return
			}

			m[tNodePair{from, to}] = speed
		}
}

func compareFileContentStable(f1, f2 string, t *testing.T) {
	b1, err1 := ioutil.ReadFile(f1)
	if err1 != nil {
		fmt.Print(err1)
	}
	str1 := string(b1)

	b2, err2 := ioutil.ReadFile(f2)
	if err2 != nil {
		fmt.Print(err2)
	}
	str2 := string(b2)

	if strings.Compare(str1, str2) != 0 {
		t.Error("Compare file content failed\n")
	}
}

func compareFileContentUnstable(f1, f2 string, t *testing.T) {
	r1 := make(map[tNodePair]int)
	loadSpeedCsv(f1, r1)

	r2 := make(map[tNodePair]int)
	loadSpeedCsv(f2, r2)

	eq := reflect.DeepEqual(r1, r2)
	if !eq {
		t.Error("TestLoadWay2Nodeids failed to generate correct map\n")
	}

}

