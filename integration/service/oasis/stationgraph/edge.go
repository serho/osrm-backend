package stationgraph

type edgeID struct {
	fromNodeID nodeID
	toNodeID   nodeID
}

type edge struct {
	distance float64
	duration float64
}
