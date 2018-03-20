package main

import(
    "fmt"
    "strings"
    "time"
)

func WriteTodos(filePath string, tasks []Task) error {
    datesWritten := []time.Time{}    

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
                fmt.Println()
            }
            dateLine := task.Date.Format("Monday - Jan 2, 2006")
            fmt.Println(dateLine)
            fmt.Println(strings.Repeat("=", len(dateLine)))

            datesWritten = append(datesWritten, task.Date)
        }

        completedSymbol := " "
        if task.Complete {
            completedSymbol = "x"
        }
        fmt.Printf("[%s] %s\n", completedSymbol, task.Description)
    }

    return nil
}
