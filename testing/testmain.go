// testmain.go
package testing

import (
	"catango/gameplay"
	"fmt"
	"strings"
)

func main() {
	// Simulated user input: "3" for players, then 3x ENTER
	input := strings.NewReader("3\n\n\n\n")
	cg := &gameplay.CLIGame{Input: input}

	playerCount := cg.Initialize()
	game := cg.BaseGame.Initialize(playerCount)
	cg.Start(game)

	playerSelector := &gameplay.CLIPlayerSelector{}
	startingPlayer := playerSelector.SelectStartingPlayer(game, cg.Input)
	fmt.Printf("Starting player is: Player %d\n", startingPlayer.ID)
	cg.SnakeBuild(game, startingPlayer, playerCount)
}
