// cli.go
package gameplay

import (
	"fmt"
	"strings"
)

var playerColors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[34m", // Blue
	"\033[33m", // Yellow
}

const resetColor = "\033[0m"

func colorText(text string, playerID int) string {
	if playerID <= 0 || playerID > len(playerColors) {
		return text
	}
	return playerColors[playerID-1] + text + resetColor
}

// reads int until a valid int is entered
func readInt(prompt string) int {
	var input int
	for {
		fmt.Print(prompt)
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			// Clear input buffer
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		return input
	}
}

// Print the raw contents of the board
func PrintRaw(game *CatanGame) {
	fmt.Println("===== CATAN GAME STATE =====\n")

	// Print players and their stats
	fmt.Println("Players:")
	for _, player := range game.Players {
		fmt.Printf("Player %d:\n", player.ID)
		fmt.Printf("  Resources: %v\n", player.Resources)
		fmt.Printf("  Victory Points: %d\n", player.VictoryPoints)
		fmt.Printf("  Development Cards: %v\n", player.DevelopmentCards)
	}
	fmt.Println()

	// Print vertices ownership
	fmt.Println("Vertices:")
	for key, vertex := range game.Board.Graph.Vertices {
		if vertex.OccupiedBy != nil {
			fmt.Printf("Vertex ID %s (Struct ID %d): Occupied by Player %d with building type %d\n",
				key, vertex.ID, vertex.OccupiedBy.ID, vertex.Building)
		}
	}
	fmt.Println()

	// Print edges ownership
	fmt.Println("Edges:")
	for key, edge := range game.Board.Graph.Edges {
		if edge.OccupiedBy != nil {
			fmt.Printf("Edge ID %s: Occupied by Player %d (Vertices %d and %d)\n",
				key, edge.OccupiedBy.ID, edge.Vertices[0].ID, edge.Vertices[1].ID)
		}
	}
	fmt.Println()

	// Print bank status
	fmt.Println("Bank:")
	fmt.Printf("  Resources: %v\n", game.Bank.Resources)
	fmt.Printf("  Remaining Dev Cards: %d\n", len(game.Bank.DevelopmentCards))
	fmt.Println()

	// Print turn and phase
	fmt.Println("Game Info:")
	fmt.Printf("  Current Turn: Player %d\n", game.Players[game.TurnIndex].ID)
	fmt.Printf("  Phase: %s\n", game.Phase)
	fmt.Println("============================\n")
}

func PrintGameBoard(game *CatanGame) {
	// Tile layout by rows
	tileRows := [][]int{
		{1, 2, 3},
		{4, 5, 6, 7},
		{8, 9, 10, 11, 12},
		{13, 14, 15, 16},
		{17, 18, 19},
	}

	// Vertex layout by rows
	vertexRows := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{8, 9, 10, 11, 12, 13, 14, 15, 16},
		{17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27},
		{28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38},
		{39, 40, 41, 42, 43, 44, 45, 46, 47},
		{48, 49, 50, 51, 52, 53, 54},
	}

	// Build port information
	portInfo := make(map[int]string)
	for _, port := range game.Board.Ports {
		for _, vertexID := range port.VertexIDs {
			portInfo[vertexID] = port.GiveResource
		}
	}

	// Print tile rows with vertices between them
	for i, tileRow := range tileRows {
		// Print vertex row above tiles (except first row)
		if i > 0 {
			printVertexRow(game, vertexRows[i-1], portInfo, i-1)
		}

		// Print tile row
		printTileRow(game, tileRow)

		// Print vertex row below tiles (except last row)
		if i < len(tileRows)-1 {
			printVertexRow(game, vertexRows[i+1], portInfo, i+1)
		}
	}
}

func printTileRow(game *CatanGame, tileIDs []int) {
	line1 := ""
	line2 := ""
	line3 := ""

	for _, id := range tileIDs {
		tile := game.Board.Tiles[id-1] // tiles are 1-indexed in display
		robber := ""
		if game.Board.RobberPosition == id-1 {
			robber = "R"
		}

		// Center the number token with padding
		numStr := fmt.Sprintf("%d", tile.NumberToken)
		if tile.NumberToken < 10 {
			numStr = " " + numStr
		}
		if tile.NumberToken == -1 {
			numStr = " D"
		}

		line1 += fmt.Sprintf("   /%s\\   ", numStr)
		line2 += fmt.Sprintf(" /%s %s\\ ", tile.Resource, robber)
		line3 += fmt.Sprintf("/%s\\ ", strings.Repeat(" ", 5))
	}

	fmt.Println(line1)
	fmt.Println(line2)
	fmt.Println(line3)
}

// cli print the edges as they appear in the valid edge placements
func PrintValidEdges() {

}

func printVertexRow(game *CatanGame, vertexIDs []int, portInfo map[int]string, rowType int) {
	line1 := ""
	line2 := ""

	for _, id := range vertexIDs {
		vertex := game.Board.Graph.Vertices[id]
		player := "0"
		building := " "
		port := " "

		if vertex.OccupiedBy != nil {
			player = fmt.Sprintf("%d", vertex.OccupiedBy.ID)
			if vertex.Building == 1 {
				building = "S" // Settlement
			} else if vertex.Building == 2 {
				building = "C" // City
			}
		}

		if p, exists := portInfo[id]; exists {
			if p == "A" {
				port = "3:1"
			} else {
				port = p + ":1"
			}
		}

		// Alternate between vertex and edge positions based on row type
		if rowType%2 == 0 {
			// Vertex row (main vertices)
			line1 += fmt.Sprintf(" %s%s%s ", port, building, player)
			line2 += "   "
		} else {
			// Edge row (between vertices)
			// Find edge between this vertex and next
			edgePlayer := "0"
			if len(vertex.AdjacentVertexes) > 0 {
				for _, adjID := range vertex.AdjacentVertexes {
					edgeKey := fmt.Sprintf("%d-%d", min(id, adjID), max(id, adjID))
					if edge, exists := game.Board.Graph.Edges[edgeKey]; exists && edge.OccupiedBy != nil {
						edgePlayer = fmt.Sprintf("%d", edge.OccupiedBy.ID)
						break
					}
				}
			}
			line1 += "   "
			line2 += fmt.Sprintf("  %s  ", edgePlayer)
		}
	}

	fmt.Println(line1)
	if rowType%2 == 1 {
		fmt.Println(line2)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
