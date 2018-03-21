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
