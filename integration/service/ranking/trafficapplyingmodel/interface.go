package trafficapplyingmodel

import "github.com/serho/osrm-backend/integration/api/osrm/route"

// Applier wraps interfaces for applying traffic on OSRM route.
type Applier interface {

	// ApplyTraffic applys traffic on a route.
	// liveTraffic and historicalSpeed are able to enable/disable the effects on the fly.
	ApplyTraffic(r *route.Route, liveTraffic bool, historicalSpeed bool) error
}
