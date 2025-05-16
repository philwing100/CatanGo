package main

import (
	"catango/helpers"
	"fmt"
	"strconv"
)

// Lumber = L
// Brick = B
// Sheep = S
// Ore = O
// Wheat = W
func printBoard(board Board) {
	var output string
	for i := 0; i < len(board.Tiles); i++ {
		var topTemp, middleTemp, bottomTemp string
		for j := 0; j < len(board.Tiles[i]); j++ {
			if board.Tiles[i][j] == nil {
				topTemp += " " + strconv.Itoa(i) + " " + strconv.Itoa(j) + "" // empty space for missing hex
				middleTemp += "    "
				bottomTemp += "    "
			} else {
				topTemp += "/ " + "  \\"
				middleTemp += "|" + board.Tiles[i][j].Resource + helpers.PadString(strconv.Itoa(board.Tiles[i][j].NumberToken), 2) + "|"
				bottomTemp += "\\   /"
			}
		}
		output += topTemp + "\n" + middleTemp + "\n" + bottomTemp + "\n"
	}

	fmt.Print(output)
}

func readInInput() {

}
