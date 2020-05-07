package cloudfinder

import (
	"github.com/serho/osrm-backend/integration/api/oasis"
	"github.com/serho/osrm-backend/integration/api/search/searchcoordinate"
	"github.com/serho/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/serho/osrm-backend/integration/service/oasis/searchhelper"
)

//@todo: This number need to be adjusted based on charge station profile
const origMaxSearchCandidateNumber = 999

type origStationFinder struct {
	oasisReq *oasis.Request
	*basicFinder
}

func newOrigStationFinder(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *origStationFinder {
	obj := &origStationFinder{
		oasisReq,
		newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (oFinder *origStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: oFinder.oasisReq.Coordinates[0].Lat,
			Lon: oFinder.oasisReq.Coordinates[0].Lon},
		origMaxSearchCandidateNumber,
		oFinder.oasisReq.CurrRange)

	oFinder.getNearbyChargeStations(req)
	return
}
