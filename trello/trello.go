package trello

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

const (
    BaseURL = "https://api.trello.com/1"

    ListNameDone = "Done"
    ListNameInProgress = "In Progress"
    ListNameNext = "Next"
)

type board struct {
    Name string  `json: "name"`
    ID string    `json: "id"`
    Closed bool  `json: "closed"`
}

type list struct {
    Name string  `json: "name"`
    ID string    `json: "id"`
}

type card struct {
    Name string  `json: "name"`
    ID string    `json: "id"`
    IDShort int  `json: "idShort"`
}

func FetchCardsFromBoard(boardID string) ([]card, error) {
    response, err := http.Get(fmt.Sprintf("%s/boards/%s/lists?cards=none&key=%s&token=%s", BaseURL, boardID, TrelloApiKey, TrelloApiToken))
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    lists := []list{}
    if err := json.NewDecoder(response.Body).Decode(&lists); err != nil {
        return nil, err
    }

    // Find the `In Progress` list
    for _, list := range lists {
        if list.Name == ListNameInProgress {
            cards, err := getCardsFromList(list.ID)
            if err != nil {
                return nil, err
            }

            return cards, nil
        }
    }

    return []card{}, nil
}

func getCardsFromList(listID string) ([]card, error) {
    response, err := http.Get(fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", BaseURL, listID, TrelloApiKey, TrelloApiToken))
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    cards := []card{}
    if err := json.NewDecoder(response.Body).Decode(&cards); err != nil {
        printResponse(response)
        return nil, err
    }

    return cards, nil
}

func FetchOpenBoards() ([]board, error) {
    response, err := http.Get(fmt.Sprintf("%s/members/%s/boards?key=%s&token=%s", BaseURL, TrelloApiUser, TrelloApiKey, TrelloApiToken))
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    boards := []board{}
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
