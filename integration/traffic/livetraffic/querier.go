package livetraffic

import (
	"github.com/serho/osrm-backend/integration/traffic"
)

// Querier defines interfaces for querying traffic flows and incidents.
type Querier traffic.LiveTrafficQuerier
