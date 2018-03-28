package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

const (
    BASE_URL = "https://api.trello.com/1"
)

type Board struct {
    Name string  `json: "name"`
    ID string    `json: "id"`
    Closed bool  `json: "closed"`
}

type List struct {
    Name string  `json: "name"`
    ID string    `json: "id"`
}

type Card struct {
    Name string  `json: "name"`
    ID string    `json: "id"`
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

func TrelloListInProgress(t *Todo, command *Command, args []string) error {
    tasks, err := getTodoTasksFrom(t.State.CurrentBoardID)
    if err != nil {
        return err
    }

    for _, task := range tasks {
        fmt.Printf("%s\n", task.Description)
    }
    fmt.Println("----")
    fmt.Printf("TODO: %d tasks on board %s.\n", len(tasks), t.State.CurrentBoardName)

    return err
}

func getOpenBoards() ([]Board, error) {
    response, err := http.Get(fmt.Sprintf("%s/members/%s/boards?key=%s&token=%s", BASE_URL, TRELLO_API_USER, TRELLO_API_KEY, TRELLO_API_TOKEN))
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

func getTodoTasksFrom(boardID string) ([]Task, error) {
    response, err := http.Get(fmt.Sprintf("%s/boards/%s/lists?cards=none&key=%s&token=%s", BASE_URL, boardID, TRELLO_API_KEY, TRELLO_API_TOKEN))
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    lists := []List{}
    if err := json.NewDecoder(response.Body).Decode(&lists); err != nil {
        return nil, err
    }

    todoTasks := []Task{}

    // TODO(joe): Do both of these at the same time?
    // TODO(joe): How do we order these?

    // Find the `In Progress` list
    for _, list := range lists {
        if list.Name == "In Progress" {
            tasks, err := getTasksFrom(list.ID)
            if err != nil {
                return nil, err
            }

            todoTasks = append(todoTasks, tasks...)
        }
    }

    // Find the `Next` list
    for _, list := range lists {
        if list.Name == "Next" {
            tasks, err := getTasksFrom(list.ID)
            if err != nil {
                return nil, err
            }

            todoTasks = append(todoTasks, tasks...)
        }
    }

    return todoTasks, nil
}

func getTasksFrom(listID string) ([]Task, error) {
    response, err := http.Get(fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", BASE_URL, listID, TRELLO_API_KEY, TRELLO_API_TOKEN))
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    cards := []Card{}
    if err := json.NewDecoder(response.Body).Decode(&cards); err != nil {
        printResponse(response)
        return nil, err
    }

    tasks := []Task{}
    for _, card := range cards {
        tasks = append(tasks, Task{0, card.Name, time.Now(), false})
    }

    return tasks, nil
}

func printResponse(r *http.Response) {
    content, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(string(content[:]))
}
