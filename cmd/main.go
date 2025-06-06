// main.go
package main

import (
	"catango/gameplay"
	"fmt"

	//"os"
	"strings"
)

func main() {
	//cg := &gameplay.CLIGame{Input: os.Stdin}
	input := strings.NewReader("3\n\n\n\n\n\n")
	cg := &gameplay.CLIGame{Input: input}

	playerCount := cg.Initialize()

	game := cg.BaseGame.Initialize(playerCount)
	cg.Start(game)

	playerSelector := &gameplay.CLIPlayerSelector{}
	startingPlayer := playerSelector.SelectStartingPlayer(game, cg.Input)
	fmt.Printf("Starting player is: Player %d\n", startingPlayer.ID)
	cg.SnakeBuild(game, startingPlayer, playerCount)
}
