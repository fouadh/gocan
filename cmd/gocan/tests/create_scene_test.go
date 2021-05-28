package tests

import (
  "com.fha.gocan/cmd/gocan/tests/support"
  create_scene "com.fha.gocan/internal/create-scene"
  "github.com/pborman/uuid"
  "os"
  "testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestCreateScene(t *testing.T) {
  ctx := support.CreateContext()
  defer os.RemoveAll(ctx.Config.EmbeddedDataPath)

  database := support.CreateDatabase(ctx)
  defer database.Stop(ctx.Ui)

  name := uuid.New()
  t.Logf("\tGiven no scene named %s exists", name)
  {
    t.Logf("\tWhen I create a scene named %s", name)
    {
      cmd := create_scene.NewCommand(ctx)

      if _, err := support.RunCommand(cmd, name); err != nil {
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

