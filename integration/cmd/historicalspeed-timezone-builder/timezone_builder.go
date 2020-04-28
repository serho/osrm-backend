package main

import (
	"time"

	"github.com/serho/osrm-backend/integration/util/unidbpatch"
	"github.com/golang/glog"
	"github.com/qedus/osmpbf"
)

type wayTimezoneInfo struct {
	wayID          int64
	timezone       string
	daylightSaving string
}

func newTimezoneBuilder(in <-chan *osmpbf.Way, out chan<- *wayTimezoneInfo) {
	startTime := time.Now()

	var waysCount, invalidWaysCount, noTimezoneInfoWaysCount, succeedWaysCount int
	for {
		way, ok := <-in
		if !ok {
			break
		}
		waysCount++

		/// The parse processing only support UniDB now since OSM almost has no timezone information.
		if !unidbpatch.IsValidWay(way.ID) {
			invalidWaysCount++
			continue
		}

		tz := wayTimezoneInfo{}
		tz.wayID = unidbpatch.TrimValidWayIDSuffix(way.ID)
		tz.timezone = way.Tags["timezone:left"]
		if len(tz.timezone) == 0 {
			tz.timezone = way.Tags["timezone:right"]
		}
		tz.daylightSaving = way.Tags["dst_pattern"]

		if len(tz.timezone) == 0 && len(tz.daylightSaving) == 0 {
			noTimezoneInfoWaysCount++
			continue // no valid timezone and daylight saving, ignore it
		}

		out <- &tz

		succeedWaysCount++
	}
	glog.Infof("Built timezone info, total ways %d, succeed ways %d, invalid ways %d, no timezone and daylight saving ways %d, takes %f seconds", waysCount, succeedWaysCount, invalidWaysCount, noTimezoneInfoWaysCount, time.Now().Sub(startTime).Seconds())
}
