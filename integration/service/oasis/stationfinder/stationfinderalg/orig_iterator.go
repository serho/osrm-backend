package stationfinderalg

import (
	"github.com/serho/osrm-backend/integration/api/nav"
	"github.com/serho/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

type origIterator struct {
	location *nav.Location
}

// NewDestIter creates origIterator
// origIterator wraps single point of orig which adopts algorithms' requirement
func NewOrigIter(location *nav.Location) *origIterator {
	return &origIterator{
		location: location,
	}
}

func (oi *origIterator) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	c := make(chan *stationfindertype.ChargeStationInfo, 1)

	go func() {
		defer close(c)
		station := stationfindertype.ChargeStationInfo{
			ID:       stationfindertype.OrigLocationID,
			Location: *oi.location,
		}
		c <- &station
	}()

	return c
}
