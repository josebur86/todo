package main

import(
    "fmt"
    "sort"
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

const (
    reviewMove = 0
    reviewKeep = 1
    reviewComplete = 2
    // TODO(joe): Implement these
    //reviewDelete = 3
    //reviewHelp = 4
    reviewQuit = 5
    reviewAskAgain = 6
)

func reviewTask(t Task) int {
    fmt.Printf("\n\t%s\n\n", t.ToString(true))

    fmt.Printf("Done? [y/n/m/q]")
    var resp string
    fmt.Scanf("%s", &resp)
    resp = strings.ToLower(resp)

    result := reviewAskAgain
    if (len(resp) > 0) {
        switch resp[0] {
        case 'y':
            result = reviewComplete
        case 'n':
            result = reviewKeep
        case 'm':
            result = reviewMove
        case 'q':
            result = reviewQuit
        }
    }

    return result
}

type ByDate []Task

func (d ByDate) Len() int { return len(d) }
func (d ByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d ByDate) Less(i, j int) bool { return d[i].Date.Before(d[j].Date) }

func ReviewTask(tasks []Task, args []string) ([]Task, bool) {
    shouldWrite := false
    headerPrinted := false

    for i := 0; i < len(tasks); {
        task := &tasks[i]
        if !task.Complete {
            if !headerPrinted {
                fmt.Println("The following tasks have not been complete:")
                headerPrinted = true
            }

            reviewStatus := reviewTask(*task)
            if reviewStatus == reviewMove {
                task.Date = time.Now()
                shouldWrite = true
                i++
            } else if reviewStatus == reviewComplete {
                task.Complete = true
                shouldWrite = true
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
    sort.Stable(ByDate(tasks))

    return tasks, shouldWrite
}

func FileTask(tasks []Task, args []string) ([]Task, bool) {
    fmt.Println("TODO(joe): Implement!")
    return tasks, false
}

func InitCommands() []CommandDefinition {
    return []CommandDefinition{
        CommandDefinition{"ls", "List all active tasks", 0, ListTasks},
        CommandDefinition{"complete", "Mark the task at the specified line number complete", 1, CompleteTask},
        CommandDefinition{"add", "Adds a task.", 1, AddTask},
        CommandDefinition{"review", "Does a end/beginning of the day review of unfinished tasks.", 0, ReviewTask},
        CommandDefinition{"file", "Prints the contents of the TODO file.", 0, FileTask},
    }
}

