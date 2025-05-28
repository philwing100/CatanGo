// game.go

// Contains all of the structs defining the game state, board, players, etc.
package gameplay

import (
	"math/rand"
	"time"
)

type Player struct {
	ID               int
	Resources        map[string]int
	VictoryPoints    int
	DevelopmentCards map[string]int
}

type Tile struct {
	ID          int
	Resource    string
	NumberToken int
}

type Vertex struct {
	ID               int
	OccupiedBy       *Player // nil if empty
	Building         int
	AdjacentVertexes [3]int // Adjacent vertices
	TileIds          [3]int
}

type Edge struct {
	ID         string
	OccupiedBy *Player // nil if empty
	Vertices   [2]*Vertex
}

type Graph struct {
	Vertices map[string]*Vertex // key is a unique coord string
	Edges    map[string]*Edge   // key is a unique coord string
}

type Port struct {
	GiveResource string
	VertexIDs    [2]int
}

type Board struct {
	Tiles          []*Tile
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
	Cli       bool
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
		Cli:       false,
	}
}

func GenerateBoard() *Board {
	Tokens := []int{2, 3, 3, 4, 4, 5, 5, 6, 6, 8, 8, 9, 9, 10, 10, 11, 11, 12}
	Resources := []string{"W", "W", "W", "W", "L", "L", "L", "L", "O", "O", "O", "B", "B", "B", "S", "S", "S", "S", "D"}
	Ports := []string{"A", "A", "A", "A", "B", "O", "S", "W", "L"} // A = 3:1, otherwise 2:1 ports
	shuffleSlice(Tokens)
	shuffleSlice(Resources)
	shuffleSlice(Ports)

	board := &Board{
		// Initialize as a slice of 19 tiles (nil initially)
		Tiles:          make([]*Tile, 19),
		RobberPosition: 0, //Will set to desert later
		Ports:          make([]Port, 0),
	}

	tokenIndex := 0
	for i := 0; i < 19; i++ {
		tile := &Tile{
			ID:       i,
			Resource: Resources[i],
		}

		if Resources[i] != "D" {
			tile.NumberToken = Tokens[tokenIndex]
			tokenIndex++
		} else {
			tile.NumberToken = 0     // Desert has no number
			board.RobberPosition = i // Robber starts on desert
		}

		board.Tiles[i] = tile
	}

	// Ports placement currently as q,r side
	//Needs to be redone as port vertex id's
	portPositions := []struct {
		VertexIDs [2]int
	}{
		{[2]int{1, 2}},
		{[2]int{4, 5}},
		{[2]int{15, 16}},
		{[2]int{27, 38}},
		{[2]int{46, 47}},
		{[2]int{51, 52}},
		{[2]int{48, 49}},
		{[2]int{29, 39}},
		{[2]int{8, 18}},
	}

	board.Graph = GenerateGraphFromHardcodedData()

	for i, pos := range portPositions {
		port := Port{
			GiveResource: Ports[i],
			VertexIDs:    pos.VertexIDs,
		}
		board.Ports = append(board.Ports, port)
	}

	return board
}

func GenerateGraphFromHardcodedData() *Graph {

	// Hardcoded graph data for the Catan board according to image
	vertices := map[string]*Vertex{
		"1":  {ID: 1, AdjacentVertexes: [3]int{9, 2}, TileIds: [3]int{1}},
		"2":  {ID: 2, AdjacentVertexes: [3]int{1, 3}, TileIds: [3]int{1}},
		"3":  {ID: 3, AdjacentVertexes: [3]int{2, 4, 11}, TileIds: [3]int{1, 2}},
		"4":  {ID: 4, AdjacentVertexes: [3]int{3, 5}, TileIds: [3]int{2}},
		"5":  {ID: 5, AdjacentVertexes: [3]int{4, 6, 13}, TileIds: [3]int{2, 3}},
		"6":  {ID: 6, AdjacentVertexes: [3]int{5, 7}, TileIds: [3]int{3}},
		"7":  {ID: 7, AdjacentVertexes: [3]int{6, 15}, TileIds: [3]int{3}},
		"8":  {ID: 8, AdjacentVertexes: [3]int{9, 18}, TileIds: [3]int{4}},
		"9":  {ID: 9, AdjacentVertexes: [3]int{1, 8, 10}, TileIds: [3]int{1, 4}},
		"10": {ID: 10, AdjacentVertexes: [3]int{9, 11, 20}, TileIds: [3]int{1, 4, 5}},
		"11": {ID: 11, AdjacentVertexes: [3]int{3, 10, 12}, TileIds: [3]int{1, 2, 5}},
		"12": {ID: 12, AdjacentVertexes: [3]int{11, 13, 22}, TileIds: [3]int{2, 5, 6}},
		"13": {ID: 13, AdjacentVertexes: [3]int{5, 12, 14}, TileIds: [3]int{2, 3, 6}},
		"14": {ID: 14, AdjacentVertexes: [3]int{13, 15, 24}, TileIds: [3]int{3, 6, 7}},
		"15": {ID: 15, AdjacentVertexes: [3]int{7, 14, 16}, TileIds: [3]int{3, 7}},
		"16": {ID: 16, AdjacentVertexes: [3]int{15, 26}, TileIds: [3]int{7}},
		"17": {ID: 17, AdjacentVertexes: [3]int{18, 28}, TileIds: [3]int{8}},
		"18": {ID: 18, AdjacentVertexes: [3]int{8, 17, 19}, TileIds: [3]int{4, 8}},
		"19": {ID: 19, AdjacentVertexes: [3]int{18, 20, 30}, TileIds: [3]int{4, 8, 9}},
		"20": {ID: 20, AdjacentVertexes: [3]int{10, 19, 21}, TileIds: [3]int{4, 5, 9}},
		"21": {ID: 21, AdjacentVertexes: [3]int{20, 22, 32}, TileIds: [3]int{5, 9, 10}},
		"22": {ID: 22, AdjacentVertexes: [3]int{12, 21, 23}, TileIds: [3]int{5, 6, 10}},
		"23": {ID: 23, AdjacentVertexes: [3]int{22, 24, 34}, TileIds: [3]int{6, 10, 11}},
		"24": {ID: 24, AdjacentVertexes: [3]int{14, 23, 25}, TileIds: [3]int{6, 7, 11}},
		"25": {ID: 25, AdjacentVertexes: [3]int{24, 26, 36}, TileIds: [3]int{7, 11, 12}},
		"26": {ID: 26, AdjacentVertexes: [3]int{16, 25, 27}, TileIds: [3]int{7, 12}},
		"27": {ID: 27, AdjacentVertexes: [3]int{26, 38}, TileIds: [3]int{12}},
		"28": {ID: 28, AdjacentVertexes: [3]int{17, 29}, TileIds: [3]int{8}},
		"29": {ID: 29, AdjacentVertexes: [3]int{28, 30, 39}, TileIds: [3]int{8, 13}},
		"30": {ID: 30, AdjacentVertexes: [3]int{19, 29, 31}, TileIds: [3]int{8, 9, 13}},
		"31": {ID: 31, AdjacentVertexes: [3]int{30, 32, 41}, TileIds: [3]int{9, 13, 14}},
		"32": {ID: 32, AdjacentVertexes: [3]int{21, 31, 33}, TileIds: [3]int{9, 10, 14}},
		"33": {ID: 33, AdjacentVertexes: [3]int{32, 34, 43}, TileIds: [3]int{10, 14, 15}},
		"34": {ID: 34, AdjacentVertexes: [3]int{23, 33, 35}, TileIds: [3]int{10, 11, 15}},
		"35": {ID: 35, AdjacentVertexes: [3]int{34, 36, 45}, TileIds: [3]int{11, 15, 16}},
		"36": {ID: 36, AdjacentVertexes: [3]int{25, 35, 37}, TileIds: [3]int{11, 12, 16}},
		"37": {ID: 37, AdjacentVertexes: [3]int{36, 38, 47}, TileIds: [3]int{12, 16}},
		"38": {ID: 38, AdjacentVertexes: [3]int{27, 37}, TileIds: [3]int{12}},
		"39": {ID: 39, AdjacentVertexes: [3]int{29, 40}, TileIds: [3]int{13}},
		"40": {ID: 40, AdjacentVertexes: [3]int{39, 41, 48}, TileIds: [3]int{13, 17}},
		"41": {ID: 41, AdjacentVertexes: [3]int{31, 40, 42}, TileIds: [3]int{13, 14, 17}},
		"42": {ID: 42, AdjacentVertexes: [3]int{41, 43, 50}, TileIds: [3]int{14, 17, 18}},
		"43": {ID: 43, AdjacentVertexes: [3]int{33, 42, 44}, TileIds: [3]int{14, 15, 18}},
		"44": {ID: 44, AdjacentVertexes: [3]int{43, 45, 52}, TileIds: [3]int{15, 18, 19}},
		"45": {ID: 45, AdjacentVertexes: [3]int{35, 44, 46}, TileIds: [3]int{15, 16, 19}},
		"46": {ID: 46, AdjacentVertexes: [3]int{45, 47, 54}, TileIds: [3]int{16, 19}},
		"47": {ID: 47, AdjacentVertexes: [3]int{37, 46}, TileIds: [3]int{16}},
		"48": {ID: 48, AdjacentVertexes: [3]int{40, 49}, TileIds: [3]int{17}},
		"49": {ID: 49, AdjacentVertexes: [3]int{48, 50}, TileIds: [3]int{17}},
		"50": {ID: 50, AdjacentVertexes: [3]int{42, 49, 51}, TileIds: [3]int{17, 18}},
		"51": {ID: 51, AdjacentVertexes: [3]int{50, 52}, TileIds: [3]int{18}},
		"52": {ID: 52, AdjacentVertexes: [3]int{44, 51, 53}, TileIds: [3]int{18, 19}},
		"53": {ID: 53, AdjacentVertexes: [3]int{52, 54}, TileIds: [3]int{19}},
		"54": {ID: 54, AdjacentVertexes: [3]int{46, 53}, TileIds: [3]int{19}},
	}

	edges := map[string]*Edge{}

	return &Graph{
		Vertices: vertices,
		Edges:    edges,
	}
}

func GenerateBank() *Bank {
	bank := &Bank{
		Resources: map[string]int{
			"B": 19, // Brick
			"L": 19, // Lumber
			"S": 19, // Sheep
			"W": 19, // Wheat
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

func determineToken(i int, resource string, tokens []int) int {
	if resource == "D" {
		return 0 // Desert gets no number token
	}
	if i >= len(tokens) {
		return 0 // Safety check
	}
	return tokens[i]
}
