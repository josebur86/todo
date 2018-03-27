package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Board struct {
    Name string  `json: "name"`
    ID string    `json: "id"`
    Closed bool  `json: "closed"`
}

func TrelloCommand(t *Todo, command *Command, args []string) error {
    if len(args) > 0 {
        for i, _ := range command.Subcommands {
            subcommand := &command.Subcommands[i]
            if (args[0] == subcommand.Name) {
                if err := subcommand.Exec(t, subcommand, args[1:]); err != nil {
                    return err
                }

                break
            }
        }
    }

    return nil
}

func TrelloBoardCommand(t *Todo, command *Command, args []string) error {
    boards, err := getOpenBoards()
    if err != nil {
        return err
    }

    if len(args) == 0 {
        for _, board := range boards {
            if board.Name == t.State.CurrentBoardName {
                fmt.Println("* ", board.Name)
            } else {
                fmt.Println("  ", board.Name)
            }
        }
    } else {
        for _, board := range boards {
            if board.Name == args[0] {
                t.State.CurrentBoardName = board.Name
                t.State.CurrentBoardID = board.ID
                fmt.Println("Switched to board", board.Name)

                if err = WriteTodoState(t); err != nil {
                    return err
                }

                break
            }
        }
    }

    return nil
}

func getOpenBoards() ([]Board, error) {
    response, err := http.Get(fmt.Sprintf("https://api.trello.com/1/members/%s/boards?key=%s&token=%s", TRELLO_API_USER, TRELLO_API_KEY, TRELLO_API_TOKEN))
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
