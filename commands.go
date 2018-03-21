package main

import(
    "fmt"
    "strconv"
    "strings"
    "time"
)

type command func([]Task, []string) ([]Task, bool)
type CommandDefinition struct {
    Name string
    Description string
    MinArgCount int
    Command command
}

func PassesFilter(task Task, filters []string) bool {
    for _, filter := range filters {
        if !strings.Contains(task.Description, filter) {
            return false
        }
    }

    return true
}

func ListTasks(tasks []Task, args []string) ([]Task, bool) {
    taskCount := 0
    for _, task := range tasks {
        if !task.Complete && PassesFilter(task, args) {
            fmt.Printf("%d %s\n", task.FileLine, task.Description)
            taskCount++
        }
    }
    fmt.Println("----")
    fmt.Printf("TODO: %d tasks in %s\n", taskCount, GlobalTodoFile)

    return tasks, false
}

func CompleteTask(tasks []Task, args []string) ([]Task, bool) {
    lineNum, err := strconv.Atoi(args[0])
    if err != nil {
        fmt.Print("Invalid line number ", args[0])
        return tasks, false
    }

    found := false
    for i, _ := range tasks {
        if tasks[i].FileLine == lineNum {
            tasks[i].Complete = true
            fmt.Print("\"", tasks[i].Description, "\" marked complete.")

            found = true
            break;
        }
    }

    if !found {
        fmt.Print("No task on line ", lineNum)
    }

    return tasks, found
}

func AddTask(tasks []Task, args []string) ([]Task, bool) {
    added := false
    if len(args) > 0 {
        description := strings.Join(args, " ")
        task := Task{-1, description, time.Now(), false}

        tasks = append(tasks, task)
        added = true
    }

    return tasks, added
}

func ReviewTask(tasks []Task, args []string) ([]Task, bool) {
    fmt.Println("TODO(joe): Not implemented yet!")

    return tasks, false
}

func InitCommands() []CommandDefinition {
    return []CommandDefinition{
        CommandDefinition{"list", "List all active tasks", 0, ListTasks},
        CommandDefinition{"ls", "List all active tasks", 0, ListTasks},
        CommandDefinition{"complete", "Mark the task at the specified line number complete", 1, CompleteTask},
        CommandDefinition{"add", "Adds a task.", 1, AddTask},
        CommandDefinition{"review", "Does a end/beginning of the day review of unfinished tasks.", 0, ReviewTask},
    }
}

