package cmd

import (
    "fmt"
    "os"
    "path"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use: "todo",
    Short: "Todo list utility.",
    Long: "Todo list utility that uses Trello as its backend.",
    Run: func(cmd *cobra.Command, args[] string) {
    },
}

func init() {
    cobra.OnInitialize(initConfig)
}

func initConfig() {
    todoDir, defined := os.LookupEnv("TODO_DIR")
    if !defined {
        fmt.Println("TODO_DIR is not defined!")
        os.Exit(1)
    }

    viper.SetConfigFile(path.Join(todoDir, ".todo.json"))

    if err := viper.ReadInConfig(); err != nil {
        fmt.Println("Unable to read todo config file.", err)
        os.Exit(1)
    }
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

