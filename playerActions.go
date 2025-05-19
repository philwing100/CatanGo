// playerActions.go
package main

func getValidRoadPlacements(game *CatanGame, player *Player) []*Edge {
	valid := []*Edge{}

	for _, edge := range game.Board.Graph.Edges {
		if edge.OccupiedBy != nil {
			continue // already taken
		}

		v1, v2 := edge.Vertices[0], edge.Vertices[1]
		if v1 == nil || v2 == nil {
			continue // malformed
		}

		// Case 1: adjacent to player's settlement/city
		if (v1.OccupiedBy == player && v1.Building != "") || (v2.OccupiedBy == player && v2.Building != "") {
			valid = append(valid, edge)
			continue
		}

		// Case 2: adjacent to player's road via connected edges
		if isConnectedToPlayerRoad(v1, player) || isConnectedToPlayerRoad(v2, player) {
			valid = append(valid, edge)
		}
	}

	return valid
}

func isConnectedToPlayerRoad(v *Vertex, player *Player) bool {
	for _, e := range v.Edges {
		if e != nil && e.OccupiedBy == player {
			return true
		}
	}
	return false
}

func getStartingSettlementLocations(game *CatanGame, player *Player) []*Vertex {
	valid := []*Vertex{}
	for _, v := range game.Board.Graph.Vertices {
		if v.OccupiedBy == nil {
			valid = append(valid, v)
		}
	}
	return valid
}
