// baseGame.go
package gameplay

import (
	"fmt"
	"strconv"
)

type BaseGame struct{}

func (bg *BaseGame) Initialize(playerNum int) *CatanGame {
	ids := make([]int, playerNum)
	for i := range ids {
		ids[i] = i + 1
	}
	game := NewCatanGame(ids)
	bg.Start(game)
	return game
}

func (bg *BaseGame) Start(game *CatanGame) {
	game.Phase = "main"
}

type BasePlayerSelector struct{}

type DiceRollFunc func(player *Player) int

func (bps *BasePlayerSelector) SelectStartingPlayer(game *CatanGame, roll DiceRollFunc) *Player {
	for len(game.Players) > 1 {
		currentRolls := make(map[int]int)

		for i, player := range game.Players {
			currentRolls[i] = roll(player)
		}

		maxRoll := -1
		for _, roll := range currentRolls {
			if roll > maxRoll {
				maxRoll = roll
			}
		}

		var remainingPlayers []*Player
		for i, roll := range currentRolls {
			if roll == maxRoll {
				remainingPlayers = append(remainingPlayers, game.Players[i])
			}
		}

		game.Players = remainingPlayers
	}

	return game.Players[0]
}

func GenerateSnakeOrder(startIndex, totalPlayers int) []int {
	var order []int

	// First pass: startIndex to end
	for i := startIndex; i < totalPlayers; i++ {
		order = append(order, i)
	}
	// Second pass: end to startIndex
	for i := totalPlayers - 1; i >= 0; i-- {
		order = append(order, i)
	}

	return order
}

func ComputeValidEdgePlacements(game *CatanGame, vertexID int) {
	//get the list of adjacent vertices

	//go through the list of edges and if an adjacent vertex is
}

func ValidateAndPlaceSettlement(vertexID int, player *Player, game *CatanGame) bool {
	vertex := game.Board.Graph.Vertices[strconv.Itoa(vertexID)]
	if vertex.OccupiedBy != nil {
		fmt.Println("Vertex already occupied by another player.")
		return false
	}

	if player.Resources["settlements"] <= 0 {
		fmt.Println("Player does not have enough settlements to place.")
		return false
	}

	// Check if the vertex is adjacent to any of the player's existing roads or settlements
	for _, adjID := range vertex.AdjacentVertexes {
		adjVertex := game.Board.Graph.Vertices[strconv.Itoa(adjID)]
		if adjVertex.OccupiedBy == player || (adjVertex.OccupiedBy != nil && adjVertex.Building == 1) {
			PlaceSettlement(vertexID, player, game)
			return true
		}
	}

	fmt.Println("Settlement cannot be placed here; no adjacent roads or settlements.")
	return false
}

func PlaceSettlement(vertexID int, player *Player, game *CatanGame) {
	vertex := game.Board.Graph.Vertices[strconv.Itoa(vertexID)]
	if vertex.OccupiedBy == nil {
		vertex.OccupiedBy = player
		player.Resources["settlements"]++
	} else {
		fmt.Print("WRONG")
		return
	}
}

func ValidateAndPlaceRoad() {

}

func PlaceRoad(vertexID1, vertexID2 int, player *Player, game *CatanGame) {
	edgeKey := fmt.Sprintf("%d-%d", vertexID1, vertexID2)
	edge := game.Board.Graph.Edges[edgeKey]

	if edge == nil {
		fmt.Println("INVALID EDGE")
		return
	}

	if edge.OccupiedBy == nil {
		edge.OccupiedBy = player
		player.Resources["roads"]++
	} else {
		fmt.Println("EDGE ALREADY OCCUPIED")
		return
	}
}
