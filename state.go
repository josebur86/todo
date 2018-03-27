package main

import (
    "encoding/json"
    "log"
    "os"
    "path"
)

type State struct {
    CurrentBoardName string     `json: "currentBoardName"`
    CurrentBoardID string       `json: "currentBoardID"`
}

func ReadTodoState(t *Todo) {
    configPath := path.Join(t.DirectoryPath, ".config.todo")
    file, err := os.Open(configPath)
    if err != nil {
        return
    }
    defer file.Close()

    var state State
    if err = json.NewDecoder(file).Decode(&state); err != nil {
        log.Fatal("Unable to parse ", configPath)
    }

    t.State = state
}

func WriteTodoState(t *Todo) error {
    configPath := path.Join(t.DirectoryPath, ".config.todo")
    file, err := os.Create(configPath)
    if err != nil {
        return err
    }
    defer file.Close()

    if err = json.NewEncoder(file).Encode(t.State); err != nil {
        return err
    }

    return nil
}

