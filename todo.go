package main

import(
    "errors"
    "fmt"
    "os"
    "path"
    "text/tabwriter"
)

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

type Todo struct {
    DirectoryPath string
    FilePath string
    Tasks []Task
    commands []Command
}

func NewTodo() (*Todo, error) {
    todoDir, defined := os.LookupEnv("TODO_DIR")
    if !defined {
        return nil, errors.New("TODO_DIR is not defined")
    }

    return &Todo{todoDir, path.Join(todoDir, "todo.md"), []Task{}, []Command{}}, nil
}

func (t *Todo) AddCommand(command Command) {
    t.commands = append(t.commands, command)
}

func (t *Todo) Execute() error {
    input := NewInput(os.Args)
    t.Tasks = ReadTodoFile(t.FilePath)

    commandExecuted := false
    for _, command := range t.commands {
        if input.Command == command.Name {
            err := command.Exec(t, input.Args)
            if err != nil {
                return err
            }

            commandExecuted = true
            break;
        }
    }

    if !commandExecuted {
        fmt.Printf("Usage:\n\t%s <command> [args]\n", os.Args[0])
        fmt.Println("Commands:")
        writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
        for _, command := range t.commands {
            fmt.Fprintf(writer, "\t%s\t- %s\n", command.Name, command.Description)
        }
        writer.Flush()
    }

    return nil
}

