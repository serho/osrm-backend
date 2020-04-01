package connectivitymap

import "github.com/Telenav/osrm-backend/integration/service/spatialindexer"

type connectivityMapBuilder struct {
	iterator      spatialindexer.PointsIterator
	finder        spatialindexer.Finder
	ranker        spatialindexer.Ranker
	distanceLimit float64
}

func newConnectivityMapBuilder(iterator spatialindexer.PointsIterator, finder spatialindexer.Finder,
	ranker spatialindexer.Ranker, distanceLimit float64) *connectivityMapBuilder {
	return &connectivityMapBuilder{
		iterator:      iterator,
		finder:        finder,
		ranker:        ranker,
		distanceLimit: distanceLimit,
	}
}

type placeIDWithNearByPlaceIDs struct {
	id  spatialindexer.PointID
	ids []IDAndDistance
}

// todo: use task pool
func (builder *connectivityMapBuilder) build() ID2NearByIDsMap {
	internalResult := make(chan placeIDWithNearByPlaceIDs, 10000)
	m := make(ID2NearByIDsMap)

	go func() {
		for p := range builder.iterator.IteratePoints() {
			nearbyIDs := builder.finder.FindNearByPointIDs(p.Location, builder.distanceLimit, spatialindexer.UnlimitedCount)
			rankedResults := builder.ranker.RankPointIDsByGreatCircleDistance(p.Location, nearbyIDs)

			ids := make([]IDAndDistance, 0, len(rankedResults))
			for _, r := range rankedResults {
				ids = append(ids, IDAndDistance{
					ID:       r.ID,
					Distance: r.Distance,
				})
			}
			internalResult <- placeIDWithNearByPlaceIDs{
				id:  p.ID,
				ids: ids,
			}
		}
		close(internalResult)
	}()

	for item := range internalResult {
		m[item.id] = item.ids
	}

	return m
}

// func (builder *connectivityMapBuilder) dump(folderPath string, m ID2NearByIDsMap) {

// }

// func (builder *connectivityMapBuilder) load(folderPath string) (ID2NearByIDsMap, float64) {

// }
