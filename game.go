// game.go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Constants for board layout
const (
	BoardRadius = 2
	BoardOffset = BoardRadius       // To make coordinates positive
	BoardSize   = BoardRadius*2 + 1 // Size of the 2D array
)

type Player struct {
	ID               int
	Resources        map[string]int
	VictoryPoints    int
	DevelopmentCards map[string]int
}

type Tile struct {
	Resource    string
	NumberToken int
	q           int
	r           int
}

type Vertex struct {
	ID         string
	OccupiedBy *Player // nil if empty
	Building   string  // "settlement", "city", or ""
	Edges      []*Edge
	Tiles      []*Tile
}

type Edge struct {
	ID         string
	OccupiedBy *Player // nil if empty
	Vertices   [2]*Vertex
	Tiles      []*Tile
}

type Graph struct {
	Vertices map[string]*Vertex // key is a unique coord string
	Edges    map[string]*Edge   // key is a unique coord string
}

type Port struct {
	GiveResource string
	HexSide      int
	q            int
	r            int
}

type Board struct {
	Tiles          [][]*Tile
	RobberPosition int
	Ports          []Port
	Graph          *Graph
}

type CatanGame struct {
	Players   []*Player
	Board     *Board
	TurnIndex int
	Phase     string
	Bank      *Bank
}

type DevelopmentCard struct {
	Type string
}

type Bank struct {
	Resources        map[string]int
	DevelopmentCards []DevelopmentCard
}

func NewCatanGame(playerIDs []int) *CatanGame {
	players := make([]*Player, 0)
	for _, id := range playerIDs {
		players = append(players, &Player{
			ID:               id,
			Resources:        make(map[string]int),
			VictoryPoints:    0,
			DevelopmentCards: make(map[string]int),
		})
	}

	board := GenerateBoard()

	return &CatanGame{
		Players:   players,
		Board:     board,
		TurnIndex: 0,
		Phase:     "setup",
		Bank:      GenerateBank(),
	}
}

func GenerateBoard() *Board {
	NumberTokens := []int{2, 3, 3, 4, 4, 5, 5, 6, 6, 8, 8, 9, 9, 10, 10, 11, 11, 12}
	Resources := []string{"W", "W", "W", "W", "L", "L", "L", "L", "O", "O", "O", "B", "B", "B", "S", "S", "S", "S", "D"}
	Ports := []string{"A", "A", "A", "A", "B", "O", "S", "W", "L"} // A = 3:1, otherwise 2:1 ports
	shuffleSlice(NumberTokens)
	shuffleSlice(Resources)
	shuffleSlice(Ports)

	// Initialize 2D array for tiles
	tiles := make([][]*Tile, BoardSize)
	for i := range tiles {
		tiles[i] = make([]*Tile, BoardSize)
	}

	board := &Board{
		Tiles:          tiles,
		RobberPosition: 0,
		Ports:          make([]Port, 0),
	}

	board.Graph = GenerateGraph(board)

	idx := 0
	desertFlag := false

	// Fill tiles
	for q := -BoardRadius; q <= BoardRadius; q++ {
		for r := -BoardRadius; r <= BoardRadius; r++ {
			s := -q - r
			if s >= -BoardRadius && s <= BoardRadius {
				if idx < len(Resources) {
					currentTile := &Tile{
						q:        q,
						r:        r,
						Resource: Resources[idx],
					}
					if Resources[idx] == "D" {
						currentTile.NumberToken = 0 // Desert has no number token
						desertFlag = true
					} else if desertFlag {
						currentTile.NumberToken = NumberTokens[idx-1]
					} else {
						currentTile.NumberToken = NumberTokens[idx]
					}

					// Place tile into the 2D slice with offset
					board.Tiles[q+BoardOffset][r+BoardOffset] = currentTile

					idx++
				}
			}
		}
	}

	// Ports placement
	portPositions := []struct {
		q, r, side int
	}{
		{-2, 0, 5},
		{-2, 1, 5},
		{-1, 2, 6},
		{1, 2, 6},
		{2, 1, 1},
		{2, 0, 1},
		{1, -2, 2},
		{0, -2, 2},
		{-1, -1, 3},
	}

	for i, pos := range portPositions {
		port := Port{
			GiveResource: Ports[i],
			q:            pos.q,
			r:            pos.r,
			HexSide:      pos.side,
		}
		board.Ports = append(board.Ports, port)
	}

	board.Graph = &Graph{
		Vertices: make(map[string]*Vertex),
		Edges:    make(map[string]*Edge),
	}

	// Define directions
	//vertexDirs := [6][2]int{{0, 1}, {-1, 1}, {-1, 0}, {0, -1}, {1, -1}, {1, 0}}
	edgeDirs := [6][2]int{{1, 0}, {0, 1}, {-1, 1}, {-1, 0}, {0, -1}, {1, -1}}

	for q := -BoardRadius; q <= BoardRadius; q++ {
		for r := -BoardRadius; r <= BoardRadius; r++ {
			s := -q - r
			if s < -BoardRadius || s > BoardRadius {
				continue
			}
			tile := tiles[q+BoardOffset][r+BoardOffset]
			if tile == nil {
				continue
			}

			for i := 0; i < 6; i++ {
				// Vertex key
				vKey := fmt.Sprintf("v-%d-%d-%d", q, r, i)
				if _, exists := board.Graph.Vertices[vKey]; !exists {
					board.Graph.Vertices[vKey] = &Vertex{
						ID:    vKey,
						Tiles: []*Tile{},
						Edges: []*Edge{},
					}
				}
				board.Graph.Vertices[vKey].Tiles = append(board.Graph.Vertices[vKey].Tiles, tile)

				// Edge key
				eKey := fmt.Sprintf("e-%d-%d-%d", q, r, i)
				if _, exists := board.Graph.Edges[eKey]; !exists {
					board.Graph.Edges[eKey] = &Edge{
						ID:    eKey,
						Tiles: []*Tile{},
						Vertices: [2]*Vertex{
							board.Graph.Vertices[fmt.Sprintf("v-%d-%d-%d", q, r, i)],
							board.Graph.Vertices[fmt.Sprintf("v-%d-%d-%d", q+edgeDirs[i][0], r+edgeDirs[i][1], (i+3)%6)],
						},
					}
				}
				board.Graph.Edges[eKey].Tiles = append(board.Graph.Edges[eKey].Tiles, tile)

				// Link vertex to edge
				board.Graph.Vertices[vKey].Edges = append(board.Graph.Vertices[vKey].Edges, board.Graph.Edges[eKey])
			}
		}
	}

	return board
}

func GenerateGraph(board *Board) *Graph {
	graph := &Graph{
		Vertices: make(map[string]*Vertex),
		Edges:    make(map[string]*Edge),
	}

	vertexID := 0
	edgeID := 0

	// Directions for hex corners and edges (0 to 5)
	vertexDirs := []int{0, 1, 2, 3, 4, 5}
	edgeDirs := []int{0, 1, 2, 3, 4, 5}

	for q := -BoardRadius; q <= BoardRadius; q++ {
		for r := -BoardRadius; r <= BoardRadius; r++ {
			s := -q - r
			if s < -BoardRadius || s > BoardRadius {
				continue
			}

			tile := board.Tiles[q+BoardOffset][r+BoardOffset]
			if tile == nil {
				continue
			}

			// Generate vertices for tile
			for _, dir := range vertexDirs {
				key := vertexKey(q, r, dir)
				if _, exists := graph.Vertices[key]; !exists {
					graph.Vertices[key] = &Vertex{
						ID:       key,
						Building: "",
						Tiles:    []*Tile{tile},
					}
					vertexID++
				} else {
					graph.Vertices[key].Tiles = append(graph.Vertices[key].Tiles, tile)
				}
			}

			// Generate edges for tile
			for _, dir := range edgeDirs {
				key := edgeKey(q, r, dir)
				if _, exists := graph.Edges[key]; !exists {
					graph.Edges[key] = &Edge{
						ID:    key,
						Tiles: []*Tile{tile},
					}
					edgeID++
				} else {
					graph.Edges[key].Tiles = append(graph.Edges[key].Tiles, tile)
				}
			}
		}
	}

	// Optional: link vertices <-> edges here

	return graph
}

func vertexKey(q, r, corner int) string {
	// Normalize corner IDs to avoid duplicates from different tiles
	return fmt.Sprintf("v-%d-%d-%d", q, r, corner)
}

func edgeKey(q, r, side int) string {
	// Normalize edge IDs similarly
	return fmt.Sprintf("e-%d-%d-%d", q, r, side)
}

func GenerateBank() *Bank {
	bank := &Bank{
		Resources: map[string]int{
			"B": 19, // Brick
			"L": 19, // Lumber
			"W": 19, // Wool
			"S": 19, // Grain
			"O": 19, // Ore
		},
		DevelopmentCards: make([]DevelopmentCard, 0, 25),
	}

	// Add development cards to the slice
	cardCounts := map[string]int{
		"Knight":         14,
		"Victory Point":  5,
		"Road Building":  2,
		"Year of Plenty": 2,
		"Monopoly":       2,
	}

	for cardType, count := range cardCounts {
		for i := 0; i < count; i++ {
			bank.DevelopmentCards = append(bank.DevelopmentCards, DevelopmentCard{Type: cardType})
		}
	}

	shuffleDevCards(bank.DevelopmentCards)

	return bank
}

func shuffleDevCards(cards []DevelopmentCard) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func createPlayer(id int) Player {
	return Player{ID: id, Resources: make(map[string]int), VictoryPoints: 0, DevelopmentCards: make(map[string]int)}
}

func shuffleSlice[T any](slice []T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func (game *CatanGame) PlaceStartingStructures(player *Player, vertexKey string, edgeKey string) error {
	vertex, ok := game.Board.Graph.Vertices[vertexKey]
	if !ok {
		return fmt.Errorf("invalid vertex key: %s", vertexKey)
	}
	if vertex.OccupiedBy != nil {
		return fmt.Errorf("vertex %s already has a settlement or city", vertexKey)
	}
	// Check distance rule: no adjacent settlements
	for _, edge := range vertex.Edges {
		for _, neighbor := range edge.Vertices {
			if neighbor != vertex && neighbor.OccupiedBy != nil {
				return fmt.Errorf("cannot place settlement: adjacent vertex %d is occupied", neighbor.ID)
			}
		}
	}

	edge, ok := game.Board.Graph.Edges[edgeKey]
	if !ok {
		return fmt.Errorf("invalid edge key: %s", edgeKey)
	}
	if edge.OccupiedBy != nil {
		return fmt.Errorf("edge %s already has a road", edgeKey)
	}
	// Check that road connects to the vertex
	if edge.Vertices[0] != vertex && edge.Vertices[1] != vertex {
		return fmt.Errorf("edge %s does not connect to vertex %s", edgeKey, vertexKey)
	}

	// Place settlement and road
	vertex.OccupiedBy = player
	vertex.Building = "settlement"
	edge.OccupiedBy = player

	player.VictoryPoints += 1

	return nil
}
