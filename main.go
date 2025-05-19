package main

func main() {

	playerSlice := []int{1, 2, 3}
	if cliGameInitialize() == 4 {
		playerSlice = append(playerSlice, 4)
	}
	gameInitialize(len(playerSlice))

	//var game CatanGame = *NewCatanGame(playerSlice)
	//printBoard((*game.Board))

}
