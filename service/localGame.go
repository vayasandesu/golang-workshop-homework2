package service

import (
	"goworkshop2/game"
)

type LocalGame struct {
	Players []game.Character `json:"players"`
}

func (state *LocalGame) Join(name string) (game.Character, error) {
	character := game.NewCharacter(name)
	state.Players = append(state.Players, character)

	return character, nil
}

func (state *LocalGame) List() ([]game.Character, error) {
	return state.Players, nil
}

func (state *LocalGame) getCharacter(name string) game.Character {
	players := state.Players
	for _, value := range players {
		if value.Name == name {
			return value
		}
	}
	return game.Character{}
}
