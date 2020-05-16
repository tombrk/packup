package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/sh0rez/packup/pkg/config"
	"github.com/sh0rez/packup/pkg/restic"
)

type api struct {
	jobs config.Jobs
}

func New(jobs config.Jobs) (http.Handler, error) {
	var r = mux.NewRouter()

	a := api{
		jobs: jobs,
	}

	r.HandleFunc("/jobs", a.jobsHandler).Methods("GET")
	r.HandleFunc("/jobs/{job}/snapshots", a.snapshotsHandler).Methods("GET")
	r.HandleFunc("/jobs/{job}/files", a.filesHandler).Methods("GET")
	r.HandleFunc("/jobs/{job}/dump", a.dumpHandler).Methods("GET")

	return r, nil
}

func (a api) Job(r *http.Request) (*restic.Restic, error) {
	name, ok := mux.Vars(r)["job"]
	if !ok {
		return nil, fmt.Errorf("URL path field 'job' missing from request")
	}

	job, ok := a.jobs[name]
	if !ok {
		return nil, fmt.Errorf("No job named '%s'. Please check your config", name)
	}

	return restic.New(job.Repo, job.Password), nil
}

func (a api) jobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := json.Marshal(a.jobs)
	if !handle500(w, err) {
		return
	}

	fmt.Fprint(w, string(data))
}

// filesHandler lists snapshot contents
func (a api) filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	snapshot := queryOrDefault(r, "snapshot", "latest")
	path := queryOrDefault(r, "path", "/")

	rst, err := a.Job(r)
	if !handleErr(w, err, http.StatusNotFound) {
		return
	}

	files, err := rst.Files(snapshot, path, false)
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

	rst, err := a.Job(r)
	if !handleErr(w, err, http.StatusNotFound) {
		return
	}

	snapshots, err := rst.Snapshots()
	if !handle500(w, err) {
		return
	}

	out, err := json.Marshal(snapshots)
	if !handle500(w, err) {
		return
	}

	fmt.Fprint(w, string(out))
}

func (a api) dumpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rst, err := a.Job(r)
	if !handleErr(w, err, http.StatusNotFound) {
		return
	}

	snapshot := queryOrDefault(r, "snapshot", "latest")
	dir := queryOrDefault(r, "path", "/")
	compress := r.URL.Query().Get("compress") == "true"

	name := filepath.Base(dir)
	if compress {
		name += ".tar.gz"
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	var dw io.Writer = w
	if compress {
		zw := gzip.NewWriter(w)
		defer zw.Close()
		dw = zw
	}

	err = rst.Dump(dw, snapshot, dir)
	if !handle500(w, err) {
		return
	}

	return
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
