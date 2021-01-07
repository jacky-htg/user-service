package databasetest

import (
	"bytes"
	"os/exec"
	"testing"
)

// StartContainer runs a mysql container to execute commands.
func StartContainer(t *testing.T) {
	t.Helper()

	cmd := exec.Command("docker", "run", "-d", "--name", "postgres_test", "--publish", "54320:5432", "--env", "POSTGRES_PASSWORD=1234", "--env", "PGDATA=/var/lib/postgresql/data/pgdata", "postgres:11-alpine")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start docker : %v", err)
	}

}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T) {
	t.Helper()

	if err := exec.Command("docker", "container", "rm", "-f", "postgres_test").Run(); err != nil {
		t.Fatalf("could not stop mysql container: %v", err)
	}
}
