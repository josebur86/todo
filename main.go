package main

import (
    "fmt"
    "os"
)

func main() {
    todo, err := NewTodo()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    todo.AddCommand(NewCommand("ls", "List all active tasks", ListTasks))
    todo.AddCommand(NewCommand("complete", "Mark the task at the specified line number complete", CompleteTask))
    todo.AddCommand(NewCommand("add", "Adds a task.", AddTask))
    todo.AddCommand(NewCommand("review", "Does a end/beginning of the day review of unfinished tasks.", ReviewTask))
    todo.AddCommand(NewCommand("archive", "Moves completed tasks out of the main TODO file and into a backup/archive file", ArchiveTasks))
    todo.AddCommand(NewCommand("file", "Prints the contents of the TODO file.", FileTask))

    trello := NewCommand("trello", "Subcommand that works with Trello.", TrelloCommand)

    // Boards
    trelloBoards := NewCommand("board", "Subcommand that works with Trello Boards.", TrelloBoardCommand)
    trelloBoards.AddCommand(NewCommand("ls", "List all the Trello Boards.", ListTrelloBoardsCommand))
    trello.AddCommand(trelloBoards)


    todo.AddCommand(trello)


    err = todo.Execute()
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
}
