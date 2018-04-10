package cache

import (
	"fmt"

	"github.com/josebur86/todo/trello"
)

type TrelloCache interface {
	GetOpenBoards() []trello.Board
}

func NewTrelloCache() TrelloCache {
	return sqliteTrelloCache{}
}

type sqliteTrelloCache struct {
}

func (s sqliteTrelloCache) GetOpenBoards() []trello.Board {
	boards, err := trello.FetchOpenBoards()
	if err != nil {
		fmt.Println("Unable to fetch open boards.", err)
	}

	return boards
}
