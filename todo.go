package main

import(
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

var GlobalTodoFile = "W:/todo.md"

type Input struct {
    Command string
    Args []string
}
func NewInput(args []string) *Input {
    command := ""
    commandArgs := []string{}
    if len(args) > 1 {
        command = args[1]
    }
    if len(args) > 2 {
        commandArgs = args[2:]
    }

    return &Input{command, commandArgs}
}

type command func([]Task, []string) []Task
type CommandDefinition struct {
    Name string
    MinArgCount int
    Command command
    RequiresWrite bool
}

func FindTaskByLine(tasks []Task, line int) *Task {
    for i, _ := range tasks {
        if tasks[i].FileLine == line {
            return &tasks[i]
        }
    }

    return nil
}

func ListIncompleteTasks(tasks []Task, args []string) []Task {
    taskCount := 0
    for _, task := range tasks {
        if !task.Complete {
            fmt.Printf("%d %s\n", task.FileLine, task.Description)
            taskCount++
        }
    }
    fmt.Println("---")
    fmt.Printf("TODO: %d tasks in %s\n", taskCount, GlobalTodoFile)

    return tasks
}

func CompleteTask(tasks []Task, args []string) []Task {
    lineNum, err := strconv.Atoi(args[0])
    if err != nil {
        fmt.Print("Invalid line number ", args[0])
        return tasks
    }

    task := FindTaskByLine(tasks, lineNum)
    if task == nil {
        fmt.Print("No task on line ", lineNum)
        return tasks
    }

    task.Complete = true

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

func main() {

    commands := []CommandDefinition{
        CommandDefinition{"", 0, ListIncompleteTasks, false},
        CommandDefinition{"complete", 1, CompleteTask, true},
        CommandDefinition{"add", 1, AddTask, true},
    }

    input := NewInput(os.Args)
    tasks := ReadTodoFile(GlobalTodoFile)

    for _, commandDef := range commands {
        if input.Command == commandDef.Name && len(input.Args) >= commandDef.MinArgCount {
            tasks = commandDef.Command(tasks, input.Args)
            if commandDef.RequiresWrite {
                if err := WriteTodos(GlobalTodoFile, tasks); err != nil {
                    fmt.Print(err)
                }
            }
        }
    }
}
