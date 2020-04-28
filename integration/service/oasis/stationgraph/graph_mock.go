package stationgraph

import "github.com/serho/osrm-backend/integration/service/oasis/chargingstrategy"

// node_0 -> node_1, duration = 30, distance = 30
// node_0 -> node_2, duration = 20, distance = 20
// node_1 -> node_3, duration = 10, distance = 10
// node_2 -> node_4, duration = 50, distance = 50
// node_2 -> node_3, duration = 50, distance = 50
// node_3 -> node_4, duration = 10, distance = 10
// Set charge information to fixed status to ignore situation of lack of energy
func NewMockGraph1() Graph {
	return &mockGraph{
		[]*node{
			{
				0,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 0.0,
					lon: 0.0,
				},
			},
			{
				1,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 1.1,
					lon: 1.1,
				},
			},
			{
				2,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 2.2,
					lon: 2.2,
				},
			},
			{
				3,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 3.3,
					lon: 3.3,
				},
			},
			{
				4,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 4.4,
					lon: 4.4,
				},
			},
		},
		[]string{
			"node_0",
			"node_1",
			"node_2",
			"node_3",
			"node_4",
		},
		map[nodeID][]*edgeIDAndData{
			// node_0 -> node_1, duration = 30, distance = 30
			// node_0 -> node_2, duration = 20, distance = 20
			0: {
				{edgeID{0, 1}, &edge{30, 30}},
				{edgeID{0, 2}, &edge{20, 20}},
			},
			// node_1 -> node_3, duration = 10, distance = 10
			1: {
				{edgeID{1, 3}, &edge{10, 10}},
			},
			// node_2 -> node_4, duration = 50, distance = 50
			// node_2 -> node_3, duration = 50, distance = 50
			2: {
				{edgeID{2, 4}, &edge{50, 50}},
				{edgeID{2, 3}, &edge{50, 50}},
			},
			// node_3 -> node_4, duration = 10, distance = 10
			3: {
				{edgeID{3, 4}, &edge{10, 10}},
			},
		},
		chargingstrategy.NewNullChargeStrategy(),
	}
}

// node_0 -> node_1, duration = 30, distance = 30
// node_0 -> node_2, duration = 20, distance = 20
// node_1 -> node_3, duration = 20, distance = 20
// node_1 -> node_4, duration = 15, distance = 15
// node_2 -> node_3, duration = 30, distance = 30
// node_2 -> node_4, duration = 20, distance = 20
// node_3 -> node_5, duration = 10, distance = 10
// node_3 -> node_6, duration = 10, distance = 10
// node_3 -> node_7, duration = 10, distance = 10
// node_4 -> node_5, duration = 15, distance = 15
// node_4 -> node_6, duration = 15, distance = 15
// node_4 -> node_7, duration = 15, distance = 15
// node_5 -> node_8, duration = 10, distance = 10
// node_6 -> node_8, duration = 20, distance = 20
// node_7 -> node_8, duration = 30, distance = 30
// Set charge information to fixed status to ignore situation of lack of energy
func NewMockGraph2() Graph {
	return &mockGraph{
		[]*node{
			{
				0,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 0.0,
					lon: 0.0,
				},
			},
			{
				1,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 1.1,
					lon: 1.1,
				},
			},
			{
				2,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 2.2,
					lon: 2.2,
				},
			},
			{
				3,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 3.3,
					lon: 3.3,
				},
			},
			{
				4,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 4.4,
					lon: 4.4,
				},
			},
			{
				5,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 5.5,
					lon: 5.5,
				},
			},
			{
				6,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 6.6,
					lon: 6.6,
				},
			},
			{
				7,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 7.7,
					lon: 7.7,
				},
			},
			{
				8,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 8.8,
					lon: 8.8,
				},
			},
		},
		[]string{
			"node_0",
			"node_1",
			"node_2",
			"node_3",
			"node_4",
			"node_5",
			"node_6",
			"node_7",
			"node_8",
		},
		map[nodeID][]*edgeIDAndData{
			// node_0 -> node_1, duration = 30, distance = 30
			// node_0 -> node_2, duration = 20, distance = 20
			0: {
				{
					edgeID{
						0,
						1,
					},
					&edge{
						30,
						30,
					},
				},
				{
					edgeID{
						0,
						2,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_1 -> node_3, duration = 20, distance = 20
			// node_1 -> node_4, duration = 15, distance = 15
			1: {
				{
					edgeID{
						1,
						3,
					},
					&edge{
						20,
						20,
					},
				},
				{
					edgeID{
						1,
						4,
					},
					&edge{
						15,
						15,
					},
				},
			},
			// node_2 -> node_3, duration = 30, distance = 30
			// node_2 -> node_4, duration = 20, distance = 20
			2: {
				{
					edgeID{
						2,
						3,
					},
					&edge{
						30,
						30,
					},
				},
				{
					edgeID{
						2,
						4,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_3 -> node_5, duration = 10, distance = 10
			// node_3 -> node_6, duration = 10, distance = 10
			// node_3 -> node_7, duration = 10, distance = 10
			3: {
				{
					edgeID{
						3,
						5,
					},
					&edge{
						10,
						10,
					},
				},
				{
					edgeID{
						3,
						6,
					},
					&edge{
						10,
						10,
					},
				},
				{
					edgeID{
						3,
						7,
					},
					&edge{
						10,
						10,
					},
				},
			},
			// node_4 -> node_5, duration = 15, distance = 15
			// node_4 -> node_6, duration = 15, distance = 15
			// node_4 -> node_7, duration = 15, distance = 15
			4: {
				{
					edgeID{
						4,
						5,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						4,
						6,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						4,
						7,
					},
					&edge{
						15,
						15,
					},
				},
			},
			// node_5 -> node_8, duration = 10, distance = 10
			5: {
				{
					edgeID{
						5,
						8,
					},
					&edge{
						10,
						10,
					},
				},
			},
			// node_6 -> node_8, duration = 20, distance = 20
			6: {
				{
					edgeID{
						6,
						8,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_7 -> node_8, duration = 30, distance = 30
			7: {
				{
					edgeID{
						7,
						8,
					},
					&edge{
						30,
						30,
					},
				},
			},
		},
		chargingstrategy.NewNullChargeStrategy(),
	}
}

// node_0 -> node_1, duration = 15, distance = 15
// node_0 -> node_2, duration = 20, distance = 20
// node_1 -> node_3, duration = 20, distance = 20
// node_1 -> node_4, duration = 15, distance = 15
// node_2 -> node_3, duration = 30, distance = 30
// node_2 -> node_4, duration = 20, distance = 20
// node_3 -> node_5, duration = 10, distance = 10
// node_3 -> node_6, duration = 10, distance = 10
// node_3 -> node_7, duration = 10, distance = 10
// node_4 -> node_5, duration = 15, distance = 15
// node_4 -> node_6, duration = 15, distance = 15
// node_4 -> node_7, duration = 15, distance = 15
// node_5 -> node_8, duration = 10, distance = 10
// node_6 -> node_8, duration = 20, distance = 20
// node_7 -> node_8, duration = 30, distance = 30
// Set charge information to fixed status to ignore situation of lack of energy
func NewMockGraph3() Graph {
	return &mockGraph{
		[]*node{
			{
				0,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 0.0,
					lon: 0.0,
				},
			},
			{
				1,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 1.1,
					lon: 1.1,
				},
			},
			{
				2,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 2.2,
					lon: 2.2,
				},
			},
			{
				3,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 3.3,
					lon: 3.3,
				},
			},
			{
				4,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 4.4,
					lon: 4.4,
				},
			},
			{
				5,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 5.5,
					lon: 5.5,
				},
			},
			{
				6,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 6.6,
					lon: 6.6,
				},
			},
			{
				7,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 7.7,
					lon: 7.7,
				},
			},
			{
				8,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 999,
					},
				},
				locationInfo{
					lat: 8.8,
					lon: 8.8,
				},
			},
		},
		[]string{
			"node_0",
			"node_1",
			"node_2",
			"node_3",
			"node_4",
			"node_5",
			"node_6",
			"node_7",
			"node_8",
		},
		map[nodeID][]*edgeIDAndData{
			// node_0 -> node_1, duration = 15, distance = 15
			// node_0 -> node_2, duration = 20, distance = 20
			0: {
				{
					edgeID{
						0,
						1,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						0,
						2,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_1 -> node_3, duration = 20, distance = 20
			// node_1 -> node_4, duration = 15, distance = 15
			1: {
				{
					edgeID{
						1,
						3,
					},
					&edge{
						20,
						20,
					},
				},
				{
					edgeID{
						1,
						4,
					},
					&edge{
						15,
						15,
					},
				},
			},
			// node_2 -> node_3, duration = 30, distance = 30
			// node_2 -> node_4, duration = 20, distance = 20
			2: {
				{
					edgeID{
						2,
						3,
					},
					&edge{
						30,
						30,
					},
				},
				{
					edgeID{
						2,
						4,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_3 -> node_5, duration = 10, distance = 10
			// node_3 -> node_6, duration = 10, distance = 10
			// node_3 -> node_7, duration = 10, distance = 10
			3: {
				{
					edgeID{
						3,
						5,
					},
					&edge{
						10,
						10,
					},
				},
				{
					edgeID{
						3,
						6,
					},
					&edge{
						10,
						10,
					},
				},
				{
					edgeID{
						3,
						7,
					},
					&edge{
						10,
						10,
					},
				},
			},
			// node_4 -> node_5, duration = 15, distance = 15
			// node_4 -> node_6, duration = 15, distance = 15
			// node_4 -> node_7, duration = 15, distance = 15
			4: {
				{
					edgeID{
						4,
						5,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						4,
						6,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						4,
						7,
					},
					&edge{
						15,
						15,
					},
				},
			},
			// node_5 -> node_8, duration = 10, distance = 10
			5: {
				{
					edgeID{
						5,
						8,
					},
					&edge{
						10,
						10,
					},
				},
			},
			// node_6 -> node_8, duration = 20, distance = 20
			6: {
				{
					edgeID{
						6,
						8,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_7 -> node_8, duration = 30, distance = 30
			7: {
				{
					edgeID{
						7,
						8,
					},
					&edge{
						30,
						30,
					},
				},
			},
		},
		chargingstrategy.NewNullChargeStrategy(),
	}
}

// node_0 -> node_1, duration = 15, distance = 15
// node_0 -> node_2, duration = 20, distance = 20
// node_1 -> node_3, duration = 20, distance = 20
// node_1 -> node_4, duration = 15, distance = 15
// node_2 -> node_3, duration = 30, distance = 30
// node_2 -> node_4, duration = 5, distance = 5
// node_3 -> node_5, duration = 10, distance = 10
// node_3 -> node_6, duration = 10, distance = 10
// node_3 -> node_7, duration = 10, distance = 10
// node_4 -> node_5, duration = 15, distance = 15
// node_4 -> node_6, duration = 15, distance = 15
// node_4 -> node_7, duration = 15, distance = 15
// node_5 -> node_8, duration = 10, distance = 10
// node_6 -> node_8, duration = 20, distance = 20
// node_7 -> node_8, duration = 30, distance = 30
// Charge information
// each station only charges 16 unit of energy
func NewMockGraph4() Graph {
	return &mockGraph{
		[]*node{
			{
				0,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 0.0,
					lon: 0.0,
				},
			},
			{
				1,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 1.1,
					lon: 1.1,
				},
			},
			{
				2,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 2.2,
					lon: 2.2,
				},
			},
			{
				3,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 3.3,
					lon: 3.3,
				},
			},
			{
				4,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 4.4,
					lon: 4.4,
				},
			},
			{
				5,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 5.5,
					lon: 5.5,
				},
			},
			{
				6,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 6.6,
					lon: 6.6,
				},
			},
			{
				7,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 16,
					},
				},
				locationInfo{
					lat: 7.7,
					lon: 7.7,
				},
			},
			{
				8,
				chargeInfo{
					targetState: chargingstrategy.State{
						Energy: 0,
					},
				},
				locationInfo{
					lat: 8.8,
					lon: 8.8,
				},
			},
		},
		[]string{
			"node_0",
			"node_1",
			"node_2",
			"node_3",
			"node_4",
			"node_5",
			"node_6",
			"node_7",
			"node_8",
		},
		map[nodeID][]*edgeIDAndData{
			// node_0 -> node_1, duration = 15, distance = 15
			// node_0 -> node_2, duration = 20, distance = 20
			0: {
				{
					edgeID{
						0,
						1,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						0,
						2,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_1 -> node_3, duration = 20, distance = 20
			// node_1 -> node_4, duration = 15, distance = 15
			1: {
				{
					edgeID{
						1,
						3,
					},
					&edge{
						20,
						20,
					},
				},
				{
					edgeID{
						1,
						4,
					},
					&edge{
						15,
						15,
					},
				},
			},
			// node_2 -> node_3, duration = 30, distance = 30
			// node_2 -> node_4, duration = 5, distance = 5
			2: {
				{
					edgeID{
						2,
						3,
					},
					&edge{
						30,
						30,
					},
				},
				{
					edgeID{
						2,
						4,
					},
					&edge{
						5,
						5,
					},
				},
			},
			// node_3 -> node_5, duration = 10, distance = 10
			// node_3 -> node_6, duration = 10, distance = 10
			// node_3 -> node_7, duration = 10, distance = 10
			3: {
				{
					edgeID{
						3,
						5,
					},
					&edge{
						10,
						10,
					},
				},
				{
					edgeID{
						3,
						6,
					},
					&edge{
						10,
						10,
					},
				},
				{
					edgeID{
						3,
						7,
					},
					&edge{
						10,
						10,
					},
				},
			},
			// node_4 -> node_5, duration = 15, distance = 15
			// node_4 -> node_6, duration = 15, distance = 15
			// node_4 -> node_7, duration = 15, distance = 15
			4: {
				{
					edgeID{
						4,
						5,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						4,
						6,
					},
					&edge{
						15,
						15,
					},
				},
				{
					edgeID{
						4,
						7,
					},
					&edge{
						15,
						15,
					},
				},
			},
			// node_5 -> node_8, duration = 10, distance = 10
			5: {
				{
					edgeID{
						5,
						8,
					},
					&edge{
						10,
						10,
					},
				},
			},
			// node_6 -> node_8, duration = 20, distance = 20
			6: {
				{
					edgeID{
						6,
						8,
					},
					&edge{
						20,
						20,
					},
				},
			},
			// node_7 -> node_8, duration = 30, distance = 30
			7: {
				{
					edgeID{
						7,
						8,
					},
					&edge{
						30,
						30,
					},
				},
			},
		},
		chargingstrategy.NewNullChargeStrategy(),
	}
}

type mockGraph struct {
	nodes      []*node
	stationIDs []string
	edges      map[nodeID][]*edgeIDAndData
	strategy   chargingstrategy.Strategy
}

// Node returns node object by its nodeID
func (graph *mockGraph) Node(id nodeID) *node {
	if graph.isValidNodeID(id) {
		return graph.nodes[id]
	}
	return nil
}

// AdjacentNodes returns a group of node ids which connect with given node id
// The connectivity between nodes is build during running time.
func (graph *mockGraph) AdjacentNodes(id nodeID) []nodeID {
	var nodeIDs []nodeID
	if graph.isValidNodeID(id) {
		edges, ok := graph.edges[id]
		if ok {
			for _, edge := range edges {
				nodeIDs = append(nodeIDs, edge.edgeId.toNodeID)
			}
		}
	}

	return nodeIDs
}

// Edge returns edge information between given two nodes
func (graph *mockGraph) Edge(from, to nodeID) *edge {
	if graph.isValidNodeID(from) && graph.isValidNodeID(to) {
		edges, ok := graph.edges[from]
		if ok {
			for _, edge := range edges {
				if edge.edgeId.toNodeID == to {
					return edge.edgeData
				}
			}
		}
	}

	return nil
}

// SetStart generates start node for the graph
func (graph *mockGraph) SetStart(stationID string, targetState chargingstrategy.State, location locationInfo) Graph {
	return graph
}

// SetEnd generates end node for the graph
func (graph *mockGraph) SetEnd(stationID string, targetState chargingstrategy.State, location locationInfo) Graph {
	return graph
}

// StartNodeID returns start node's ID for given graph
func (graph *mockGraph) StartNodeID() nodeID {
	return invalidNodeID
}

// EndNodeID returns end node's ID for given graph
func (graph *mockGraph) EndNodeID() nodeID {
	return invalidNodeID
}

// ChargeStrategy returns charge strategy used for graph construction
func (graph *mockGraph) ChargeStrategy() chargingstrategy.Strategy {
	return graph.strategy
}

// StationID returns original stationID from internal nodeID
func (graph *mockGraph) StationID(id nodeID) string {
	if id < 0 || int(id) >= len(graph.stationIDs) {
		return invalidStationID
	}
	return graph.stationIDs[id]
}

func (graph *mockGraph) isValidNodeID(id nodeID) bool {
	if id < 0 || int(id) >= len(graph.nodes) {
		return false
	}
	return true
}
