// Stores data from the TODO server locally so it can be queried if the
// server is not available.
package cache

import (
	"fmt"

	"github.com/josebur86/todo/trello"
)

// A BoardCache stores data about any one Trello Board.
type BoardCache interface {
	// GetBoardLists returns the open Lists for a Board.
	GetBoardLists() []trello.List

	// GetBoardLists returns the Cards for the given Trello List.
	GetCards(list trello.List) []trello.Card
}

func NewBoardCache(boardName string) BoardCache {
	boards, err := trello.FetchOpenBoards()
	if err != nil {
		fmt.Println("Unable to fetch open boards.")
	}

	for _, board := range boards {
		if board.Name == boardName {
			cache := sqliteBoardCache{boardName, board.ID}
			return cache
		}
	}

	return nil
}

// A BoardCache that stores its data in a SQLite database.
type sqliteBoardCache struct {
	boardName string
	boardID   string
}

func (s sqliteBoardCache) GetBoardLists() []trello.List {
	lists, err := trello.FetchBoardLists(s.boardID)
	if err != nil {
		fmt.Println("Unable to fetch lists from", s.boardName)
	}

	return lists
}

func (s sqliteBoardCache) GetCards(list trello.List) []trello.Card {
	cards, err := trello.FetchCardsFromList(list.ID)
	if err != nil {
		fmt.Println("Unable to fetch cards from", list.Name)
	}

	return cards
}
