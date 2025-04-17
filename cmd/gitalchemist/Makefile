# ------------------------------------------------------------
# globals

BINARYNAME=gitalchemist



# ------------------------------------------------------------
# unit test
test:
	go test

# verbose unit test
testv:
	go test -v


# test coverage
COVERAGE_TMP_FILE=cover.txt
COVERAGE_HTML_FILE=cover.html

cover:
	go test -cover -coverprofile $(COVERAGE_TMP_FILE) $(TIMEOUT)
	go tool cover -func=$(COVERAGE_TMP_FILE)
	go tool cover -html=$(COVERAGE_TMP_FILE) -o $(COVERAGE_HTML_FILE)
	echo "chromium-browser $(COVERAGE_HTML_FILE) &"
	echo "firefox $(COVERAGE_HTML_FILE) &"

.PHONY: test testv cover 

# ------------------------------------------------------------
# acceptance tests
# These tests require that the binary compiles successfully.
# It is also used for cleaning the output directory before
# exectution of the tests.

# default is basic_workflow, specify other on make call.
# example: make acctest TASK=cmd_mv
TASK?=basic_workflow
TESTTAG=-tags acctest

# test one task
acctest: clean
	go test $(TESTTAG) -run TestAcceptanceTask/$(TASK)

# test one task verbose
acctestv: clean
	go test -v $(TESTTAG) -run TestAcceptanceTask/$(TASK)

# test all tasks
acctestall: clean
	go test $(TESTTAG) -run TestAcceptanceTask

# test all tasks verbose
acctestallv: clean
	go test -v $(TESTTAG) -run TestAcceptanceTask

# test runall
acctestrunall: clean
	go test $(TESTTAG) -run TestAcceptanceRunAll

# test runall verbose
acctestrunallv: clean
	go test -v $(TESTTAG) -run TestAcceptanceRunAll

.PHONY: acctest acctestv acctestall acctestallv acctestrunall acctestrunallv

# run github workflow in batch mode in current branch
TESTOS?=linux
CURRENTBRANCH=$(shell git  branch --show-current)
WORKFLOW=gitalchemist-acceptance-tests 
acctestworkflow:
	gh workflow run $(WORKFLOW) --ref $(CURRENTBRANCH) \
		-f os=$(TESTOS)

.PHONY: acctestworkflow

# ------------------------------------------------------------
# remove temporary test output
clean: build
	./$(BINARYNAME) -clean

cleanlinux:
	rm -rf cwd

.PHONY: clean cleanlinux

# ------------------------------------------------------------
# build

VERSIONTAG=$(shell git describe --tags)
BUILDFLAGS=-ldflags="-X 'main.Version=$(VERSIONTAG)'"

# local development
build:
	go build $(BUILDFLAGS)

# build all os variants
all: win linux mac

# build windows exe for amd64
win:
	GOOS=windows GOARCH=amd64 go build -o $(BINARYNAME).exe $(BUILDFLAGS)

# build linux binary for amd64
linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARYNAME) $(BUILDFLAGS)

# build macosx binary for arm64
mac:
	GOOS=darwin GOARCH=arm64 go build -o $(BINARYNAME).mac $(BUILDFLAGS)

.PHONY: build all win linux mac

# ------------------------------------------------------------
# make github release

RELEASENAME=$(BINARYNAME)-$(REL)
GIT_HOST=github.com
REPOPATH=https://github.com/HMS-Analytical-Software/goGitAlchemist

release:
ifeq ($(strip $(REL)),)
	echo "call make release REL=<release>"
	echo "example: make release REL=v0.0.1"
	exit 2
else
	git tag $(RELEASENAME)
	git push --tags
	$(MAKE) all
	gh release create $(RELEASENAME) \
		--notes "$(RELEASENAME)" \
		--title "$(RELEASENAME)" \
		--repo $(REPOPATH) \
		$(BINARYNAME) \
		$(BINARYNAME).mac \
		$(BINARYNAME).exe 
endif

.PHONY: release

