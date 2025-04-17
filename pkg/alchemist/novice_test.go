package alchemist

import (
	"bytes"
	"log"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNovice(t *testing.T) {
	dir, from, to := "dir", "from", "to"
	var buf bytes.Buffer
	novice := newNovice(log.New(&buf, "", 0), Options{Verbose: true})

	err := novice.git(dir, "init")
	if err != nil {
		t.Errorf("ERROR: got error: %v", err)
	}
	err = novice.copy(from, to)
	if err != nil {
		t.Errorf("ERROR: got error: %v", err)
	}
	err = novice.makedir(dir)
	if err != nil {
		t.Errorf("ERROR: got error: %v", err)
	}

	want := `[DEBUG] "dir": git []string{"init"}
[DEBUG] copy "from" to "to"
[DEBUG] makedir "dir"
`
	got := buf.String()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ERROR: got- want+\n%s\n", diff)
	}

	// copy with makedir
	buf.Reset()
	err = novice.copy(
		filepath.Join(TestDataDir, "source.txt"),
		filepath.Join(TestDataDir, "new_page", "new_dir"))
	if err != nil {
		t.Errorf("ERROR: got error: %v", err)
	}
	got = buf.String()
	want = wantNoviceLog

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ERROR: got- want+\n%s\n", diff)
	}
}
