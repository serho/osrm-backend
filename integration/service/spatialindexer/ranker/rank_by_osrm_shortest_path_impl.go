package ranker

import (
	"strconv"
	"sync"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/golang/glog"
)

// pointsLimit4SingleTableRequest defines point count limitation for each single table request.
// During pre-processing, its possible to have situation to calculate distance between thousnads of points.
// The situation here is 1-to-N table request, use pointsLimit4SingleTableRequest to limit N
const pointsLimit4SingleTableRequest = 1000

func rankPointsByOSRMShortestPath(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo, oc *osrmconnector.OSRMConnector) []*spatialindexer.RankedPointInfo {
	if len(nearByIDs) == 0 {
		glog.Warning("When try to rankPointsByGreatCircleDistanceToCenter, input array is empty\n")
		return nil
	}

	var wg sync.WaitGroup
	pointWithDistanceC := make(chan *spatialindexer.RankedPointInfo, len(nearByIDs))
	startIndex := 0
	endIndex := 0
	for {
		if startIndex >= len(nearByIDs) {
			break
		}
		endIndex = startIndex + pointsLimit4SingleTableRequest
		if endIndex >= len(nearByIDs) {
			endIndex = len(nearByIDs) - 1
		}

		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()

			rankedPoints, err := calcShortestPathDistance(center, nearByIDs, oc, startIndex, endIndex)

			if err != nil {
				// @todo: add retry logic when failed
			} else {
				for _, item := range rankedPoints {
					pointWithDistanceC <- item
				}
			}

		}(&wg)

		startIndex = endIndex + 1
	}

	wg.Wait()
	close(pointWithDistanceC)

	rankAgent := newRankAgent(len(nearByIDs))
	return rankAgent.RankByDistance(pointWithDistanceC)

}

func calcShortestPathDistance(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo, oc *osrmconnector.OSRMConnector, startIndex, endIndex int) ([]*spatialindexer.RankedPointInfo, error) {
	req := generateTableRequest(center, nearByIDs, startIndex, endIndex)
	respC := oc.Request4Table(req)
	resp := <-respC

	if resp.Err != nil {
		glog.Errorf("Failed to generate table response for \n %s with \n err =%v \n", req.RequestURI(), resp.Err)
		return nil, resp.Err
	}

	result := make([]*spatialindexer.RankedPointInfo, 0, endIndex-startIndex+1)
	for i := 0; i < endIndex-startIndex+1; i++ {
		result = append(result, &spatialindexer.RankedPointInfo{
			PointInfo: spatialindexer.PointInfo{
				ID:       nearByIDs[startIndex+i].ID,
				Location: nearByIDs[startIndex+i].Location,
			},
			Distance: *resp.Resp.Distances[0][i],
		})
	}
	return result, nil
}

// generateTableRequest generates table requests from center to [startIndex, endIndex] of nearByIDs
func generateTableRequest(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo, startIndex, endIndex int) *table.Request {
	if startIndex < 0 || startIndex > endIndex || endIndex >= len(nearByIDs) {
		glog.Fatalf("startIndex should be smaller equal to endIndex and both of them should in the range of len(nearByIDs), while (startIndex, endIndex, len(nearByIDs)) = (%d, %d, %d)",
			startIndex, endIndex, len(nearByIDs))
	}

	req := table.NewRequest()
	req.Coordinates = append(ConvertLocation2Coordinates(center),
		ConvertPointInfos2Coordinates(nearByIDs, startIndex, endIndex)...)

	req.Sources = append(req.Sources, strconv.Itoa(0))
	pointsCount4Sources := 1
	for i := startIndex; i < endIndex; i++ {
		str := strconv.Itoa(endIndex - startIndex + pointsCount4Sources)
		req.Destinations = append(req.Destinations, str)
	}

	return req
}

func ConvertLocation2Coordinates(location spatialindexer.Location) coordinate.Coordinates {
	result := make(coordinate.Coordinates, 0, 1)
	result = append(result, coordinate.Coordinate{
		Lat: location.Lat,
		Lon: location.Lon,
	})
	return result
}

func ConvertPointInfos2Coordinates(nearByIDs []*spatialindexer.PointInfo, startIndex, endIndex int) coordinate.Coordinates {
	result := make(coordinate.Coordinates, 0, endIndex-startIndex+1)
	for i := startIndex; i <= endIndex; i++ {
		result = append(result, coordinate.Coordinate{
			Lat: nearByIDs[i].Location.Lat,
			Lon: nearByIDs[i].Location.Lon,
		})
	}
	return result
}
