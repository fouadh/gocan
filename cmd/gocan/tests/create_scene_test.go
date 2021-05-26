package tests

import (
  "bytes"
  "com.fha.gocan/internal/create-scene"
  init_db "com.fha.gocan/internal/init-db"
  embeddedpostgres "github.com/fergusstrange/embedded-postgres"
  "github.com/jmoiron/sqlx"
  "github.com/pborman/uuid"
  "github.com/spf13/cobra"
  "testing"
)

const succeed = "\u2713"
const failed = "\u2717"

var postgres *embeddedpostgres.EmbeddedPostgres

func connect(dsn string) (*sqlx.DB, error) {
  db, err := sqlx.Connect("postgres", dsn)
  return db, err
}

func TestCreateScene(t *testing.T) {
  database := embeddedpostgres.NewDatabase()
  if err := database.Start(); err != nil {
    t.Fatalf("%s Cannot start the database: %+v", failed, err)
  }

  defer func() {
    if err := database.Stop(); err != nil {
      t.Fatalf("%s Cannot stop the database: %+v", failed, err)
    }
  }()

  db, err := connect("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
  if err != nil {
    t.Fatalf("%s Cannot connect to the database: %+v", failed, err)
  }

  /*if err := goose.Up(db.DB, "../../../internal/init-db/migrations"); err != nil {
    t.Fatalf("%s Cannot run the migration scripts: %+v", failed, err)
  }*/
  ui := FakeUI{}
  init_db.InitDb("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", &ui)

  name := uuid.New()
  t.Logf("\tGiven no scene named %s exists", name)
  {
    t.Logf("\tWhen I create a scene named %s", name)

    {
      cmd := create_scene.BuildCreateSceneCmd(db, &ui)

      if _, err := runCommand(cmd, name); err != nil {
        t.Fatalf("\t%s Failed to execute create scene command: %+v", failed, err)
      }
      // check that the scene has been added to the db
      var id string
      if err := db.Get(&id, "select id from scenes where name=$1", name); err != nil {
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
// try to create a scene with an empty name
// try to create a scene with a name that is too long for the db

