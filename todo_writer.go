package main

import(
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
)

type taskDate struct{
    Year int
    Month time.Month
    Day int
}

func NewTaskDate(date time.Time) *taskDate {
    year, month, day := date.Date()
    return &taskDate{year, month, day}
}

func (d *taskDate) Equal(date time.Time) bool {
    year, month, day := date.Date()
    return d.Year == year && d.Month == month && d.Day == day
}

func WriteTodos(filePath string, tasks []Task) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)

    datesWritten := []taskDate{}
    for _, task := range tasks {
        dateWritten := false
        for _, date := range datesWritten {
            if date.Equal(task.Date) {
                dateWritten = true
                break;
            }
        }

        if !dateWritten {
            if len(datesWritten) > 0 {
                writer.WriteByte('\n')
            }
            dateLine := task.Date.Format("Monday - Jan 2, 2006")
            writer.WriteString(dateLine)
            writer.WriteByte('\n')
            writer.WriteString(strings.Repeat("=", len(dateLine)))
            writer.WriteByte('\n')

            datesWritten = append(datesWritten, *NewTaskDate(task.Date))
        }

        completedSymbol := " "
        if task.Complete {
            completedSymbol = "x"
        }
        writer.WriteString(fmt.Sprintf("[%s] %s\n", completedSymbol, task.Description))
    }

    writer.Flush()

    return nil
}
