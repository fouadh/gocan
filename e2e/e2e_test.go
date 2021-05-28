package e2e

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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
}

func startDatabase(t *testing.T) {
	output := startCommand(t, 10 * time.Second, "start-db")
	if strings.Contains(output, "database system is ready to accept connections") {
		t.Logf("%s Database started", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Database could not be started", failed)
	}
}

func setupDatabase(t *testing.T) (string) {
	dir, err := ioutil.TempDir("", "gocan")
	if err != nil {
		t.Fatalf("Cannot create temp directory: %s", err.Error())
	}
	output := runCommand(t, "setup-db", "--path", dir)
	if strings.Contains(output, "Database configured") {
		t.Logf("%s Database configured", succeed)
	} else {
		t.Log(output)
		t.Fatalf("%s Database configuration failed", failed)
	}
	return dir
}

func startCommand(t *testing.T, sleep time.Duration, args ...string) string {
	command := exec.Command("../bin/gocan", args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	command.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	command.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := command.Start()

	if err != nil {
		t.Logf("Failed to execute the command: %s", err)
	}

	time.Sleep(sleep)

	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	t.Logf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	return outStr
}

func runCommand(t *testing.T, args ...string) string {
	command := exec.Command("../bin/gocan", args...)

	out, err := command.CombinedOutput()
	if err != nil {
		t.Logf("Failed to execute the command: %s", err)
	}

	return string(out)
}