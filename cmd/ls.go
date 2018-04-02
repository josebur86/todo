package cmd

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

func init() {
    rootCmd.AddCommand(lsCmd)
}

type Task struct {
    Description string
}

var lsCmd = &cobra.Command{
    Use: "ls",
    Short: "List tasks and other todo related items.",
    Long: "Prints todo tasks from the current board, by default.",
    Run: runLS,
}

const (
    BASE_URL = "https://api.trello.com/1"
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
}

func runLS(cmd *cobra.Command, args []string) {
    boardName := viper.GetString("CurrentBoardName")
    if len(args) == 0 || args[0] == "tasks" {
        boardID := viper.GetString("CurrentBoardID")

        tasks, err := fetchTasksFrom(boardID)
        if err != nil {
            fmt.Println("Unable to fetch tasks from", boardName, err)
            os.Exit(1)
        }

        for _, task := range tasks {
            fmt.Printf("%s\n", task.Description)
        }
        fmt.Println("----")
        fmt.Printf("TODO: %d tasks on board %s.\n", len(tasks), boardName)
    } else if args[0] == "boards" {
        boards, err := fetchOpenBoards()
        if err != nil {
            fmt.Println("Unable to fetch boards", err)
            os.Exit(1)
        }

        for _, board := range boards {
            if board.Name == boardName {
                fmt.Printf("* %s\n", board.Name)
            } else {
                fmt.Printf("  %s\n", board.Name)
            }
        }
    }
}

func fetchTasksFrom(boardID string) ([]Task, error) {
    response, err := http.Get(fmt.Sprintf("%s/boards/%s/lists?cards=none&key=%s&token=%s", BASE_URL, boardID, TRELLO_API_KEY, TRELLO_API_TOKEN))
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    lists := []list{}
    if err := json.NewDecoder(response.Body).Decode(&lists); err != nil {
        return nil, err
    }

    todoTasks := []Task{}

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

    return todoTasks, nil
}

func getTasksFrom(listID string) ([]Task, error) {
    response, err := http.Get(fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", BASE_URL, listID, TRELLO_API_KEY, TRELLO_API_TOKEN))
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    cards := []card{}
    if err := json.NewDecoder(response.Body).Decode(&cards); err != nil {
        printResponse(response)
        return nil, err
    }

    tasks := []Task{}
    for _, card := range cards {
        tasks = append(tasks, Task{card.Name})
    }

    return tasks, nil
}

func fetchOpenBoards() ([]board, error) {
    response, err := http.Get(fmt.Sprintf("%s/members/%s/boards?key=%s&token=%s", BASE_URL, TRELLO_API_USER, TRELLO_API_KEY, TRELLO_API_TOKEN))
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
