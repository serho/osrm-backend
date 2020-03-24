// package spatialindexer answers query of nearest pois for given point
//
// Sample 1: Build connectivity for charge stations during pro-processing
// indexer := NewS2Indexer().Build(poiCsvFile)
// for _, stationPoint := range chargeStations {
// 	nearbyStations := indexer.FindNearByIDs(stationPoint, 800km, -1)
// 	rankedStations := indexer.RankingIDsByShortestDistance(stationPoint, nearbyStations)
// }
//
//
// Sample 2: Dump s2Indexer's content to folder
// indexer.Dump(folderPath)
//
//
// Sample 3: Query reachable charge stations with current energy level
// indexer := NewS2Indexer().Load(folderPath)
// nearbyStations := indexer.FindNearByIDs(currentPoint, currentEnergyLevel, -1)
//
package spatialindexer
