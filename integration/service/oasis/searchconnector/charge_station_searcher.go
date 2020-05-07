package searchconnector

import "github.com/serho/osrm-backend/integration/api/search/nearbychargestation"

// ChargeStationsResponse contains charge stations response and error
type ChargeStationsResponse struct {
	Resp *nearbychargestation.Response
	Err  error
}

// ChargeStationSearcher is the interface provides ability to search nearby charge station.
type ChargeStationSearcher interface {
	ChargeStationSearch(req *nearbychargestation.Request) <-chan ChargeStationsResponse
}
