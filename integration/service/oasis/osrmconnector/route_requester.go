package osrmconnector

import "github.com/serho/osrm-backend/integration/api/osrm/route"

type RouteRequster interface {
	Request4Route(r *route.Request) <-chan RouteResponse
}
