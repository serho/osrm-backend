package placeconnectivitymap

import (
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer/s2indexer"
	"github.com/golang/glog"
)

type PlaceConnectivityMap struct {
	finder   spatialindexer.Finder
	ranker   spatialindexer.Ranker
	iterator spatialindexer.PointsIterator
}

func NewPlaceConnectivityMap(filePath string) *PlaceConnectivityMap {
	indexer := s2indexer.NewS2Indexer().Build(filePath)
	if indexer == nil {
		glog.Error("Failed to NewPlaceConnectivityMap due to empty indexer is generated.  Check your input\n")
		return nil
	}

	return nil
}

func (pcm *PlaceConnectivityMap) GenerateConnectivityMap(filePath string) {
	glog.Info("Successfully finished GenerateConnectivityMap\n")
}
