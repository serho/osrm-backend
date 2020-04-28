package cloudfinder

import (
	"github.com/serho/osrm-backend/integration/pkg/api/oasis"
	"github.com/serho/osrm-backend/integration/pkg/api/search/searchcoordinate"
	"github.com/serho/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/serho/osrm-backend/integration/service/oasis/searchhelper"
)

//@todo: This number need to be adjusted based on charge station profile
const destMaxSearchCandidateNumber = 999

type destStationFinder struct {
	oasisReq *oasis.Request
	*basicFinder
}

func newDestStationFinder(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *destStationFinder {
	obj := &destStationFinder{
		oasisReq,
		newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (dFinder *destStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: dFinder.oasisReq.Coordinates[1].Lat,
			Lon: dFinder.oasisReq.Coordinates[1].Lon},
		destMaxSearchCandidateNumber,
		dFinder.oasisReq.MaxRange-dFinder.oasisReq.SafeLevel)
	dFinder.getNearbyChargeStations(req)
}
