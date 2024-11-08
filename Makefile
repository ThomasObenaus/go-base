#.DEFAULT_GOAL := all
name := "go-base"

all: tools test finish

unit_test_packages := $(shell go list ./... | grep -v "/test/mocks")

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"


test: sep gen-mocks ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ${unit_test_packages} -covermode=count -coverprofile=coverage.out

cover-upload: sep ## Uploads the unittest coverage to coveralls (for this the GO_BASE_COVERALLS_REPO_TOKEN has to be set correctly).
	# for this to get working you have to export the repo_token for your repo at coveralls.io
	# i.e. export GO_BASE_COVERALLS_REPO_TOKEN=<your token>
	@${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=circleci -repotoken=${GO_BASE_COVERALLS_REPO_TOKEN}

gen-mocks: sep ## Generates test doubles (mocks).
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go install github.com/golang/mock/mockgen@latest
	@mockgen -source=health/check.go -destination test/mocks/health/mock_check.go
	@mockgen -source=stop/interfaces.go -destination stop/mock_stop_test.go -package stop
	@mockgen -source=stop/interfaces.go -destination shutdown/mock_stop_test.go -package shutdown
	@mockgen -source=shutdown/interfaces.go -destination shutdown/mock_interfaces_test.go -package shutdown
	@mockgen -source=signal/signal.go -destination signal/mock_signal_test.go -package signal

tools: sep ## Installs needed tools
	@echo "--> Install needed tools"
	@go install golang.org/x/tools/cmd/cover@latest
	@go install github.com/mattn/goveralls@latest

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="