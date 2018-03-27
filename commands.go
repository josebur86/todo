package main

import(
    "errors"
    "fmt"
    "path"
    "sort"
    "strconv"
    "strings"
    "time"
)

type Exec func(*Todo, *Command, []string) error
type Command struct {
    Name string
    Description string
    Exec Exec
    Subcommands []Command
}

func NewCommand(name, description string, exec Exec) Command {
    return Command{name, description, exec, []Command{}}
}

func (c *Command) AddCommand(command Command) {
    c.Subcommands = append(c.Subcommands, command)
}

func PassesFilter(task Task, filters []string) bool {
    for _, filter := range filters {
        if !strings.Contains(task.Description, filter) {
            return false
        }
    }

    return true
}

func ListTasks(t *Todo, command *Command, args []string) (error) {
    taskCount := 0
    for _, task := range t.Tasks {
        if !task.Complete && PassesFilter(task, args) {
            fmt.Printf("%d %s\n", task.FileLine, task.Description)
            taskCount++
        }
    }
    fmt.Println("----")
    fmt.Printf("TODO: %d tasks in %s.\n", taskCount, t.FilePath)

    return nil
}

func CompleteTask(t *Todo, command *Command, args []string) error {
    if len(args) < 1 {
        return errors.New("No line number")
    }

    lineNum, err := strconv.Atoi(args[0])
    if err != nil {
        return errors.New(fmt.Sprintf("Invalid line number %s", args[0]))
    }

    found := false
    for i, _ := range t.Tasks {
        if t.Tasks[i].FileLine == lineNum {
            t.Tasks[i].Complete = true
            fmt.Print("\"", t.Tasks[i].Description, "\" marked complete.")

            found = true
            break;
        }
    }

    if !found {
        return errors.New(fmt.Sprintf("No task on line %d", lineNum))
    }

    err = WriteTodos(t.FilePath, t.Tasks)
    if err != nil {
        return err
    }

    return nil
}

func AddTask(t *Todo, command *Command, args []string) error {
    if len(args) > 0 {
        description := strings.Join(args, " ")
        task := Task{-1, description, time.Now(), false}

        t.Tasks = append(t.Tasks, task)

        err := WriteTodos(t.FilePath, t.Tasks)
        if err != nil {
            return err
        }
    }

    return nil
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

// TODO(joe): Simplify?
type ByDate []Task
func (d ByDate) Len() int { return len(d) }
func (d ByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d ByDate) Less(i, j int) bool { return d[i].Date.Before(d[j].Date) }

func ReviewTask(t *Todo, command *Command, args []string) error {
    shouldWrite := false
    headerPrinted := false

    for i := 0; i < len(t.Tasks); {
        task := &t.Tasks[i]
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

    if shouldWrite {
        sort.Stable(ByDate(t.Tasks))
        err := WriteTodos(t.FilePath, t.Tasks)
        if err != nil {
            return err
        }
    }

    return nil
}

func FileTask(t *Todo, command *Command, args []string) error {
    fmt.Println("TODO(joe): Implement!")
    return nil
}

func ArchiveTasks(t *Todo, command *Command, args []string) error {
    fmt.Println("The following tasks will be archived:")

    toArchive := []Task{}
    for i := 0; i < len(t.Tasks); {
        task := t.Tasks[i]
        if task.Complete {
            fmt.Printf("\t%d %s\n", task.FileLine, task.Description)

            t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
            toArchive = append(toArchive, task)
        } else {
            i++
        }
    }

    if len(toArchive) > 0 {
        // Write the archive file.
        // TODO(joe): Check to see if the file is already there and append to it.
        archivePath := path.Join(t.DirectoryPath, fmt.Sprintf("%s.md", toArchive[0].Date.Format("Jan-02-2006")))
        fmt.Println(archivePath)
        if err := WriteTodos(archivePath, toArchive); err != nil {
            return err
        }

        // Write back the main todo file.
        if err := WriteTodos(t.FilePath, t.Tasks); err != nil {
            return err
        }
    }

    return nil
}

