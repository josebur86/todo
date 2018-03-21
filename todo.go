package main

import(
    "fmt"
    "os"
    "strconv"
)

var GlobalTodoFile = "W:/todo.md"

func FindTaskByLine(tasks []Task, line int) *Task {
    for i, _ := range tasks {
        if tasks[i].FileLine == line {
            return &tasks[i]
        }
    }

    return nil
}

func main() {
    args := os.Args
    argCount := len(args)

    tasks := ReadTodoFile(GlobalTodoFile)

    if argCount == 1 {
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
    } else if argCount > 1 {
        command := args[1]

        // TODO(joe): this is a mess
        if command == "complete" {
            if argCount >= 3 {
                taskLineNum, err := strconv.Atoi(args[2])
                if err != nil {
                    fmt.Println("Invalid line number")
                } else  {
                    task := FindTaskByLine(tasks, taskLineNum)
                    if task != nil {
                        task.Complete = true
                        if err := WriteTodos(GlobalTodoFile, tasks); err != nil {
                            fmt.Print("Unable to write Todos", err)
                        } else {
                            fmt.Print("Completed \"", task.Description, "\"")
                        }
                    } else {
                        fmt.Print("No task on line ", taskLineNum)
                    }
                }
            } else {
                fmt.Println("todo complete: not enough arguments")
            }
        }
    }
}
