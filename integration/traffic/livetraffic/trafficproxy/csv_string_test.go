package trafficproxy

import "testing"

func TestFlowCSVString(t *testing.T) {

	cases := []struct {
		f                      Flow
		csvString              string
		humanFriendlyCSVString string
	}{
		{
			Flow{WayID: 829733412, Speed: 20.280001, TrafficLevel: TrafficLevel_FREE_FLOW, Timestamp: 1579419488000, Offset: 2, Limit: 3},
			"829733412,20.280001,7,1579419488000,2,3",
			"829733412,20.280001,FREE_FLOW,1579419488000,2,3",
		},
		{
			Flow{WayID: -129639168, Speed: 31.389999, TrafficLevel: TrafficLevel_FREE_FLOW, Timestamp: 1579419488000, Offset: 2, Limit: 3},
			"-129639168,31.389999,7,1579419488000,2,3",
			"-129639168,31.389999,FREE_FLOW,1579419488000,2,3",
		},
	}

	for _, c := range cases {
		cs := c.f.CSVString()
		if cs != c.csvString {
			t.Errorf("flow: %v, expect csv string %s, but got %s", c.f, c.csvString, cs)
		}

		hs := c.f.HumanFriendlyCSVString()
		if hs != c.humanFriendlyCSVString {
			t.Errorf("flow: %v, expect human friendly csv string %s, but got %s", c.f, c.humanFriendlyCSVString, hs)
		}
	}

}

func TestIncidentCSVString(t *testing.T) {

	cases := []struct {
		incident               Incident
		s                      string
		humanFriendlyCSVString string
	}{
		{
			Incident{
				IncidentID:            "TTI-f47b8dba-59a3-372d-9cec-549eb252e2d5-TTR46312939215361-1",
				AffectedWayIDs:        []int64{100663296, -1204020275, 100663296, -1204020274, 100663296, -916744017, 100663296, -1204020245, 100663296, -1194204646, 100663296, -1204394608, 100663296, -1194204647, 100663296, -129639168, 100663296, -1194204645},
				IncidentType:          IncidentType_MISCELLANEOUS,
				IncidentSeverity:      IncidentSeverity_CRITICAL,
				IncidentLocation:      &Location{Lat: 44.181220, Lon: -117.135840},
				Description:           "Construction on I-84 EB near MP 359, Drive with caution.",
				FirstCrossStreet:      "",
				SecondCrossStreet:     "",
				StreetName:            "I-84 E",
				EventCode:             500,
				AlertCEventQuantifier: 0,
				IsBlocking:            false,
				Timestamp:             1579419488000,
			},
			"TTI-f47b8dba-59a3-372d-9cec-549eb252e2d5-TTR46312939215361-1,\"100663296,-1204020275,100663296,-1204020274,100663296,-916744017,100663296,-1204020245,100663296,-1194204646,100663296,-1204394608,100663296,-1194204647,100663296,-129639168,100663296,-1194204645\",5,1,44.181220,-117.135840,\"Construction on I-84 EB near MP 359, Drive with caution.\",,,I-84 E,500,0,0,1579419488000",
			"TTI-f47b8dba-59a3-372d-9cec-549eb252e2d5-TTR46312939215361-1,\"100663296,-1204020275,100663296,-1204020274,100663296,-916744017,100663296,-1204020245,100663296,-1194204646,100663296,-1204394608,100663296,-1194204647,100663296,-129639168,100663296,-1194204645\",MISCELLANEOUS,CRITICAL,44.181220,-117.135840,\"Construction on I-84 EB near MP 359, Drive with caution.\",,,I-84 E,500,0,false,1579419488000",
		},
	}

	for _, c := range cases {
		s := c.incident.CSVString()
		if s != c.s {
			t.Errorf("incident: %v, expect csv string %s, but got %s", c.incident, c.s, s)
		}
	}

}
