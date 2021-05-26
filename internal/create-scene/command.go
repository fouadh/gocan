package create_scene

import (
  "fmt"
  "github.com/jmoiron/sqlx"
  "github.com/pborman/uuid"
  "github.com/spf13/cobra"
)

func BuildCreateSceneCmd(db *sqlx.DB) *cobra.Command {
  return &cobra.Command{
    Use: "create-scene",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      fmt.Println("create scene")
      id := uuid.NewUUID().String()
      name := args[0]
      db.Exec("insert into scenes(id, name) values($1, $2)", id, name)
      return nil
    },
  }
}
