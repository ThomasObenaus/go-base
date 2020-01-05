#.DEFAULT_GOAL := all
name := "go-base"

all: test finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"


test: sep ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ./health ./shutdown ./config ./buildinfo ./logging -covermode=count -coverprofile=coverage.out

gen-mocks: sep ## Generates test doubles (mocks).
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@mockgen -source=health/check.go -destination test/mocks/health/mock_check.go

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="