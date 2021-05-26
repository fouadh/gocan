package tests

import (
  "bytes"
  "com.fha.gocan/internal/create-scene"
  "fmt"
  embeddedpostgres "github.com/fergusstrange/embedded-postgres"
  "github.com/jmoiron/sqlx"
  "github.com/pborman/uuid"
  "testing"
)

const succeed = "\u2713"
const failed = "\u2717"

var postgres *embeddedpostgres.EmbeddedPostgres

func setupDb() {
  fmt.Println("----- init")
  postgres = embeddedpostgres.NewDatabase()
  postgres.Start()
  // start the db
  // run the migration scripts
}

func connect() (*sqlx.DB, error) {
  db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
  return db, err
}

func TestCreateScene(t *testing.T) {
  database := embeddedpostgres.NewDatabase()
  if err := database.Start(); err != nil {
    t.Fatal(err)
  }

  defer func() {
    if err := database.Stop(); err != nil {
      t.Fatal(err)
    }
  }()

  _, err := connect()
  if err != nil {
    t.Fatal(err)
  }

  name := uuid.New()
  t.Logf("\tGiven no scene named %s exists", name)
  {
    t.Logf("\tWhen I create a scene named %s",  name)
    {
      cmd := create_scene.BuildCreateSceneCmd()
      // run the command
      buf := new(bytes.Buffer)
      cmd.SetOut(buf)
      cmd.SetErr(buf)
      cmd.SetArgs([]string{name})

      if _, err := cmd.ExecuteC(); err != nil {
        t.Fatalf("\t%s Failed to execute create scene command: %+v", failed, err)
      }
      // check that the scene has been added to the db
      t.Logf("\t%s Then the scene must have been added to the database", succeed)
    }
  }
}

// todo
// try to create a scene with an empty name
// try to create a scene with a name that is too long for the db

