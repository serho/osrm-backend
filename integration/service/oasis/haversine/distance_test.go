package haversine

import (
	"testing"

	"github.com/serho/osrm-backend/integration/util"
)

func TestGreatCircleDistance(t *testing.T) {
	// expect value got from http://www.onlineconversion.com/map_greatcircle_distance.htm
	expect := 111595.4865288326
	actual := GreatCircleDistance(32.333, 122.323, 31.333, 122.423)
	if !util.FloatEquals(expect, actual) {
		t.Errorf("Expected GreatCircleDistance returns %v, got %v", expect, actual)
	}
}
