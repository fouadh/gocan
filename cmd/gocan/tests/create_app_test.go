package tests

import (
	"com.fha.gocan/cmd/gocan/tests/support"
	create_app "com.fha.gocan/internal/create-app"
	create_scene "com.fha.gocan/internal/create-scene"
	context "com.fha.gocan/internal/platform"
	"github.com/pborman/uuid"
	"os"
	"testing"
)

func TestCreateApp(t *testing.T) {
	ctx := support.CreateContext()
	defer os.RemoveAll(ctx.Config.EmbeddedDataPath)

	database := support.CreateDatabase(ctx)
	defer database.Stop(ctx.Ui)

	t.Log("\tGiven a scene has been created")
	{
		scene := createScene(ctx)
		name := uuid.New()
		t.Logf("\tWhen I create an app named %s in this scene", name)
		{
			cmd := create_app.NewCommand(ctx)
			if _, err := support.RunCommand(cmd, name, "--scene", scene); err != nil {
				t.Fatalf("\t%s Failed to execute create app command: %+v", failed, err)
			}

			var id string
			connection, _ := ctx.DataSource.GetConnection()
			if err := connection.Get(&id, "select id from apps where name=$1", name); err != nil {
				t.Errorf("\t%s Failed retrieving created app: %+v", failed, err)
			} else {
				t.Logf("\t%s Then the app must have been added to the database", succeed)
			}
		}
	}
}

func createScene(ctx *context.Context) string {
	name := uuid.New()
	request := create_scene.CreateSceneRequest{Name: name}
	create_scene.CreateScene(request, ctx)
	return name
}

