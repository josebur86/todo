package trello

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	BaseURL = "https://api.trello.com/1"

	ListNameDone       = "Done"
	ListNameInProgress = "In Progress"
	ListNameNext       = "Next"
)

type Board struct {
	Name   string `json: "name"`
	ID     string `json: "id"`
	Closed bool   `json: "closed"`
}

type List struct {
	Name string `json: "name"`
	ID   string `json: "id"`
}

type Card struct {
	Name    string `json: "name"`
	ID      string `json: "id"`
	IDShort int    `json: "idShort"`
}

func FetchBoardLists(boardID string) ([]List, error) {
	response, err := http.Get(fmt.Sprintf("%s/boards/%s/lists?cards=none&key=%s&token=%s", BaseURL, boardID, TrelloApiKey, TrelloApiToken))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	lists := []List{}
	if err := json.NewDecoder(response.Body).Decode(&lists); err != nil {
		return nil, err
	}

	return lists, nil
}

func FetchCardsFromBoard(boardID string) ([]Card, error) {
	lists, err := FetchBoardLists(boardID)
	if err != nil {
		return nil, err
	}

	// Find the `In Progress` list
	for _, list := range lists {
		if list.Name == ListNameInProgress {
			cards, err := FetchCardsFromList(list.ID)
			if err != nil {
				return nil, err
			}

			return cards, nil
		}
	}

	return []Card{}, nil
}

func FetchCardsFromList(listID string) ([]Card, error) {
	response, err := http.Get(fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", BaseURL, listID, TrelloApiKey, TrelloApiToken))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	cards := []Card{}
	if err := json.NewDecoder(response.Body).Decode(&cards); err != nil {
		printResponse(response)
		return nil, err
	}

	return cards, nil
}

func FetchOpenBoards() ([]Board, error) {
	response, err := http.Get(fmt.Sprintf("%s/members/%s/boards?key=%s&token=%s", BaseURL, TrelloApiUser, TrelloApiKey, TrelloApiToken))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	boards := []Board{}
	if err := json.NewDecoder(response.Body).Decode(&boards); err != nil {
		return nil, err
	}

	// Remove closed boards
	for i := 0; i < len(boards); {
		if boards[i].Closed {
			boards = append(boards[:i], boards[i+1:]...)
		} else {
			i++
		}
	}

	return boards, nil
}

func printResponse(r *http.Response) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(content[:]))
}
