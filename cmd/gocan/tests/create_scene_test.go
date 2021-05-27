package tests

import (
  "bytes"
  create_scene "com.fha.gocan/internal/create-scene"
  context "com.fha.gocan/internal/platform"
  "com.fha.gocan/internal/platform/db"
  "github.com/pborman/uuid"
  "github.com/spf13/cobra"
  "testing"
)

const succeed = "\u2713"
const failed = "\u2717"

const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

func TestCreateScene(t *testing.T) {
  database := db.EmbeddedDatabase{}
  ui := FakeUI{}
  database.Start(&ui)
  defer database.Stop(&ui)
  db.Migrate(dsn, &ui)

  name := uuid.New()
  t.Logf("\tGiven no scene named %s exists", name)
  {
    t.Logf("\tWhen I create a scene named %s", name)
    {
      ctx, _ := context.New(&ui)
      cmd := create_scene.NewCommand(ctx)

      if _, err := runCommand(cmd, name); err != nil {
        t.Fatalf("\t%s Failed to execute create scene command: %+v", failed, err)
      }

      var id string
      connection, _ := ctx.DataSource.GetConnection()
      if err := connection.Get(&id, "select id from scenes where name=$1", name); err != nil {
        t.Errorf("\t%s Failed retrieving created scene: %+v", failed, err)
      } else {
        t.Logf("\t%s Then the scene must have been added to the database", succeed)
      }
    }
  }
}

func runCommand(cmd *cobra.Command, args ...string) (string, error) {
  buf := new(bytes.Buffer)
  cmd.SetOut(buf)
  cmd.SetErr(buf)
  cmd.SetArgs(args)

  _, err := cmd.ExecuteC()

  return buf.String(), err
}

// todo
// try to create a scene with a name that is too long for the db

