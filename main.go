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

    todo.AddCommand(Command{"ls", "List all active tasks", ListTasks})
    todo.AddCommand(Command{"complete", "Mark the task at the specified line number complete", CompleteTask})
    todo.AddCommand(Command{"add", "Adds a task.", AddTask})
    todo.AddCommand(Command{"review", "Does a end/beginning of the day review of unfinished tasks.", ReviewTask})
    todo.AddCommand(Command{"archive", "Moves completed tasks out of the main TODO file and into a backup/archive file", ArchiveTasks})
    todo.AddCommand(Command{"file", "Prints the contents of the TODO file.", FileTask})

    err = todo.Execute()
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
}
