// interfaces.go
package gameplay

type GameInitializer interface {
	Initialize() int
}

type GameStarter interface {
	Start(game *CatanGame)
}

type PlayerSelector interface {
	SelectStartingPlayer(game *CatanGame) *Player
}
