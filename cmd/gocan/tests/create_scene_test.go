package tests

import (
  "bytes"
  "com.fha.gocan/internal/create-scene"
  "github.com/pborman/uuid"
  "testing"
)

func TestCreateScene(t *testing.T) {
  // start the db
  // run the migration scripts

  name := uuid.New()
  cmd := create_scene.BuildCreateSceneCmd()

  // run the command
  buf := new(bytes.Buffer)
  cmd.SetOut(buf)
  cmd.SetErr(buf)
  cmd.SetArgs([]string{name})

  cmd.ExecuteC()

  // check that the scene has been added to the db
}

// todo
// try to create a scene with an empty name
// try to create a scene with a name that is too long for the db

