package main

import(
    "time"
)

type Workday struct {
    Date time.Time
    Tasks []Task
}

type Task struct {
    FileLine int
    Description string
    Complete bool
}


