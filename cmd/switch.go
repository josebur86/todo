package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

func init() {
    rootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
    Use: "switch",
    Short: "Switch to another board.",
    Run: runSwitch,
    Args: cobra.MinimumNArgs(1),
}

func runSwitch(cmd *cobra.Command, args []string) {
    newBoard := args[0]

    boards, err := fetchOpenBoards()
    if err != nil {
        fmt.Println("Unable to fetch boards", err)
        os.Exit(1)
    }

    for _, board := range boards {
        if board.Name == newBoard {
            viper.Set("CurrentBoardName", board.Name)
            viper.Set("CurrentBoardID", board.ID)
            viper.WriteConfig()

            fmt.Println("Switched to", board.Name)
            break
        }
    }
}
