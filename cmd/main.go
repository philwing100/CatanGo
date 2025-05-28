// main.go
package main

import (
	"catango/gameplay"
	"fmt"
)

func main() {
	cg := &gameplay.CLIGame{}
	playerCount := cg.Initialize()

	game := cg.BaseGame.Initialize(playerCount)
	cg.Start(game)

	playerSelector := &gameplay.CLIPlayerSelector{}
	startingPlayer := playerSelector.SelectStartingPlayer(game)
	fmt.Printf("Starting player is: Player %d\n", startingPlayer.ID)
	cg.SnakeBuild(game, startingPlayer, playerCount)
}
