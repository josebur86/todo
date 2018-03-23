package main

import(
    "errors"
    "log"
    "strings"
    "time"
)

// TODO(joe): This assumes that we are going to get all the tasks from a file which is fine for the
// initial release. It would be cool to aggregate tasks from other sources.

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

func ParseWeekdaySection(cursor *FileCursor) ([]Task, error) {
    date, err := time.Parse("Monday - Jan 02, 2006", cursor.Line())
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
            task.Date = date

            tasks = append(tasks, *task)
        } else {
            break;
        }
    }

    return tasks, nil
}

func ReadTodoFile(filePath string) []Task {
    cursor := NewFileCursor(filePath)
    defer cursor.Close()

    tasks := []Task{}
    for cursor.Advance() {
        if isWeekdaySectionLine(cursor.Line()) {
            tasksForDay, err := ParseWeekdaySection(cursor)
            if err != nil {
                log.Print("Unable to parse a weekday section", err)
                continue
            }

            // STUDY(joe): The ... is kinda cool. It turns the slice into a variadic arguments
            tasks = append(tasks, tasksForDay...)
        }
    }

    return tasks
}

