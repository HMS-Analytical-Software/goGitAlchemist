# ------------------------------------------------------------
#  globals

COVERAGE_TMP_FILE=cover.txt
COVERAGE_HTML_FILE=cover.html

.PHONY: test testv cover dockertest dockertestv

# ------------------------------------------------------------
# run unit test
test:
	go test

# ------------------------------------------------------------
# run verbose unit test
testv:
	go test -v

# ------------------------------------------------------------
# run special tests in docker container
TESTNAME=-run TestStartDocker
TESTTAGS=-tags teststartdocker

dockertest:
	go test $(TESTTAGS) $(TESTNAME)

dockertestv:
	go test -v $(TESTTAGS) $(TESTNAME)

# ------------------------------------------------------------
# check test coverage
	
cover:
	go test -cover -coverprofile $(COVERAGE_TMP_FILE) $(TIMEOUT)
	go tool cover -func=$(COVERAGE_TMP_FILE)
	go tool cover -html=$(COVERAGE_TMP_FILE) -o $(COVERAGE_HTML_FILE)
	echo "chromium-browser $(COVERAGE_HTML_FILE) &"
	echo "firefox $(COVERAGE_HTML_FILE) &"


