package main

import (
	"catango/gameplay"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Check if running in debug mode via env var
	testMode := os.Getenv("DEBUG_MODE") == "true"

	var cg *gameplay.CLIGame
	if testMode {
		input := strings.NewReader("3\n\n\n\n\n\n") // Replace with appropriate test inputs
		cg = &gameplay.CLIGame{Input: input}
	} else {
		cg = &gameplay.CLIGame{Input: os.Stdin}
	}

	playerCount := cg.Initialize()
	game := cg.BaseGame.Initialize(playerCount)

	cg.Start(game)

	playerSelector := &gameplay.CLIPlayerSelector{}
	startingPlayer := playerSelector.SelectStartingPlayer(game, cg.Input)

	fmt.Printf("Starting player is: Player %d\n", startingPlayer.ID)
	cg.SnakeBuild(game, startingPlayer, playerCount)
}
