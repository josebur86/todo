package main

import(
    "fmt"
    "strconv"
    "strings"
    "time"
)

type command func([]Task, []string) []Task
type CommandDefinition struct {
    Name string
    MinArgCount int
    Command command
    RequiresWrite bool
}

func PassesFilter(task Task, filters []string) bool {
    for _, filter := range filters {
        if !strings.Contains(task.Description, filter) {
            return false
        }
    }

    return true
}

func ListTasks(tasks []Task, args []string) []Task {
    taskCount := 0
    for _, task := range tasks {
        if !task.Complete && PassesFilter(task, args) {
            fmt.Printf("%d %s\n", task.FileLine, task.Description)
            taskCount++
        }
    }
    fmt.Println("----")
    fmt.Printf("TODO: %d tasks in %s\n", taskCount, GlobalTodoFile)

    return tasks
}

func CompleteTask(tasks []Task, args []string) []Task {
    lineNum, err := strconv.Atoi(args[0])
    if err != nil {
        fmt.Print("Invalid line number ", args[0])
        return tasks
    }

    found := false
    for i, _ := range tasks {
        if tasks[i].FileLine == lineNum {
            tasks[i].Complete = true
            found = true
            break;
        }
    }

    if !found {
        fmt.Print("No task on line ", lineNum)
    }

    return tasks
}

func AddTask(tasks []Task, args []string) []Task {
    if len(args) > 0 {
        description := strings.Join(args, " ")
        task := Task{-1, description, time.Now(), false}

        tasks = append(tasks, task)
    }

    return tasks
}

func InitCommands() []CommandDefinition {
    return []CommandDefinition{
        CommandDefinition{"list", 0, ListTasks, false},
        CommandDefinition{"ls", 0, ListTasks, false},
        CommandDefinition{"complete", 1, CompleteTask, true},
        CommandDefinition{"add", 1, AddTask, true},
    }
}

