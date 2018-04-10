package cmd

import (
	"fmt"
	"os"

	"github.com/josebur86/todo/cache"
	"github.com/josebur86/todo/trello"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

type Task struct {
	Name    string
	Handle  int
	Locator string
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List tasks and other todo related items.",
	Long:  "Prints todo tasks from the current board, by default.",
	Run:   runLS,
}

func runLS(cmd *cobra.Command, args []string) {
	boardName := viper.GetString("CurrentBoardName")
	if len(args) == 0 || args[0] == "tasks" {
		boardID := viper.GetString("CurrentBoardID")

		cards, err := trello.FetchCardsFromBoard(boardID)
		if err != nil {
			fmt.Println("Unable to fetch cards from", boardName, err)
			os.Exit(1)
		}

		for _, card := range cards {
			fmt.Printf("%d %s\n", card.IDShort, card.Name)
		}
		fmt.Println("----")
		fmt.Printf("TODO: %d cards on board %s.\n", len(cards), boardName)
	} else if args[0] == "boards" {
		boards := cache.NewTrelloCache().GetOpenBoards()

		for _, board := range boards {
			if board.Name == boardName {
				fmt.Printf("* %s\n", board.Name)
			} else {
				fmt.Printf("  %s\n", board.Name)
			}
		}
	} else if args[0] == "lists" {
		fmt.Println(boardName, "has the following lists:")
		lists := cache.NewBoardCache(boardName).GetBoardLists()
		for _, list := range lists {
			fmt.Printf("  %s\n", list.Name)
		}
	}
}
