# ------------------------------------------------------------
#  globals

COVERAGE_TMP_FILE=cover.txt
COVERAGE_HTML_FILE=cover.html

.PHONY: test testv cover

# ------------------------------------------------------------
# run unit test
test:
	go test

# ------------------------------------------------------------
# run verbose unit test
testv:
	go test -v

# ------------------------------------------------------------
# check test coverage
# special case here: verbose mode required	

cover:
	go test -v -cover -coverprofile $(COVERAGE_TMP_FILE) $(TIMEOUT)
	go tool cover -func=$(COVERAGE_TMP_FILE)
	go tool cover -html=$(COVERAGE_TMP_FILE) -o $(COVERAGE_HTML_FILE)
	echo "chromium-browser $(COVERAGE_HTML_FILE) &"
	echo "firefox $(COVERAGE_HTML_FILE) &"


