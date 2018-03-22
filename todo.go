package main

import(
    "fmt"
    "os"
    "text/tabwriter"
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
        fmt.Printf("Usage:\n\t%s <command> [args]\n", os.Args[0])
        fmt.Println("Commands:")
        writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
        for _, commandDef := range commands {
            fmt.Fprintf(writer, "\t%s\t- %s\n", commandDef.Name, commandDef.Description)
        }
        writer.Flush()
    }
}
