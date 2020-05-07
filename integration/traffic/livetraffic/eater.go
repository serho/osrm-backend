package livetraffic

import "github.com/serho/osrm-backend/integration/traffic/livetraffic/trafficproxy"

// Eater is the interface that wraps the basic Eat method.
type Eater interface {

	// Eat consumes traffic responses.
	Eat(trafficproxy.TrafficResponse)
}
