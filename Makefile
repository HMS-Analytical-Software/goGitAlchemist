# ------------------------------------------------------------
# run all static and unit tests
testall: vulncheck vet lint test


# ------------------------------------------------------------
# run all unit tests
test:
	go test ./...

# ------------------------------------------------------------
# vulnerability check
#
# https://go.dev/blog/vuln
# https://go.dev/blog/govulncheck
# https://go.dev/doc/tutorial/govulncheck
# https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
# Update: 
#   go install golang.org/x/vuln/cmd/govulncheck@latest
vulncheck:
	govulncheck ./...

# ------------------------------------------------------------
# static code analysis
# https://pkg.go.dev/cmd/vet
vet:
	go vet -lostcancel=false ./...

# ------------------------------------------------------------
# static code analysis
# https://golangci-lint.run/
# Update: 
#   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint:
	golangci-lint run

# ------------------------------------------------------------
# start the godoc server and the browser for development
# non exported content can be viewed by appending "/?m=all" to the url in the
# browser
docserver:
	godoc &
	xdg-open http://localhost:6060/pkg/github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/

