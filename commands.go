package main

import(
    "bufio"
    "fmt"
    "strconv"
    "strings"
    "time"
    "os"
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

const (
    reviewMove = 0
    reviewKeep = 1
    reviewAskAgain = 2
    reviewQuit = 3
)

func reviewTask(t Task) int {
    fmt.Printf("\n\t%s\n\n", t.ToString(true))

    fmt.Printf("Move to today? [y/n/q]")
    input := bufio.NewReader(os.Stdin)
    resp, err := input.ReadString('\n')
    if err != nil {
        fmt.Println("Unable to process input", err)
    }
    resp = strings.ToLower(resp)

    result := reviewAskAgain
    if (len(resp) > 0) {
        switch resp[0] {
        case 'y':
            result = reviewMove
        case 'n':
            result = reviewKeep
        case 'q':
            result = reviewQuit
        }
    }

    return result
}

func ReviewTask(tasks []Task, args []string) ([]Task, bool) {
    headerPrinted := false
    for i := 0; i < len(tasks); {
        task := tasks[i]
        if !task.Complete {
            if !headerPrinted {
                fmt.Println("The following tasks have not been complete:")
                headerPrinted = true
            }

            reviewStatus := reviewTask(task)
            if reviewStatus == reviewMove {
                // TODO(joe): Actually move
                i++
            } else if reviewStatus == reviewAskAgain {
                continue
            } else if reviewStatus == reviewQuit {
                break
            } else {
                i++
            }
        }
    }

    return tasks, false
}

func InitCommands() []CommandDefinition {
    return []CommandDefinition{
        CommandDefinition{"ls", "List all active tasks", 0, ListTasks},
        CommandDefinition{"complete", "Mark the task at the specified line number complete", 1, CompleteTask},
        CommandDefinition{"add", "Adds a task.", 1, AddTask},
        CommandDefinition{"review", "Does a end/beginning of the day review of unfinished tasks.", 0, ReviewTask},
    }
}

