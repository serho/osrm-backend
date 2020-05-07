package ranker

import (
	"sort"

	"github.com/serho/osrm-backend/integration/service/oasis/spatialindexer"
)

// rankAgent accepts items to be ranked then returns ranking result
type rankAgent struct {
	rankedPoints []*spatialindexer.RankedPointInfo
}

func newRankAgent(pointNum int) *rankAgent {
	return &rankAgent{
		rankedPoints: make([]*spatialindexer.RankedPointInfo, 0, pointNum),
	}
}

type rankItems []*spatialindexer.RankedPointInfo

func (r rankItems) Len() int {
	return len(r)
}

func (r rankItems) Less(i, j int) bool {
	return r[i].Distance < r[j].Distance
}

func (r rankItems) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r *rankAgent) RankByDistance(input <-chan *spatialindexer.RankedPointInfo) []*spatialindexer.RankedPointInfo {
	for p := range input {
		r.rankedPoints = append(r.rankedPoints, p)
	}

	sort.Sort(rankItems(r.rankedPoints))

	return r.rankedPoints
}
