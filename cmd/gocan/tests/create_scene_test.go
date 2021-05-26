package tests

import (
  "bytes"
  "com.fha.gocan/internal/create-scene"
  "github.com/pborman/uuid"
  "testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestCreateScene(t *testing.T) {
  // start the db
  // run the migration scripts

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

