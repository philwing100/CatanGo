// baseGame.go
package gameplay

import (
	"fmt"
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
	// Use a local copy to find the starting player
	contenders := game.Players

	for len(contenders) > 1 {
		currentRolls := make(map[int]int)

		for i, player := range contenders {
			currentRolls[i] = roll(player)
		}

		maxRoll := -1
		for _, r := range currentRolls {
			if r > maxRoll {
				maxRoll = r
			}
		}

		var nextRound []*Player
		for i, r := range currentRolls {
			if r == maxRoll {
				nextRound = append(nextRound, contenders[i])
			}
		}

		contenders = nextRound // Only update the local contenders slice
	}

	// Do not mutate game.Players
	return contenders[0]
}

func GenerateSnakeOrder(game *CatanGame, startingPlayer *Player, totalPlayers int) []int {
	var order []int //Build a stack to create the snake

	startIdx := -1
	for i, player := range game.Players {
		if player.ID == startingPlayer.ID {
			startIdx = i
			break
		}
	}
	if startIdx == -1 {
		// fallback: default to 0 if not found
		startIdx = 0
	} // Find the index of the starting player

	// Wraparound loop to create the first placement order
	for i := startIdx; i < totalPlayers; i++ {
		order = append(order, game.Players[i].ID)
	}
	for i := 0; i < startIdx; i++ {
		order = append(order, game.Players[i].ID)
	}

	for i := totalPlayers - 1; i >= 0; i-- {
		order = append(order, order[i])
	}
	fmt.Printf("%v", order)

	return order
}

// Get the list of valid edges for a given vertex regardless of who owns it
// Check one VertexID and then see what edges (max of 3) are available
func ComputeValidEdgePlacements(game *CatanGame, vertexID int) []int {
	var EdgeIDs []int
	//vertex := game.Board.Graph.Vertices[strconv.Itoa(vertexID)]

	return EdgeIDs
}

func GetAdjacentVertices(vertexID int, game *CatanGame) []int {
	var adjacentVertexIDs []int
	vertex, exists := game.Board.Graph.Vertices[vertexID]
	if !exists {
		return adjacentVertexIDs // Return empty if vertex does not exist
	}

	for _, adjID := range vertex.AdjacentVertexes {
		if adjVertex, exists := game.Board.Graph.Vertices[adjID]; exists {
			adjacentVertexIDs = append(adjacentVertexIDs, adjVertex.ID)
		}
	}

	return adjacentVertexIDs
}

// get the list of valid vertex placements for a player at the beginning of a game
// after the beginning of the game the player can only place settlements on vertices adjacent to their roads or settlements
// input the game
// Also needs to verify that the vertex is at least two spaces away from another player's settlement
func ComputeValidVertexPlacements(game *CatanGame) []int {
	var VertexIDs []int

	for _, vertex := range game.Board.Graph.Vertices {
		if vertex.OccupiedBy != nil || vertex.Building != 0 {
			continue // Skip if already occupied or has a building
		}

		valid := true
		for _, adjID := range vertex.AdjacentVertexes {
			adjVertex, exists := game.Board.Graph.Vertices[adjID]
			if !exists {
				continue // skip if adjacent vertex doesn't exist (defensive check)
			}
			if adjVertex.OccupiedBy != nil {
				valid = false
				break
			}
		}

		if valid {
			VertexIDs = append(VertexIDs, vertex.ID)
		}
	}

	return VertexIDs
}

// Checks if vertex is empty, player can afford it, and if it is adjacent to a road
func ValidateAndPlaceSettlement(vertexID int, player *Player, game *CatanGame) bool {

	return false
}

// Assume validation has already been done
func PlaceSettlement(vertexID int, player *Player, game *CatanGame) {
	vertex := game.Board.Graph.Vertices[vertexID]
	if vertex.OccupiedBy == nil {
		vertex.OccupiedBy = player
		vertex.Building = 1       // 1 for settlement
		player.VictoryPoints += 1 // Increment player's victory points
	}
}

// Checks if the place the player wants to build a road is empty
func RoadEmptySpace() {

}

// Validates that the player can place a road,
func ValidateAndPlaceRoad() {

}

// Assumes they can afford it, used solo when road building card or snake build is used
func PlaceRoad(vertexID1, vertexID2 int, player *Player, game *CatanGame) {
	edgeKey := fmt.Sprintf("%d-%d", vertexID1, vertexID2)

	var Road = Edge{
		ID:         edgeKey,
		OccupiedBy: player,
		Vertices:   [2]*Vertex{game.Board.Graph.Vertices[vertexID1], game.Board.Graph.Vertices[vertexID2]},
	}

	game.Board.Graph.Edges[edgeKey] = &Road
}

// Pass in the player and what the player wants to buy
// Returns true if the player can afford the resource, false otherwise
func CanPlayerAfford() {

}

func BankToPlayerResource(game *CatanGame, player *Player, resource string, amount int) bool {
	if game.Bank.Resources[resource] >= amount {
		game.Bank.Resources[resource] -= amount
		player.Resources[resource] += amount
		return true
	}
	return false // Not enough resources in the bank
}

func PlayerToBankResource() {

}

func GetTileByID(game *CatanGame, tileID int) *Tile {
	for _, tile := range game.Board.Tiles {
		if tile.ID == tileID {
			return tile
		}
	}
	return nil // Return nil if no tile found with the given ID
}

func GetVertexByID(game *CatanGame, vertexID int) *Vertex {
	vertex, exists := game.Board.Graph.Vertices[vertexID]
	if !exists {
		return nil // Return nil if vertex does not exist
	}
	return vertex
}
