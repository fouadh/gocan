package e2e

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestE2E(t *testing.T) {
	dir := setupDatabase(t)
	defer os.RemoveAll(dir)
	defer runCommand(t, "stop-db")
	startDatabase(t)
	sceneName := "a-scene"
	createScene(t, sceneName)
	assertSceneCanBeRetrieved(t, sceneName)
	appName := "an-app"
	createApp(t, appName, sceneName)
	assertAppCanBeRetrieved(t, appName, sceneName)
	importHistory(t, appName, sceneName)
	assertAppSummaryCanBeRetrieved(t, appName, sceneName)
	assertDevsCanBeRetrieved(t, appName, sceneName)
	assertRevisionsCanBeRetrieved(t, appName, sceneName)
}

func importHistory(t *testing.T, appName string, sceneName string) {
	appFolder := createTempFolder(t)
	defer os.RemoveAll(appFolder)

	data := `
	public class Hello {
      public static void main(String[] args) {
		System.out.println("hello from gocan !");
      }
	}
`
	if err := ioutil.WriteFile(appFolder+"/Hello.java", []byte(data), 0755); err != nil {
		t.Fatalf("Unable to create file")
	}

	out, err := exec.Command("/bin/sh", "-c", "cd "+appFolder+"; git init; git config user.name \"Developer 1\"; git add .; git commit -m 'init repo'").Output()
	if err != nil {
		t.Log(string(out))
		t.Fatalf("%s Initializing git repo", failed)
		return
	}

	output := runCommand(t, "import-history", appName, "--scene", sceneName, "--directory", appFolder)
	if strings.Contains(output, "Importing history") && strings.Contains(output, "OK") {
		t.Logf("%s History imported", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s History Import failed", failed)
	}
}

func assertDevsCanBeRetrieved(t *testing.T, appName string, sceneName string) {
	output := runCommand(t, "devs", appName, "--scene", sceneName)
	if strings.Contains(output, "Developer 1") &&
		strings.Contains(output, "OK") {
		t.Logf("%s Developers retrieved", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Retrieving developers", failed)
	}
}

func assertAppSummaryCanBeRetrieved(t *testing.T, appName string, sceneName string) {
	output := runCommand(t, "app-summary", appName, "--scene", sceneName)
	if strings.Contains(output, appName) &&
		strings.Contains(output, "1") &&
		strings.Contains(output, "OK") {
		t.Logf("%s App summary retrieved", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Retrieving app summary", failed)
	}
}

func assertAppCanBeRetrieved(t *testing.T, appName string, sceneName string) {
	output := runCommand(t, "apps", "--scene", sceneName)
	if strings.Contains(output, appName) &&
		strings.Contains(output, "OK") {
		t.Logf("%s App %s retrieved", succeed, appName)
	} else {
		t.Log(output)
		t.Fatalf("%s Retrieving app failed", failed)
	}
}

func assertSceneCanBeRetrieved(t *testing.T, sceneName string) {
	output := runCommand(t, "scenes")
	if strings.Contains(output, sceneName) &&
		strings.Contains(output, "OK") {
		t.Logf("%s Scene %s retrieved", succeed, sceneName)
	} else {
		t.Log(output)
		t.Fatalf("%s Retrieving scene failed", failed)
	}
}

func assertRevisionsCanBeRetrieved(t *testing.T, appName string, sceneName string) {
	output := runCommand(t, "revisions", appName, "--scene", sceneName)
	if strings.Contains(output, "Hello.java") &&
		strings.Contains(output, "OK") {
		t.Logf("%s Revisions retrieved", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Retrieving revisions failed", failed)
	}
}

func createTempFolder(t *testing.T) string {
	dir, err := ioutil.TempDir("", "gocan-")
	if err != nil {
		t.Log(err)
		t.Fatalf("%s Creating temp folder", failed)
	}
	return dir
}

func createApp(t *testing.T, appName string, sceneName string) {
	output := runCommand(t, "create-app", appName, "--scene", sceneName)
	if strings.Contains(output, "Creating the app") && strings.Contains(output, "OK") {
		t.Logf("%s App created", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s App creation failed", failed)
	}
}

func createScene(t *testing.T, sceneName string) {
	output := runCommand(t, "create-scene", sceneName)
	match, _ := regexp.MatchString("Scene .* created", output)
	if match {
		t.Logf("%s Scene created", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Scene creation failed", failed)
	}
}

func startDatabase(t *testing.T) {
	output := startCommand(t, 10*time.Second, "start-db")
	if strings.Contains(output, "database system is ready to accept connections") {
		t.Logf("%s Database started", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Database could not be started", failed)
	}
}

func setupDatabase(t *testing.T) string {
	dir := createTempFolder(t)
	output := runCommand(t, "setup-db", "--directory", dir, "--port", "5433")
	if strings.Contains(output, "Database has been configured") {
		t.Logf("%s Database configured", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Database configuration failed", failed)
	}
	return dir
}

func startCommand(t *testing.T, sleep time.Duration, args ...string) string {
	command := exec.Command("../bin/gocan", args...)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out

	err := command.Start()

	if err != nil {
		t.Logf("Failed to execute the command: %s", err)
	}

	time.Sleep(sleep)

	return string(out.Bytes())
}

func runCommand(t *testing.T, args ...string) string {
	command := exec.Command("../bin/gocan", args...)
	out, err := command.CombinedOutput()

	if err != nil {
		t.Fatalf("Failed to execute the command: %s", err)
	}

	return string(out)
}
