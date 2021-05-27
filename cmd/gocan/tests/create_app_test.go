package tests

import (
	create_app "com.fha.gocan/internal/create-app"
	context "com.fha.gocan/internal/platform"
	"com.fha.gocan/internal/platform/db"
	"github.com/pborman/uuid"
	"testing"
)

func TestCreateApp(t *testing.T) {
	database := db.EmbeddedDatabase{}
	ui := FakeUI{}
	database.Start(&ui)
	defer database.Stop(&ui)
	ctx := context.New(dsn, &ui)

	t.Log("\tGiven a scene has been created")
	{
		scene := createScene(ctx)
		name := uuid.New()
		t.Logf("\tWhen I create an app named %s in this scene", name)
		{
			cmd := create_app.NewCommand(ctx)
			if _, err := runCommand(cmd, name, "--scene", scene); err != nil {
				t.Fatalf("\t%s Failed to execute create app command: %+v", failed, err)
			}
			t.Logf("\t%s Then the app must have been added to the database", succeed)
		}
	}
}

func createScene(ctx *context.Context) string {
	return ""
}

