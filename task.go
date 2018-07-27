package main

import (
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"desc"`
	Date        time.Time `json:"date"`
	Complete    bool      `json:"is_complete"`
}
