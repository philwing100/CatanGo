// cliGame.go
package gameplay

import (
	"bufio"
	"catango/helpers"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CLIGame struct {
	BaseGame // Embed the base implementation
}

func (cg *CLIGame) Initialize() int {
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

func (cg *CLIGame) Start(game *CatanGame) {
	fmt.Println("Game is starting!")
	fmt.Println("Current Phase:", game.Phase)
	PrintGameBoard(game)
	cg.BaseGame.Start(game) // Call base implementation
	fmt.Println("Game phase set to:", game.Phase)
}

type CLIPlayerSelector struct {
	BasePlayerSelector // Embed the base implementation
}

func (cps *CLIPlayerSelector) SelectStartingPlayer(game *CatanGame) *Player {
	reader := bufio.NewReader(os.Stdin)

	rollFunc := func(player *Player) int {
		fmt.Printf("Player %d, press ENTER to roll the die...", player.ID)
		reader.ReadString('\n')
		roll := helpers.RollDie()
		fmt.Printf("Player %d rolled a %d\n", player.ID, roll)
		return roll
	}

	fmt.Println("\n=== Starting Player Selection ===")
	winner := cps.BasePlayerSelector.SelectStartingPlayer(game, rollFunc)
	fmt.Printf("🎉 Player %d will go first!\n", winner.ID)
	return winner
}

func (cg *CLIGame) SnakeBuild(game *CatanGame, startingPlayer *Player, playerCount int) {
	fmt.Println("\n=== Starting Build Phase ===")
	order := GenerateSnakeOrder(game, startingPlayer, playerCount)

	for _, playerID := range order {
		player := game.Players[playerID]
		PrintRaw(game)
		var ValidVertices = ComputeValidVertexPlacements(game)
		var vertexID1 int
		var vertexID2 int

		i := 0
		for i < 1 {
			vertexID1 = readInt(fmt.Sprintf("Enter the ID of the vertex where Player %d wants to build a settlement: ", player.ID))
			if helpers.ContainsInt(ValidVertices, vertexID1) {
				PlaceSettlement(vertexID1, player, game)
				if player.VictoryPoints == 1 {
					//if the player has one victory point at this point then they get their resources
					for _, tileID := range GetVertexByID(game, vertexID1).TileIds {
						BankToPlayerResource(game, player, GetTileByID(game, tileID).Resource, 1)
					}
				}
				i++
			}
		}

		fmt.Printf("Player %d's turn to select a road from that settlement", player.ID)
		for i > 0 {
			vertexID2 = readInt(fmt.Sprintf("Enter the ID of the vertex where Player %d wants to build a road: ", player.ID))
			if helpers.ContainsInt(GetAdjacentVertices(vertexID1, game), vertexID2) {
				PlaceRoad(vertexID1, vertexID2, player, game)
				i--
			}
		}
	}

	fmt.Println("Snake building phase completed!")
	PrintGameBoard(game)
}
