// initialSetup.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func cliGameInitialize() int {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Welcome to Catan!")
		fmt.Print("Please enter the number of players (3 or 4): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		input = strings.TrimSpace(input)
		playerNum, err := strconv.Atoi(input)
		if err != nil || (playerNum != 3 && playerNum != 4) {
			fmt.Println("Invalid input. Please enter either 3 or 4.")
			continue
		}

		return playerNum
	}
}

func gameInitialize(playerNum int) {
	ids := make([]int, playerNum)
	for i := range ids {
		ids[i] = i + 1
	}

	game := NewCatanGame(ids)
	fmt.Println("Game initialized with", playerNum, "players.")

	gameStart(game)
}

func gameStart(game *CatanGame) {
	fmt.Println("Beginning setup phase...")

	// Snake order placement (forward, then reverse)
	PlaceStartingStructures(game, true)  // First round
	PlaceStartingStructures(game, false) // Second round

	game.Phase = "main"
	fmt.Println("Game phase set to:", game.Phase)
}

func PlaceStartingStructures(game *CatanGame, forward bool) {
	playerOrder := game.Players
	if !forward {
		// Reverse the order
		for i, j := 0, len(playerOrder)-1; i < j; i, j = i+1, j-1 {
			playerOrder[i], playerOrder[j] = playerOrder[j], playerOrder[i]
		}
	}

	reader := bufio.NewReader(os.Stdin)

	for _, player := range playerOrder {
		for {
			fmt.Printf("Player %d: Enter settlement vertex key: ", player.ID)
			vertexKey, _ := reader.ReadString('\n')
			vertexKey = strings.TrimSpace(vertexKey)

			fmt.Printf("Player %d: Enter road edge key: ", player.ID)
			edgeKey, _ := reader.ReadString('\n')
			edgeKey = strings.TrimSpace(edgeKey)

			err := game.PlaceStartingStructures(player, vertexKey, edgeKey)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Println("Placement successful.")
			break
		}
	}
}
