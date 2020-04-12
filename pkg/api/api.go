package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sh0rez/packup/pkg/restic"
)

type api struct {
	rst *restic.Restic
}

func New() http.Handler {
	var mux = http.NewServeMux()

	a := api{
		rst: restic.New(),
	}

	mux.HandleFunc("/snapshots", a.snapshotsHandler)
	mux.HandleFunc("/files", a.filesHandler)
	// mux.HandleFunc("/dump", a.dumpHandler)

	return mux
}

// filesHandler lists snapshot contents
func (a api) filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	snapshot := queryOrDefault(r, "snapshot", "latest")
	path := queryOrDefault(r, "path", "/")

	files, err := a.rst.Files(snapshot, path, false)
	if !handle500(w, err) {
		return
	}

	out, err := json.Marshal(files)
	if !handle500(w, err) {
		return
	}

	fmt.Fprint(w, string(out))
}

// snapshotsHandler lists snapshots
func (a api) snapshotsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	snapshots, err := a.rst.Snapshots()
	if !handle500(w, err) {
		return
	}

	out, err := json.Marshal(snapshots)
	if !handle500(w, err) {
		return
	}

	fmt.Fprint(w, string(out))
}

func handle500(w http.ResponseWriter, err error) bool {
	return handleErr(w, err, http.StatusInternalServerError)
}
func handleErr(w http.ResponseWriter, err error, status int) bool {
	if err == nil {
		return true
	}

	http.Error(w, err.Error(), status)
	return false
}

func queryOrDefault(r *http.Request, key, def string) string {
	str := r.URL.Query().Get(key)
	if str == "" {
		return def
	}
	return str
}
