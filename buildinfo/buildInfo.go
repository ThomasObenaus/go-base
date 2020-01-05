package buildinfo

import (
	"encoding/json"
	"net/http"
	"strings"
)

// BuildInfo contains information about the build.
type BuildInfo struct {
	Version   string `json:"version,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	Revision  string `json:"revision,omitempty"`
	Branch    string `json:"branch,omitempty"`
}

// Print prints the build information using the given print function
func (bi *BuildInfo) Print(printFun func(format string, a ...interface{}) (n int, err error)) {

	version := strings.TrimSpace(bi.Version)
	if len(version) == 0 {
		version = "n/a"
	}

	buildTime := strings.TrimSpace(bi.BuildTime)
	if len(buildTime) == 0 {
		buildTime = "n/a"
	}

	revision := strings.TrimSpace(bi.Revision)
	if len(revision) == 0 {
		revision = "n/a"
	}

	branch := strings.TrimSpace(bi.Branch)
	if len(branch) == 0 {
		branch = "n/a"
	}

	printFun("-----------------------------------------------------------------\n")
	printFun("BuildInfo\n")
	printFun("-----------------------------------------------------------------\n")
	printFun("\tVersion:\t%s\n", version)
	printFun("\tBuild-Time:\t%s\n", buildTime)
	printFun("\tRevision:\t%s on %s\n", revision, branch)
	printFun("-----------------------------------------------------------------\n")
}

// BuildInfo represents the build-info end-point of sokar
func (bi *BuildInfo) BuildInfo(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(bi); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
