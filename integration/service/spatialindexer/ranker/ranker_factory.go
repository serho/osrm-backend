package ranker

import (
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/Telenav/osrm-backend2/integration/oasis/osrmconnector"
)

const (
	SimpleRanker    = "SimpleRanker"
	OSRMBasedRanker = "OSRMBasedRanker"
)

func CreateRanker(rankerType string, oc *osrmconnector.OSRMConnector) spatialindexer.Ranker {
	switch rankerType {
	case SimpleRanker:
		return nil
	case OSRMBasedRanker:
		return nil
	default:
		return nil
	}
}
