package main

import (
    "bufio"
    "log"
    "os"
)

type FileCursor struct {
    file *os.File
    scanner *bufio.Scanner
    currentLine int
}

func NewFileCursor(filePath string) *FileCursor {
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Error opening %s: ", filePath, err)
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
