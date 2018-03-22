package main

import(
    "fmt"
    "time"
)

type Task struct {
    FileLine int
    Description string
    Date time.Time
    Complete bool
}

func (t *Task) ToString(withDate bool) string {
    if withDate {
        return fmt.Sprintf("From %s: \"%s\"", t.Date.Format("Monday - Jan 2, 2006"), t.Description)
    } else {
        return fmt.Sprintf("%d %s", t.FileLine, t.Description)
    }
}

