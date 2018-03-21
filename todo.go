package main

import(
    "fmt"
    "os"
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

func main() {
    input := NewInput(os.Args)
    commands := InitCommands()
    tasks := ReadTodoFile(GlobalTodoFile)

    commandExecuted := false
    for _, commandDef := range commands {
        if input.Command == commandDef.Name && len(input.Args) >= commandDef.MinArgCount {
            tasks, writeRequired := commandDef.Command(tasks, input.Args)
            if writeRequired {
                if err := WriteTodos(GlobalTodoFile, tasks); err != nil {
                    fmt.Print(err)
                }
            }

            commandExecuted = true
            break;
        }
    }

    if !commandExecuted {
        fmt.Println("Usage: ")
        for _, commandDef := range commands {
            fmt.Println("  ", commandDef.Name, " - ", commandDef.Description)
        }
    }
}
