package connectivitymap

import (
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/golang/glog"
)

// IDAndDistance wraps ID and distance information
type IDAndDistance struct {
	ID       spatialindexer.PointID
	Distance float64
}

// ID2NearByIDsMap is a mapping between ID and its nearby IDs
type ID2NearByIDsMap map[spatialindexer.PointID][]IDAndDistance

type ConnectivityMap struct {
	id2nearByIDs       ID2NearByIDsMap
	distanceLimitation float64
	statistic          *statistic
}

func (cm *ConnectivityMap) DistanceLimitation() float64 {
	return cm.distanceLimitation
}

func NewPlaceConnectivityMap(distanceLimitation float64) *ConnectivityMap {
	return &ConnectivityMap{
		distanceLimitation: distanceLimitation,
		statistic:          newStatistic(),
	}
}

func (cm *ConnectivityMap) Build() {
	glog.Info("Successfully finished GenerateConnectivityMap\n")
}

func (cm *ConnectivityMap) Dump(folderPath string) {
}

func (cm *ConnectivityMap) Load(folderPath string) {
}

func (cm *ConnectivityMap) QueryConnectivity(placeInfo spatialindexer.PointInfo, limitDistance float64) {
	// for each everything recorded in data, apply limit option on that
}
