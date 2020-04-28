package ranking

import (
	"github.com/serho/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/serho/osrm-backend/integration/pkg/api/osrm/route/options"
)

func pickupRoutes(routes []*route.Route, num int) []*route.Route {
	if len(routes) <= num {
		return routes
	}
	return routes[:num]
}

func cleanupAnnotations(routes []*route.Route, annotations string) {
	if annotations != options.AnnotationsValueFalse {
		return // return all annotations even if want some
	}

	for _, route := range routes {
		for _, leg := range route.Legs {
			leg.Annotation = nil
		}
	}
}
