//go:build teststartdocker

package alchemist

import (
	"go/build"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"testing"
)

const (
	dockerCmd    = "docker"
	sleepSeconds = "300"
)

// TestStartDocker starts a docker container, modifies it and
// executes go test within the container.
func TestStartDocker(t *testing.T) {
	containerID, userName := startDocker(t)
	if containerID == "" {
		t.Fatalf("ERROR: got no container id")
	}

	args := []string{
		"exec", "-u", userName, containerID,
		"go", "test", "-tags", "testrundocker", "-run", "TestDocker",
	}
	if testing.Verbose() {
		args = append(args, "-v")
	}
	t.Logf("INFO: executing %s %v", dockerCmd, args)

	cmd := exec.Command(dockerCmd, args...)
	if testing.Verbose() {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}

	t.Logf("INFO: docker exec -it %s /bin/bash", containerID)

	err := cmd.Run()
	if err != nil {
		t.Fatalf("ERROR: got error %v", err)
	}

	if !testing.Verbose() {
		err = exec.Command(dockerCmd, "stop", containerID).Run()
		if err != nil {
			t.Fatalf("ERROR: stop docker container: %v", err)
		}
	}
}

// startDocker starts a golang docker image that fits the
// currently used go verison and creates a group and a user
// for non-privileged test runs.
func startDocker(t *testing.T) (string, string) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("ERROR: get working directory: %v ", err)
	}

	gopath := os.Getenv("$GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	version := strings.TrimPrefix(runtime.Version(), "go")

	output, err := exec.Command(dockerCmd,
		"run", "-d", "--rm",
		"--mount", "type=bind,src="+wd+"/../..,dst=/gitalchemist",
		"--mount", "type=bind,src="+gopath+",dst=/go",
		"-w", "/gitalchemist/pkg/alchemist",
		"golang:"+version, "sleep", sleepSeconds,
	).CombinedOutput()

	if err != nil {
		t.Logf("INFO: got %s", output)
		t.Fatalf("ERROR: start container: %v ", err)
	}

	containerID := strings.TrimSpace(string(output))

	my, err := user.Current()
	if err != nil {
		t.Fatalf("ERROR: determine current user: %v", err)
	}

	cmd := exec.Command(dockerCmd, "exec", containerID,
		"useradd", "-m", "-u", my.Uid, my.Username)
	if testing.Verbose() {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		t.Fatalf("ERROR: create user in container: %v", err)
	}

	return containerID, my.Username
}
