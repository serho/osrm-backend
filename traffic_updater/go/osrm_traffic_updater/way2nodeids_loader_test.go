package main

import (
	"testing"
	"reflect"
	"time"
)

func TestLoadWay2Nodeids(t *testing.T) {
	
	way2nodeids := make(map[uint64][]int64) 
	c := make(chan bool, 1)
	defer close(c)
	loadWay2NodeidsTable("./testdata/id-mapping-delta.csv.snappy", way2nodeids, c)

	// test map result
	way2nodeidsExpect := make(map[uint64][]int64)
	generateMockWay2nodeids(way2nodeidsExpect)
	eq := reflect.DeepEqual(way2nodeids, way2nodeidsExpect)
	if !eq {
		t.Error("TestLoadWay2Nodeids failed to generate correct map\n")
	}

	// test channel
	select {
	case b, ok := <-c:
		if !ok || !b {
			t.Error("TestLoadWay2Nodeids failed to set channel result")
		} 
	case <-time.After(1 * time.Second):
			t.Error("TestLoadWay2Nodeids timeout")
	}
}

func generateMockWay2nodeids(way2nodeids map[uint64][]int64) {
	way2nodeids[24418325] = []int64{84760891102, 19496208102}
	way2nodeids[24418332] = []int64{84762609102,244183320001101,84762607102}
	way2nodeids[24418343] = []int64{84760849102,84760850102}
	way2nodeids[24418344] = []int64{84760846102,84760858102}
}

