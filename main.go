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

    todo.AddCommand(Command{"ls", "List all active tasks", 0, ListTasks})
    todo.AddCommand(Command{"complete", "Mark the task at the specified line number complete", 1, CompleteTask})
    todo.AddCommand(Command{"add", "Adds a task.", 1, AddTask})
    todo.AddCommand(Command{"review", "Does a end/beginning of the day review of unfinished tasks.", 0, ReviewTask})
    todo.AddCommand(Command{"archive", "Moves completed tasks out of the main TODO file and into a backup/archive file", 0, ReviewTask})
    todo.AddCommand(Command{"file", "Prints the contents of the TODO file.", 0, FileTask})

    err = todo.Execute()
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
}
