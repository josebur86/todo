package main

import(
    "bufio"
    "errors"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
)

var GlobalTodoFile = "W:/todo.md"

type FileCursor struct {
    file *os.File
    scanner *bufio.Scanner
    currentLine int
}

func NewFileCursor(filePath string) *FileCursor {
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Error opening %s: ", GlobalTodoFile, err)
    }

    scanner := bufio.NewScanner(file)
    return &FileCursor{file, scanner, 0}
}

func (f *FileCursor) Advance() bool {
    result := f.scanner.Scan()
    if result {
        f.currentLine++
    }

    return result
}

func (f *FileCursor) LineNumber() int {
    return f.currentLine
}

func (f *FileCursor) Line() string {
    return f.scanner.Text()
}

func (f *FileCursor) Close() {
    f.file.Close()
}

func isLineSectionHeader(line string, headers []string) bool {
    for _, header := range headers {
        if strings.HasPrefix(line, header) {
            return true
        }
    }

    return false
}

func isWeekdaySectionLine(line string) bool {
    weekDayHeaders := []string{
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
        "Sunday",
    }

    return isLineSectionHeader(line, weekDayHeaders)
}

func isGeneralSectionLine(line string) bool {
    generalHeaders := []string{
        "Future",
    }

    return isLineSectionHeader(line, generalHeaders)
}

func isTaskLine(line string) bool {
    return strings.HasPrefix(line, "[ ]") || strings.HasPrefix(line, "[x]")
}

func ParseTask(cursor *FileCursor) (*Task, error) {
    line := cursor.Line()

    task := Task{}
    task.Description = strings.TrimLeft(line, "[x] ")
    task.FileLine = cursor.LineNumber()

    statusIndex := strings.Index(line, "]")-1
    switch completeChar := line[statusIndex]; completeChar {
        case 'x':
            task.Complete = true
        case ' ':
            task.Complete = false
        default:
            return nil, errors.New("Unknown task state")
    }

    return &task, nil
}

func ParseWeekdaySection(cursor *FileCursor) (*Workday, error) {
    fields := strings.Split(cursor.Line(), " - ")
    date, err := time.Parse("Jan 02, 2006", fields[1])
    if err != nil {
        return nil, err
    }
    cursor.Advance()

    tasks := []Task{}        
    for cursor.Advance() {
        if isTaskLine(cursor.Line()) {
            task, err := ParseTask(cursor)
            if err != nil {
                log.Print("Unable to parse task.", err)
                continue;
            }

            tasks = append(tasks, *task)
        } else {
            break;
        }
    }

    return &Workday{date, tasks}, nil
}

func main() {
    cursor := NewFileCursor(GlobalTodoFile)
    defer cursor.Close()

    days := []Workday{}
    
    for cursor.Advance() {
        if isWeekdaySectionLine(cursor.Line()) {
            day, err := ParseWeekdaySection(cursor)
            if err != nil {
                log.Print("Unable to parse a weekday section", err)
                continue
            }

            days = append(days, *day)
        }
    }

    // By default, print out all the uncompleted tasks
    taskCount := 0
    for _, day := range days {
        for _, task := range day.Tasks {
            if !task.Complete {
                fmt.Printf("%d %s\n", task.FileLine, task.Description)
                taskCount++
            }
        }
    }
    fmt.Println("---")
    fmt.Printf("TODO: %d tasks in %s\n", taskCount, GlobalTodoFile)
}
