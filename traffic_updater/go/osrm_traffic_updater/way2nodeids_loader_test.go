package main

import (
	"testing"
	"reflect"
	"sync"
	"fmt"
)


func TestLoadWay2Nodeids(t *testing.T) {
	// load result into sources
	var sources [TASKNUM]chan way2Nodes
	for i := range sources {
		//fmt.Printf("&&& current i is %d\n", i)
		sources[i] = make(chan way2Nodes, 10000)
	}
	go loadWay2NodeidsTable("./testdata/id-mapping-delta.csv.snappy", sources)

	allWay2NodesChan := make(chan way2Nodes, 10000)
	var wgs sync.WaitGroup
	wgs.Add(TASKNUM)
	for i:= 0; i < TASKNUM; i++ {
		//fmt.Printf("### current i is %d\n", i)
		go mergeChannels(sources[i], allWay2NodesChan, &wgs)
	}
	wgs.Wait()

	// dump result into map
	way2nodeids := make(map[uint64][]int64)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for elem := range allWay2NodesChan {
			way2nodeids[elem.w] = elem.nodes
		}
		wg.Done()
	}()

	close(allWay2NodesChan)
	wg.Wait()

	// test map result
	way2nodeidsExpect := make(map[uint64][]int64)
	generateMockWay2nodeids(way2nodeidsExpect)
	eq := reflect.DeepEqual(way2nodeids, way2nodeidsExpect)
	if !eq {
		t.Error("TestLoadWay2Nodeids failed to generate correct map\n")
	}
}

func mergeChannels(f <-chan way2Nodes, t chan<- way2Nodes, w *sync.WaitGroup) {
	fmt.Printf("Enter mergeChannels\n")
	for elem := range f {
		fmt.Printf("@@@ in merge channel %v\n", elem)
		t <- elem
	}
	fmt.Printf("merge channel is done\n")
	w.Done()
	fmt.Printf("$$$ set done\n")
}



func generateMockWay2nodeids(way2nodeids map[uint64][]int64) {
	way2nodeids[24418325] = []int64{84760891102, 19496208102}
	way2nodeids[24418332] = []int64{84762609102,244183320001101,84762607102}
	way2nodeids[24418343] = []int64{84760849102,84760850102}
	way2nodeids[24418344] = []int64{84760846102,84760858102}
}

