package main

import (
  "fmt"
  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use: "gocan",
}

var createCmd = &cobra.Command{
  Use: "create-scene",
  RunE: func(cmd *cobra.Command, args []string) error {
    fmt.Println("create scene")
    return nil
  },
}

func main() {
  rootCmd.AddCommand(createCmd)
  rootCmd.Execute()
}
