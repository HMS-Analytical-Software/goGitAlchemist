//go:build !windows

package alchemist

// noAccessDir is used for testing access errors.
const noAccessDir = "/"

// wantTestCompleteLog represents the output of the novice in verbose mode.
var wantTestVerboseLog = `[INFO] execute formula test_workflow
[INFO] 1/14: make directory repodir/remotes/create_add_commit
[DEBUG] makedir "repodir/remotes/create_add_commit"
[INFO] 1/14: init --bare --initial-branch=main .
[DEBUG] "repodir/remotes/create_add_commit": git []string{"init", "--bare", "--initial-branch=main", "."}
[INFO] 1/14: clone remotes/create_add_commit workflow
[DEBUG] "repodir": git []string{"clone", "remotes/create_add_commit", "workflow"}
[INFO] 1/14: remote set-url origin ../remotes/create_add_commit
[DEBUG] "repodir/workflow": git []string{"remote", "set-url", "origin", "../remotes/create_add_commit"}
[INFO] 1/14: config user.name Richard Red
[DEBUG] "repodir/workflow": git []string{"config", "user.name", "Richard Red"}
[INFO] 1/14: config user.email richard@pw-compa.ny
[DEBUG] "repodir/workflow": git []string{"config", "user.email", "richard@pw-compa.ny"}
[INFO] 1/14: config init.defaultBranch main
[DEBUG] "repodir/workflow": git []string{"config", "init.defaultBranch", "main"}
[INFO] 2/14: copy cfgdir/test_workflow/files/project_plan_v1.md to repodir/workflow/project_plan.md
[DEBUG] copy "cfgdir/test_workflow/files/project_plan_v1.md" to "repodir/workflow/project_plan.md"
[INFO] 3/14: add 1 files
[DEBUG] "repodir/workflow": git []string{"add", "project_plan.md"}
[INFO] 4/14: commit
[DEBUG] "repodir/workflow": git []string{"commit", "--date=format:relative:5.hours.ago", "-m", "Added first file", "--author=Richard Red <richard@pw-compa.ny>"}
[INFO] 5/14: copy cfgdir/test_workflow/files/project_plan_v3.md to repodir/workflow/project_plan.md
[DEBUG] copy "cfgdir/test_workflow/files/project_plan_v3.md" to "repodir/workflow/project_plan.md"
[INFO] 5/14: add 1 files
[DEBUG] "repodir/workflow": git []string{"add", "."}
[INFO] 5/14: commit
[DEBUG] "repodir/workflow": git []string{"commit", "--date=format:relative:5.hours.ago", "-m", "removed unnecessary parts of the project plan", "--author=Richard Red <richard@pw-compa.ny>"}
[INFO] 6/14: copy cfgdir/test_workflow/files/folder1 to repodir/workflow/folder1
[DEBUG] copy "cfgdir/test_workflow/files/folder1" to "repodir/workflow/folder1"
[INFO] 6/14: add 1 files
[DEBUG] "repodir/workflow": git []string{"add", "."}
[INFO] 6/14: commit
[DEBUG] "repodir/workflow": git []string{"commit", "--date=format:relative:5.hours.ago", "-m", "added folder1", "--author=Richard Red <richard@pw-compa.ny>"}
[INFO] 7/14: commit -m my message
[DEBUG] "repodir/workflow": git []string{"commit", "-m", "my message"}
[INFO] 8/14: push origin main
[DEBUG] "repodir/workflow": git []string{"push", "origin", "main"}
[INFO] 9/14: merge feature/start_project with main  (delete: false)
[INFO] 9/14: checkout main
[DEBUG] "repodir/workflow": git []string{"checkout", "main"}
[INFO] 9/14: merge feature/start_project
[DEBUG] "repodir/workflow": git []string{"merge", "feature/start_project"}
[INFO] 10/14: merge feature/other with main  (delete: true)
[INFO] 10/14: checkout main
[DEBUG] "repodir/workflow": git []string{"checkout", "main"}
[INFO] 10/14: merge feature/other
[DEBUG] "repodir/workflow": git []string{"merge", "feature/other"}
[INFO] 10/14: branch -d feature/other
[DEBUG] "repodir/workflow": git []string{"branch", "-d", "feature/other"}
[INFO] 11/14: push (true)
[DEBUG] "repodir/workflow": git []string{"push", "origin", "main"}
[INFO] 12/14: mv main.py generator.py
[DEBUG] "repodir/workflow": git []string{"mv", "main.py", "generator.py"}
[INFO] 13/14: remove and commit 1 files
[INFO] 13/14: rm notes-timeline.txt
[DEBUG] "repodir/workflow": git []string{"rm", "notes-timeline.txt"}
[INFO] 13/14: commit
[DEBUG] "repodir/workflow": git []string{"commit", "--date=format:relative:5.hours.ago", "-m", "clean up timeline notes", "--author=Richard Red <richard@pw-compa.ny>"}
[INFO] 14/14: remove and commit 1 files
[INFO] 14/14: rm notes-timeline_nonexistent.txt
[DEBUG] "repodir/workflow": git []string{"rm", "notes-timeline_nonexistent.txt"}
[INFO] 14/14: commit
[DEBUG] "repodir/workflow": git []string{"commit", "--date=format:relative:5.hours.ago", "-m", "clean up timeline notes that don't exist (this step should fail)", "--author=Richard Red <richard@pw-compa.ny>"}
`

// wantTestCompleteLog represents the output of the novice in normal mode.
var wantTestCompleteLog = `[INFO] execute formula test_workflow
[INFO] 1/14: make directory repodir/remotes/create_add_commit
[INFO] 1/14: init --bare --initial-branch=main .
[INFO] 1/14: clone remotes/create_add_commit workflow
[INFO] 1/14: remote set-url origin ../remotes/create_add_commit
[INFO] 1/14: config user.name Richard Red
[INFO] 1/14: config user.email richard@pw-compa.ny
[INFO] 1/14: config init.defaultBranch main
[INFO] 2/14: copy cfgdir/test_workflow/files/project_plan_v1.md to repodir/workflow/project_plan.md
[INFO] 3/14: add 1 files
[INFO] 4/14: commit
[INFO] 5/14: copy cfgdir/test_workflow/files/project_plan_v3.md to repodir/workflow/project_plan.md
[INFO] 5/14: add 1 files
[INFO] 5/14: commit
[INFO] 6/14: copy cfgdir/test_workflow/files/folder1 to repodir/workflow/folder1
[INFO] 6/14: add 1 files
[INFO] 6/14: commit
[INFO] 7/14: commit -m my message
[INFO] 8/14: push origin main
[INFO] 9/14: merge feature/start_project with main  (delete: false)
[INFO] 9/14: checkout main
[INFO] 9/14: merge feature/start_project
[INFO] 10/14: merge feature/other with main  (delete: true)
[INFO] 10/14: checkout main
[INFO] 10/14: merge feature/other
[INFO] 10/14: branch -d feature/other
[INFO] 11/14: push (true)
[INFO] 12/14: mv main.py generator.py
[INFO] 13/14: remove and commit 1 files
[INFO] 13/14: rm notes-timeline.txt
[INFO] 13/14: commit
[INFO] 14/14: remove and commit 1 files
[INFO] 14/14: rm notes-timeline_nonexistent.txt
[INFO] 14/14: commit
`

// wantTestShortLog represents the output of the novice when only two
// steps are executed.
var wantTestShortLog = `[INFO] execute formula test_workflow
[INFO] 1/14: make directory repodir/remotes/create_add_commit
[INFO] 1/14: init --bare --initial-branch=main .
[INFO] 1/14: clone remotes/create_add_commit workflow
[INFO] 1/14: remote set-url origin ../remotes/create_add_commit
[INFO] 1/14: config user.name Richard Red
[INFO] 1/14: config user.email richard@pw-compa.ny
[INFO] 1/14: config init.defaultBranch main
[INFO] 2/14: copy cfgdir/test_workflow/files/project_plan_v1.md to repodir/workflow/project_plan.md
`

// wantNoviceLog is the log output from the novice.
var wantNoviceLog = `[DEBUG] copy "testdata/source.txt" to "testdata/new_page/new_dir"
`
