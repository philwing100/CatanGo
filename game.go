package main

import (
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
	ID               string
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
}

type CatanGame struct {
	Players   []*Player
	Board     *Board
	TurnIndex int
	Phase     string
	Bank      map[string]int
}

func NewCatanGame(playerIDs []string) *CatanGame {
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

	return board
}

func shuffleSlice[T any](slice []T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
