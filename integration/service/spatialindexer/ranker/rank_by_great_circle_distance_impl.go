package ranker

import (
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/blevesearch/bleve/geo"
	"github.com/golang/glog"
)

func rankPointsByGreatCircleDistanceToCenter(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	if len(nearByIDs) == 0 {
		glog.Warning("When try to rankPointsByGreatCircleDistanceToCenter, input array is empty\n")
		return nil
	}

	pointWithDistanceC := make(chan *spatialindexer.RankedPointInfo, len(nearByIDs))
	go func() {
		defer close(pointWithDistanceC)

		for _, p := range nearByIDs {
			pointWithDistanceC <- &spatialindexer.RankedPointInfo{
				PointInfo: spatialindexer.PointInfo{
					ID:       p.ID,
					Location: p.Location,
				},
				// geo.Haversin's unit is kilometer, convert to meter
				Distance: geo.Haversin(center.Lon, center.Lat, p.Location.Lon, p.Location.Lat) * 1000,
			}
		}
	}()

	rankAgent := newRankAgent(len(nearByIDs))
	return rankAgent.RankByDistance(pointWithDistanceC)
}
