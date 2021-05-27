package create_scene

import (
  "com.fha.gocan/internal/platform/db"
  terminal2 "com.fha.gocan/internal/platform/terminal"
  "fmt"
  "github.com/pborman/uuid"
  "github.com/spf13/cobra"
)


func BuildCreateSceneCmd(datasource db.DataSource, ui terminal2.UI) *cobra.Command {
  return &cobra.Command{
    Use: "create-scene",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      id := uuid.NewUUID().String()
      name := args[0]
      connection := datasource.GetConnection()

      ui.Say("Creating scene...")
      _, err := connection.Exec("insert into scenes(id, name) values($1, $2)", id, name)

      if err != nil {
        ui.Failed(fmt.Sprintf("Scene could not be created: %v", err))
      } else {
        ui.Ok()
      }

      return err
    },
  }
}
