package oasis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/serho/osrm-backend/integration/api/oasis"
	"github.com/serho/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/serho/osrm-backend/integration/service/oasis/osrmhelper"
	"github.com/serho/osrm-backend/integration/service/oasis/stationfinder"
	"github.com/golang/glog"
)

// Handler handles oasis request and provide response
type Handler struct {
	osrmConnector *osrmconnector.OSRMConnector
	finder        stationfinder.StationFinder
}

// New creates new Handler object
func New(osrmBackend, finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath string) (*Handler, error) {
	// @todo: need make sure connectivity is on and continues available
	//        simple request to guarantee server is alive after init
	if len(osrmBackend) == 0 {
		err := fmt.Errorf("empty osrmBackend end point")
		return nil, err
	}

	finder, err := stationfinder.CreateStationsFinder(finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath)
	if err != nil {
		glog.Errorf("Failed in Handler's New() when try to call CreateStationsFinder(), met error = %+v\n", err)
		return nil, err
	}

	return &Handler{
		osrmConnector: osrmconnector.NewOSRMConnector(osrmBackend),
		finder:        finder,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Handle incoming request %s from remote addr %s", req.RequestURI, req.RemoteAddr)

	// parse oasis request
	oasisReq, err := oasis.ParseRequestURL(req.URL)
	if err != nil || len(oasisReq.Coordinates) != 2 {
		glog.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// generate route response based on given oasis's orig/destination
	routeResp, err := osrmhelper.RequestRoute4InputOrigDest(oasisReq, h.osrmConnector)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// check whether orig and dest is reachable
	if len(routeResp.Routes) == 0 {
		info := "Orig and destination is not reachable for request " + oasisReq.RequestURI() + "."
		glog.Info(info)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, info)
		return
	}

	// check whether has enough energy
	b, remainRange, err := hasEnoughEnergy(oasisReq.CurrRange, oasisReq.SafeLevel, routeResp)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", err)
		return
	}
	if b {
		generateOASISResponse4NoChargeNeeded(w, routeResp, remainRange)
		return
	}

	// check whether could achieve by single charge
	overlap := getOverlapChargeStations4OrigDest(oasisReq, routeResp.Routes[0].Distance, h.osrmConnector, h.finder)
	if len(overlap) > 0 {
		generateResponse4SingleChargeStation(w, oasisReq, overlap, h.osrmConnector)
		return
	}

	// generate result for multiple charge
	generateSolutions4MultipleCharge(w, oasisReq, routeResp, h.osrmConnector, h.finder)
	return
}

func generateFakeOASISResponse(w http.ResponseWriter, req *oasis.Request) {
	w.WriteHeader(http.StatusOK)

	r := oasis.GenerateFakeResponse(req)

	json.NewEncoder(w).Encode(r)
}
