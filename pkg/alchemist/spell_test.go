package alchemist

import (
	"path/filepath"
	"testing"

	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
)

// TestSpellCast tests casting different spells using the assistantSpy
// test double.
func TestSpellCast(t *testing.T) {

	// values used in multiple places
	bareDir, cloneDir, repoDir, newDir := "baredir", "clonedir", "repodir", "newdir"
	fromFile, toFile := "from.txt", filepath.Join(newDir, "to.txt")
	filePair1 := fromFile + " => " + fromFile
	filePair2 := "to.txt => " + toFile

	repoBareDir := filepath.Join(repoDir, bareDir)
	repoCloneDir := filepath.Join(repoDir, cloneDir)

	testCases := []struct {
		name    string        // test case name
		spell   caster        // type under test
		spy     *assistantSpy // test double assistant, records the calls
		want    [][]string    // wanted call recordings by the spy
		wantErr string        // wanted error message ("" if no error is expected)
	}{{
		name:  "initRepoSpell ok",
		spell: initRepoSpell{Bare: bareDir, CloneTo: cloneDir},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{"makedir", repoBareDir},
			[]string{repoBareDir, gitCmd, "init", "--bare", "--initial-branch=main", "."},
			[]string{repoDir, gitCmd, "clone", bareDir, cloneDir},
			[]string{repoCloneDir, gitCmd, "remote", "set-url", "origin", filepath.Join("..", bareDir)},
			[]string{repoCloneDir, gitCmd, "config", "user.name", author[defaultUser]},
			[]string{repoCloneDir, gitCmd, "config", "user.email", email[defaultUser]},
			[]string{repoCloneDir, gitCmd, "config", "init.defaultBranch", defaultBranch},
		},
	}, {
		name:    "initRepoSpell error 1",
		spell:   initRepoSpell{Bare: bareDir, CloneTo: cloneDir},
		spy:     &assistantSpy{errorAt: 1}, // return error at the first call
		wantErr: "make dir " + repoBareDir + ": spy error: 1",
	}, {
		name:  "initRepoSpell error 2",
		spell: initRepoSpell{Bare: bareDir, CloneTo: cloneDir},
		spy:   &assistantSpy{errorAt: 2}, // return error at the second call
		want: [][]string{
			[]string{"makedir", repoBareDir},
		},
		wantErr: gitCmd + " init --bare --initial-branch=main .: spy error: 2",
	}, {
		name:  "initRepoSpell error 3",
		spell: initRepoSpell{Bare: bareDir, CloneTo: cloneDir},
		spy:   &assistantSpy{errorAt: 3},
		want: [][]string{
			[]string{"makedir", repoBareDir},
			[]string{repoBareDir, gitCmd, "init", "--bare", "--initial-branch=main", "."},
		},
		wantErr: gitCmd + " clone " + bareDir + " " + cloneDir + ": spy error: 3",
	}, {
		name:  "createFileSpell ok",
		spell: createFileSpell{Source: fromFile, Target: toFile},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{"copy", fromFile, filepath.Join(repoDir, toFile)},
		},
	}, {
		name:  "createFileSpell error",
		spell: createFileSpell{Source: fromFile, Target: toFile},
		spy:   &assistantSpy{errorAt: 1},
		wantErr: "copy " + fromFile + " " + filepath.Join(repoDir, toFile) +
			": spy error: 1",
	}, {
		name:  "addSpell ok",
		spell: addSpell{Files: []string{fromFile, toFile}},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "add", fromFile},
			[]string{repoDir, gitCmd, "add", toFile},
		},
	}, {
		name:  "addSpell no files",
		spell: addSpell{},
		spy:   &assistantSpy{},
	}, {
		name:  "addSpell error second",
		spell: addSpell{Files: []string{fromFile, toFile}},
		spy:   &assistantSpy{errorAt: 2},
		want: [][]string{
			[]string{repoDir, gitCmd, "add", fromFile},
		},
		wantErr: gitCmd + " add " + toFile + ": spy error: 2",
	}, {
		name:  "commitSpell ok",
		spell: commitSpell{Author: "red", Message: "hello"},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "commit", "--date=" + gitCommitDateFormat,
				"-m", "hello", "--author=" + author["red"] + " <" + email["red"] + ">"},
		},
	}, {
		name:  "commitSpell unknown author",
		spell: commitSpell{Author: "skywalker", Message: "hello"},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "commit", "--date=" + gitCommitDateFormat,
				"-m", "hello", "--author=skywalker"},
		},
	}, {
		name:  "commitSpell error",
		spell: commitSpell{Author: "skywalker", Message: "hello"},
		spy:   &assistantSpy{errorAt: 1},
		wantErr: gitCmd + " commit --date=" + gitCommitDateFormat +
			" -m hello --author=skywalker: spy error: 1",
	}, {
		name: "createAddCommitSpell ok",
		spell: createAddCommitSpell{
			Files:   []string{filePair1, filePair2},
			Author:  "red",
			Message: "hello",
		},
		spy: &assistantSpy{},
		want: [][]string{
			[]string{"copy", fromFile, filepath.Join(repoDir, fromFile)},
			[]string{"copy", "to.txt", filepath.Join(repoDir, toFile)},
			[]string{repoDir, gitCmd, "add", "."},
			[]string{repoDir, gitCmd, "commit", "--date=" + gitCommitDateFormat,
				"-m", "hello", "--author=" + author["red"] + " <" + email["red"] + ">"},
		},
	}, {
		name: "createAddCommitSpell error",
		spell: createAddCommitSpell{
			Files:   []string{filePair1, filePair2},
			Author:  "red",
			Message: "hello",
		},
		spy: &assistantSpy{errorAt: 1},
		wantErr: "copy " + fromFile + " " + filepath.Join(repoDir, fromFile) +
			": spy error: 1",
	}, {
		name:  "gitSpell ok",
		spell: gitSpell{Command: `add -m "my msg" .`},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "add", "-m", "my msg", "."},
		},
	}, {
		name:  "gitSpell git prefix",
		spell: gitSpell{Command: gitCmd + ` add -m "my msg" .`},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "add", "-m", "my msg", "."},
		},
	}, {
		name:    "gitSpell error",
		spell:   gitSpell{Command: `add  -m "my msg" .`},
		spy:     &assistantSpy{errorAt: 1},
		wantErr: gitCmd + " add -m my msg .: spy error: 1",
	}, {
		name: "mergeSpell ok",
		spell: mergeSpell{
			Source:       "develop",
			Target:       "main",
			DeleteSource: true,
		},
		spy: &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "checkout", "main"},
			[]string{repoDir, gitCmd, "merge", "develop"},
			[]string{repoDir, gitCmd, "branch", "-d", "develop"},
		},
	}, {
		name: "mergeSpell no delete",
		spell: mergeSpell{
			Source: "develop",
			Target: "main",
		},
		spy: &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "checkout", "main"},
			[]string{repoDir, gitCmd, "merge", "develop"},
		},
	}, {
		name: "mergeSpell checkout error",
		spell: mergeSpell{
			Source:       "develop",
			Target:       "main",
			DeleteSource: true,
		},
		spy:     &assistantSpy{errorAt: 1},
		wantErr: gitCmd + " checkout main: spy error: 1",
	}, {
		name: "mergeSpell merge error",
		spell: mergeSpell{
			Source:       "develop",
			Target:       "main",
			DeleteSource: true,
		},
		spy: &assistantSpy{errorAt: 2},
		want: [][]string{
			[]string{repoDir, gitCmd, "checkout", "main"},
		},
		wantErr: gitCmd + " merge develop: spy error: 2",
	}, {
		name: "mergeSpell delete error",
		spell: mergeSpell{
			Source:       "develop",
			Target:       "main",
			DeleteSource: true,
		},
		spy: &assistantSpy{errorAt: 3},
		want: [][]string{
			[]string{repoDir, gitCmd, "checkout", "main"},
			[]string{repoDir, gitCmd, "merge", "develop"},
		},
		wantErr: gitCmd + " branch -d develop: spy error: 3",
	}, {
		name:  "pushSpell ok",
		spell: pushSpell{Main: true},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "push", "origin", "main"},
		},
	}, {
		name:  "pushSpell no main",
		spell: pushSpell{},
		spy:   &assistantSpy{},
	}, {
		name:    "pushSpell error",
		spell:   pushSpell{Main: true},
		spy:     &assistantSpy{errorAt: 1},
		wantErr: gitCmd + " push origin main: spy error: 1",
	}, {
		name:  "moveSpell ok",
		spell: moveSpell{Source: fromFile, Target: toFile},
		spy:   &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "mv", fromFile, toFile},
		},
	}, {
		name:    "moveSpell error",
		spell:   moveSpell{Source: fromFile, Target: toFile},
		spy:     &assistantSpy{errorAt: 1},
		wantErr: gitCmd + " mv " + fromFile + " " + toFile + ": spy error: 1",
	}, {
		name: "removeAndCommitSpell ok",
		spell: removeAndCommitSpell{
			Files:   []string{fromFile, toFile},
			Author:  "red",
			Message: "hello",
		},
		spy: &assistantSpy{},
		want: [][]string{
			[]string{repoDir, gitCmd, "rm", fromFile},
			[]string{repoDir, gitCmd, "rm", toFile},
			[]string{repoDir, gitCmd, "commit", "--date=" + gitCommitDateFormat,
				"-m", "hello", "--author=" + author["red"] + " <" + email["red"] + ">"},
		},
	}, {
		name: "removeAndCommitSpell rm error",
		spell: removeAndCommitSpell{
			Files:   []string{fromFile, toFile},
			Author:  "red",
			Message: "hello",
		},
		spy: &assistantSpy{errorAt: 2},
		want: [][]string{
			[]string{repoDir, gitCmd, "rm", fromFile},
		},
		wantErr: gitCmd + " rm " + toFile + ": spy error: 2",
	}, {
		name: "removeAndCommitSpell commit error",
		spell: removeAndCommitSpell{
			Files:   []string{fromFile, toFile},
			Author:  "red",
			Message: "hello",
		},
		spy: &assistantSpy{errorAt: 3},
		want: [][]string{
			[]string{repoDir, gitCmd, "rm", fromFile},
			[]string{repoDir, gitCmd, "rm", toFile},
		},
		wantErr: gitCmd + " commit --date=" + gitCommitDateFormat +
			" -m hello --author=" + getAuthor("red") + ": spy error: 3",
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			err := c.spell.cast(c.spy, Options{RepoDir: repoDir})

			got := c.spy.calls
			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("ERROR: got-, want+\n%v\n", diff)
			}

			// check error after checking the calls to prevent skipping
			// when error is wanted
			check.ErrorString(t, err, c.wantErr)
		})
	}
}
