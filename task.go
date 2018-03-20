package main

import(
    "time"
)

type Task struct {
    FileLine int
    Description string
    Date time.Time
    Complete bool
}

