package main

import(
    "fmt"
    "io/ioutil"
    "log"
)

var GlobalTodoFile = "W:/todo.md"

func main() {
    fmt.Printf("TODO File: %s\n", GlobalTodoFile)

    contents, err := ioutil.ReadFile(GlobalTodoFile)
    if err != nil {
        log.Fatal("Error reading %s: ", GlobalTodoFile, err)
    }

    fmt.Print(string(contents))
}
