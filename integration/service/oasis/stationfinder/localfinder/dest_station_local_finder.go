package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type destStationLocalFinder struct {
	basicFinder *basicLocalFinder
}

func newDestStationFinder(localFinder spatialindexer.Finder, oasisReq *oasis.Request) *destStationLocalFinder {
	if len(oasisReq.Coordinates) != 2 {
		glog.Errorf("Try to create newOrigStationFinder use incorrect oasis request, len(oasisReq.Coordinates) should be 2 but got %d.\n", len(oasisReq.Coordinates))
		return nil
	}
	if oasisReq.MaxRange <= oasisReq.SafeLevel {
		glog.Errorf("Try to create newOrigStationFinder use incorrect oasis request, SafeLevel should be smaller than MaxRange.\n")
		return nil
	}

	obj := &destStationLocalFinder{
		basicFinder: newBasicLocalFinder(localFinder),
	}
	obj.basicFinder.getNearbyChargeStations(spatialindexer.Location{
		Lat: oasisReq.Coordinates[1].Lat,
		Lon: oasisReq.Coordinates[1].Lon},
		oasisReq.MaxRange-oasisReq.SafeLevel)

	return obj
}

// NearbyStationsIterator provides channel which contains near by station information for dest
func (localFinder *destStationLocalFinder) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	return localFinder.basicFinder.IterateNearbyStations()
}

// Stop stops functionality of finder
func (localFinder *destStationLocalFinder) Stop() {
	localFinder.basicFinder.Stop()
}
