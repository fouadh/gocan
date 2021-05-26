package create_scene

import (
  terminal2 "com.fha.gocan/internal/platform/terminal"
  "fmt"
  "github.com/jmoiron/sqlx"
  "github.com/pborman/uuid"
  "github.com/spf13/cobra"
)

func BuildCreateSceneCmd(db *sqlx.DB, ui terminal2.UI) *cobra.Command {
  return &cobra.Command{
    Use: "create-scene",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      ui.Say("Creating scene...")
      id := uuid.NewUUID().String()
      name := args[0]

      _, err := db.Exec("insert into scenes(id, name) values($1, $2)", id, name)

      if err != nil {
        ui.Failed(fmt.Sprintf("Scene could not be created: %v", err))
      } else {
        ui.Ok()
      }

      return err
    },
  }
}
