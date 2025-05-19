package main

import (
	"catango/helpers"
	"fmt"
	"strconv"
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

func printBoard(board Board) {
	for i := 0; i < len(board.Tiles); i++ {
		var topRow, midRow, botRow string

		// indent every other row for hex alignment
		indent := ""
		if i%2 != 0 {
			indent = "   "
		}
		topRow += indent
		midRow += indent
		botRow += indent

		for j := 0; j < len(board.Tiles[i]); j++ {
			tile := board.Tiles[i][j]
			if tile == nil {
				topRow += "      "
				midRow += "      "
				botRow += "      "
			} else {
				// Match port if any
				portLabel := ""
				for _, port := range board.Ports {
					if port.q == tile.q && port.r == tile.r {
						if port.GiveResource == "any" {
							portLabel = "3:1"
						} else {
							portLabel = port.GiveResource[:1] + ":2"
						}
					}
				}

				resource := tile.Resource
				token := ""
				if tile.NumberToken > 0 {
					token = helpers.PadString(strconv.Itoa(tile.NumberToken), 2)
				} else {
					token = "  "
				}

				label := fmt.Sprintf("%s%s", resource, token)
				if portLabel == "" {
					portLabel = "‾‾‾"
				} else if portLabel == "?" {
					portLabel = " ? "
				}

				topRow += " /" + portLabel + "\\ "
				midRow += "| " + label + " |"
				botRow += " \\___/ "
			}
		}
		fmt.Println(topRow)
		fmt.Println(midRow)
		fmt.Println(botRow)
	}
	fmt.Println("\nPlayer Structures:")
	for _, v := range board.Graph.Vertices {
		if v.OccupiedBy != nil {
			structure := "S"
			if v.Building == "city" {
				structure = "C"
			}
			fmt.Printf("Vertex %s: %s\n", v.ID, colorText(structure, v.OccupiedBy.ID))
		}
	}
	for _, e := range board.Graph.Edges {
		if e.OccupiedBy != nil {
			fmt.Printf("Edge %s: %s\n", e.ID, colorText("R", e.OccupiedBy.ID))
		}
	}

}

func listAvailableSettlementVertices(graph *Graph) []string {
	var available []string

	for id, vertex := range graph.Vertices {
		if vertex == nil || vertex.OccupiedBy != nil {
			continue
		}

		valid := true

		for _, edge := range vertex.Edges {
			if edge == nil {
				continue
			}

			for _, neighbor := range edge.Vertices {
				if neighbor == nil {
					continue
				}
				if neighbor != vertex && neighbor.OccupiedBy != nil {
					valid = false
					break
				}
			}

			if !valid {
				break
			}
		}

		if valid {
			available = append(available, id)
		}
	}

	return available
}

func listAvailableRoadEdgesForVertex(vertex *Vertex) []string {
	var available []string
	for _, edge := range vertex.Edges {
		if edge.OccupiedBy == nil {
			available = append(available, edge.ID)
		}
	}
	return available
}
