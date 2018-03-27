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

func ListTrelloBoardsCommand(t *Todo, command *Command, args []string) error {
    response, err := http.Get(fmt.Sprintf("https://api.trello.com/1/members/%s/boards?key=%s&token=%s", TRELLO_API_USER, TRELLO_API_KEY, TRELLO_API_TOKEN))
    if err != nil {
        return err
    }

    defer response.Body.Close()

    boards := []Board{}
    if err := json.NewDecoder(response.Body).Decode(&boards); err != nil {
        return err
    }

    for _, board := range boards {
        if !board.Closed {
            fmt.Println(board.Name, board.ID)
        }
    }

    return nil
}
