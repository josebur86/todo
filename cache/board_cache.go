// Stores data from the TODO server locally so it can be queried if the
// server is not available.
package cache

import (
    "github.com/josebur86/todo/trello"
)

// A BoardCache stores data about any one Trello Board
type BoardCache interface {
    // Syncs the card data for the list.
    SyncList(listName string) error
}

func NewBoardCache(boardName string) *BoardCache {

}

type sqlBoardCache struct {
    boardName string
    boardID string
}

