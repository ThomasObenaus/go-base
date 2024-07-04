#.DEFAULT_GOAL := all
name := "go-base"

all: tools test finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"


test: sep gen-mocks ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ./health ./shutdown ./shutdown/v2/*/ ./config ./buildinfo ./logging -covermode=count -coverprofile=coverage.out

cover-upload: sep ## Uploads the unittest coverage to coveralls (for this the GO_BASE_COVERALLS_REPO_TOKEN has to be set correctly).
	# for this to get working you have to export the repo_token for your repo at coveralls.io
	# i.e. export GO_BASE_COVERALLS_REPO_TOKEN=<your token>
	@${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=circleci -repotoken=${GO_BASE_COVERALLS_REPO_TOKEN}

gen-mocks: sep ## Generates test doubles (mocks).
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go install github.com/golang/mock/mockgen@latest
	@mockgen -source=health/check.go -destination test/mocks/health/mock_check.go
	@mockgen -source=shutdown/stopable.go -destination shutdown/mock_stopable_test.go -package shutdown
	@mockgen -source=shutdown/stopable.go -destination test/mocks/shutdown/mock_stopable.go
	@mockgen -source=shutdown/shutdownHandler.go -destination test/mocks/shutdown/mock_shutdownHandler.go
	@mockgen -source=shutdown/v2/stop/interfaces.go -destination shutdown/v2/stop/mock_stop_test.go -package stop
	@mockgen -source=shutdown/v2/stop/interfaces.go -destination shutdown/v2/mock_stop_test.go -package v2
	@mockgen -source=shutdown/v2/signal/signal.go -destination shutdown/v2/signal/mock_signal_test.go -package signal
	@mockgen -source=shutdown/v2/interfaces.go -destination shutdown/v2/mock_interfaces_test.go -package v2

tools: sep ## Installs needed tools
	@echo "--> Install needed tools"
	@go install golang.org/x/tools/cmd/cover@latest
	@go install github.com/mattn/goveralls@latest

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="