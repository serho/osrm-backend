package ranking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/code"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route/options"
	"github.com/Telenav/osrm-backend/integration/service/ranking/strategy/rankbyduration"

	"github.com/golang/glog"
)

// Handler represents a handler for ranking.
type Handler struct {
	nodes2WayQuerier waysnodes.WaysQuerier
	trafficApplier   trafficapplyingmodel.Applier

	osrmBackend string
}

// New creates a new handler for ranking.
func New(osrmBackend string, nodes2WayQuerier waysnodes.WaysQuerier, trafficApplier trafficapplyingmodel.Applier) *Handler {
	if nodes2WayQuerier == nil {
		glog.Fatal("nil nodes2WayQuerier")
		return nil
	}

	return &Handler{
		nodes2WayQuerier,
		trafficApplier,
		osrmBackend,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Handle incoming request %s from remote addr %s", req.RequestURI, req.RemoteAddr)

	// parse incoming request
	osrmRequest, err := route.ParseRequestURL(req.URL)
	if err != nil {
		glog.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// support traffic options on the fly
	originalRequestValues := api.ParseQueryDiscardError(req.URL.RawQuery)
	enableLiveTraffic, enableHistoricalSpeed := parseTrafficOptions(originalRequestValues.Get(liveTrafficQueryKey), originalRequestValues.Get(historicalSpeedQueryKey))

	// modify
	originalAlternativesNum := osrmRequest.AlternativesNumber()
	originalAnnotations := osrmRequest.Annotations
	osrmRequest.Alternatives = strconv.FormatUint(uint64(flags.alternatives), 10)
	osrmRequest.Annotations = options.AnnotationsValueTrue

	// route against backend OSRM
	osrmResponse, osrmHTTPStatus, err := h.routeByOSRM(osrmRequest)
	w.WriteHeader(osrmHTTPStatus)
	if err != nil {
		glog.Warning(err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if osrmResponse.Code == code.OK {
		if err := h.retrieveWayIDs(osrmResponse.Routes); err != nil {
			glog.Warning(err)
			fmt.Fprintf(w, "Retrieve ways from nodes failed, err: %v", err)
			return
		}

		if h.trafficApplier != nil {
			// update speeds,durations,datasources by traffic
			osrmResponse.Routes, err = h.updateRoutesByTraffic(osrmResponse.Routes, enableLiveTraffic, enableHistoricalSpeed)
			if err != nil {
				glog.Warning(err)
				fmt.Fprintf(w, "Apply traffic on routes failed, err: %v", err)
				return
			}
		}

		// rank
		osrmResponse.Routes = rankbyduration.Rank(osrmResponse.Routes)

		// pick up
		osrmResponse.Routes = pickupRoutes(osrmResponse.Routes, originalAlternativesNum)

		// cleanup annotations if necessary
		cleanupAnnotations(osrmResponse.Routes, originalAnnotations)
	}

	// return
	json.NewEncoder(w).Encode(osrmResponse)
}
