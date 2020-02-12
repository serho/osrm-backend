package osrmconnector

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/pkg/backend"
	"github.com/golang/glog"
)

type osrmHTTPClient struct {
	osrmBackendEndpoint string
	httpclient          http.Client
	requestC            chan *request
}

func newOsrmHTTPClient(osrmEndpoint string) *osrmHTTPClient {
	return &osrmHTTPClient{
		osrmBackendEndpoint: osrmEndpoint,
		httpclient:          http.Client{Timeout: backend.Timeout()},
		requestC:            make(chan *request),
	}
}

func (oc *osrmHTTPClient) submitRouteReq(r *route.Request) <-chan RouteResponse {
	var url string
	if !strings.HasPrefix(oc.osrmBackendEndpoint, "http://") {
		url += "http://"
	}
	url = url + oc.osrmBackendEndpoint + r.RequestURI()

	req := newOsrmRequest(url, OSRMRoute)
	oc.requestC <- req
	glog.Info("[osrmHTTPClient]Submit route request " + url)
	return req.routeRespC
}

func (oc *osrmHTTPClient) submitTableReq(r *table.Request) <-chan TableResponse {
	var url string
	if !strings.HasPrefix(oc.osrmBackendEndpoint, "http://") {
		url += "http://"
	}
	url = url + oc.osrmBackendEndpoint + r.RequestURI()

	req := newOsrmRequest(url, OSRMTable)
	oc.requestC <- req
	glog.Info("[osrmHTTPClient]Submit table request " + url)
	return req.tableRespC
}

func (oc *osrmHTTPClient) start() {
	glog.Info("osrm connector started.\n")
	c := make(chan message)

	for {
		select {
		case req := <-oc.requestC:
			go oc.send(req, c)
		case m := <-c:
			go oc.response(&m)
		}
	}
}

type message struct {
	req  *request
	resp *http.Response
	err  error
}

func (oc *osrmHTTPClient) send(req *request, c chan<- message) {
	resp, err := oc.httpclient.Get(req.url)
	glog.Infof("[osrmHTTPClient] send function succeed with request %s" + req.url)
	m := message{req: req, resp: resp, err: err}
	c <- m
}

func (oc *osrmHTTPClient) response(m *message) {
	defer close(m.req.routeRespC)
	defer close(m.req.tableRespC)

	var routeResp RouteResponse
	var tableResp TableResponse
	if m.err != nil || m.resp == nil {
		glog.Warningf("osrm request %s failed, err %v\n", m.req.url, m.err)

		if m.req.t == OSRMRoute {
			routeResp.Err = m.err
			m.req.routeRespC <- routeResp
		} else if m.req.t == OSRMTable {
			tableResp.Err = m.err
			m.req.tableRespC <- tableResp
		} else {
			glog.Fatal("Unsupported request type for osrmHTTPClient")
		}

		return
	}
	defer m.resp.Body.Close()

	if m.req.t == OSRMRoute {
		routeResp.Err = json.NewDecoder(m.resp.Body).Decode(&routeResp.Resp)
		m.req.routeRespC <- routeResp
	} else if m.req.t == OSRMTable {
		tableResp.Err = json.NewDecoder(m.resp.Body).Decode(&tableResp.Resp)
		m.req.tableRespC <- tableResp
	} else {
		glog.Fatal("Unsupported request type for osrmHTTPClient")
	}
	glog.Infof("[osrmHTTPClient] Response for request %s" + m.req.url + " is generated.")

	return
}
