package main

import(
    "fmt"
)

var GlobalTodoFile = "W:/todo.md"

func main() {
    tasks := ReadTodoFile(GlobalTodoFile)

    // By default, print out all the uncompleted tasks
    taskCount := 0
    for _, task := range tasks {
        if !task.Complete {
            fmt.Printf("%d %s\n", task.FileLine, task.Description)
            taskCount++
        }
    }
    fmt.Println("---")
    fmt.Printf("TODO: %d tasks in %s\n", taskCount, GlobalTodoFile)
}
