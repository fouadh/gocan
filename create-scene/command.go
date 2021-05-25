package create_scene

import (
  "fmt"
  "github.com/spf13/cobra"
)

func BuildCreateSceneCmd() *cobra.Command {
  return &cobra.Command{
    Use: "create-scene",
    RunE: func(cmd *cobra.Command, args []string) error {
      fmt.Println("create scene")
      return nil
    },
  }
}
