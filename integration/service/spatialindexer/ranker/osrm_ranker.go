package ranker

import (
	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
)

type osrmRanker struct {
	oc *osrmconnector.OSRMConnector
}

func newOsrmRanker(oc *osrmconnector.OSRMConnector) *osrmRanker {
	return &osrmRanker{
		oc: oc,
	}
}

func (ranker *osrmRanker) RankPointIDsByGreatCircleDistance(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return rankPointsByGreatCircleDistanceToCenter(center, nearByIDs)
}

func (ranker *osrmRanker) RankPointIDsByShortestDistance(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return rankPointsByOSRMShortestPath(center, nearByIDs, ranker.oc)
}
