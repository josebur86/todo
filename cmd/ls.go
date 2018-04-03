package cmd

import (
    "fmt"
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

