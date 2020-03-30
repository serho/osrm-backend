package ranker

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/genericoptions"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
)

func TestGenerateTableRequest(t *testing.T) {
	cases := []struct {
		center     spatialindexer.Location
		nearByIDs  []*spatialindexer.PointInfo
		startIndex int
		endIndex   int
		expect     *table.Request
	}{
		// case 1: test 0 -> {1, 2, 3, 4, 5}
		{
			center: spatialindexer.Location{
				Lat: 0,
				Lon: 0,
			},
			nearByIDs: []*spatialindexer.PointInfo{
				&spatialindexer.PointInfo{
					ID: 1,
					Location: spatialindexer.Location{
						Lat: 1.1,
						Lon: 1.1,
					},
				},
				&spatialindexer.PointInfo{
					ID: 2,
					Location: spatialindexer.Location{
						Lat: 2.2,
						Lon: 2.2,
					},
				},
				&spatialindexer.PointInfo{
					ID: 3,
					Location: spatialindexer.Location{
						Lat: 3.3,
						Lon: 3.3,
					},
				},
				&spatialindexer.PointInfo{
					ID: 4,
					Location: spatialindexer.Location{
						Lat: 4.4,
						Lon: 4.4,
					},
				},
				&spatialindexer.PointInfo{
					ID: 5,
					Location: spatialindexer.Location{
						Lat: 5.5,
						Lon: 5.5,
					},
				},
			},
			startIndex: 0,
			endIndex:   4,
			expect: &table.Request{
				Service: "table",
				Version: "v1",
				Profile: "driving",
				Coordinates: coordinate.Coordinates{
					coordinate.Coordinate{
						Lat: 0,
						Lon: 0,
					},
					coordinate.Coordinate{
						Lat: 1.1,
						Lon: 1.1,
					},
					coordinate.Coordinate{
						Lat: 2.2,
						Lon: 2.2,
					},
					coordinate.Coordinate{
						Lat: 3.3,
						Lon: 3.3,
					},
					coordinate.Coordinate{
						Lat: 4.4,
						Lon: 4.4,
					},
					coordinate.Coordinate{
						Lat: 5.5,
						Lon: 5.5,
					},
				},
				Sources: genericoptions.Elements{
					"0",
				},
				Destinations: genericoptions.Elements{
					"1",
					"2",
					"3",
					"4",
					"5",
				},
				Annotations: "distance,duration",
			},
		},
		// case 2: test 0 -> {2, 3, 4}
		{
			center: spatialindexer.Location{
				Lat: 0,
				Lon: 0,
			},
			nearByIDs: []*spatialindexer.PointInfo{
				&spatialindexer.PointInfo{
					ID: 1,
					Location: spatialindexer.Location{
						Lat: 1.1,
						Lon: 1.1,
					},
				},
				&spatialindexer.PointInfo{
					ID: 2,
					Location: spatialindexer.Location{
						Lat: 2.2,
						Lon: 2.2,
					},
				},
				&spatialindexer.PointInfo{
					ID: 3,
					Location: spatialindexer.Location{
						Lat: 3.3,
						Lon: 3.3,
					},
				},
				&spatialindexer.PointInfo{
					ID: 4,
					Location: spatialindexer.Location{
						Lat: 4.4,
						Lon: 4.4,
					},
				},
				&spatialindexer.PointInfo{
					ID: 5,
					Location: spatialindexer.Location{
						Lat: 5.5,
						Lon: 5.5,
					},
				},
			},
			startIndex: 1,
			endIndex:   3,
			expect: &table.Request{
				Service: "table",
				Version: "v1",
				Profile: "driving",
				Coordinates: coordinate.Coordinates{
					coordinate.Coordinate{
						Lat: 0,
						Lon: 0,
					},
					coordinate.Coordinate{
						Lat: 2.2,
						Lon: 2.2,
					},
					coordinate.Coordinate{
						Lat: 3.3,
						Lon: 3.3,
					},
					coordinate.Coordinate{
						Lat: 4.4,
						Lon: 4.4,
					},
				},
				Sources: genericoptions.Elements{
					"0",
				},
				Destinations: genericoptions.Elements{
					"2",
					"3",
					"4",
				},
				Annotations: "distance,duration",
			},
		},
	}

	for _, c := range cases {
		actual := generateTableRequest(c.center, c.nearByIDs, c.startIndex, c.endIndex)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During TestGenerateTableRequest, expect table request is \n%+v\n but actual is \n%+v\n", c.expect, actual)
		}
	}

}
