package spatialindexer

// Location for poi point
// @todo: will be replaced by the one in map
type Location struct {
	Latitude  float64
	Longitude float64
}

type PointInfo struct {
	ID       PointID
	Location Location
}

type RankedPointInfo struct {
	PointInfo
	Distance float64
}

type PointID int64

type Finder interface {
	FindNearByIDs(center Location, radius float64, limitCount int) []PointInfo
}

type Ranker interface {
	RankingIDsByGreatCircleDistance(center Location, nearByIDs []PointInfo) []RankedPointInfo

	RankingIDsByShortestDistance(center Location, nearByIDs []PointInfo) []RankedPointInfo
}
