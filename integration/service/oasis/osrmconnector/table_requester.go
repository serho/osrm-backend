package osrmconnector

import "github.com/serho/osrm-backend/integration/pkg/api/osrm/table"

type TableRequster interface {
	Request4Table(r *table.Request) <-chan TableResponse
}
