package ranker

import (
	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
)

// func rankPointsByGreatCircleDistanceToCenter(center spatialindexer.Location, nearByIDs []*spatialindexer.PointInfo, oc *osrmconnector.OSRMConnector) []*spatialindexer.RankedPointInfo {
// 	if len(nearByIDs) == 0 {
// 		glog.Warning("When try to rankPointsByGreatCircleDistanceToCenter, input array is empty\n")
// 		return nil
// 	}

// 	pointWithDistanceC := make([]*spatialindexer.RankedPointInfo, len(nearByIDs))
// }

func generateTableRequest(center spatialindexer.Location,
	nearByIDs []*spatialindexer.PointInfo,
	oc *osrmconnector.OSRMConnector,
	outputC chan *spatialindexer.RankedPointInfo,
	startIndex int) {
	// todo
}
