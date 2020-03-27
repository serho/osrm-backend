package ranker

import "github.com/Telenav/osrm-backend/integration/service/spatialindexer"

type simpleRanker struct {
}

func (ranker *simpleRanker) RankPointIDsByGreatCircleDistance(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return rankPointsByGreatCircleDistanceToCenter(center, nearByIDs)
}

func (ranker *simpleRanker) RankPointIDsByShortestDistance(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return ranker.RankPointIDsByGreatCircleDistance(center, nearByIDs)
}
