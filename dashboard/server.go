package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/c-bata/goptuna"
	"github.com/gorilla/mux"
)

var (
	storage      goptuna.Storage
	storageMutex sync.RWMutex
)

func NewServer(s goptuna.Storage) (http.Handler, error) {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	storage = s

	router := mux.NewRouter()
	// HTML
	router.HandleFunc("/", handleGetIndex).Methods("GET")
	// Static files
	err := registerStaticFileRoutes(router, "static")
	if err != nil {
		return nil, err
	}
	// JSON API
	router.HandleFunc("/api/studies", handleGetAllStudySummary).Methods("GET")
	router.HandleFunc("/api/studies/{study_id:[0-9]+}/trials", handleGetTrials).Methods("GET")
	return router, nil
}

func writeErrorResponse(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(struct {
		Reason string `json:"reason"`
	}{
		Reason: reason,
	})
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	htmlStr := `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Goptuna Dashboard</title>
    <script async src="static/bundle.js"></script>
</head>

<body>
    <noscript>You need to enable JavaScript to run this dashboard.</noscript>
    <div id="dashboard"></div>
</body>
</html>
`
	fmt.Fprintf(w, htmlStr)
}

func handleGetAllStudySummary(w http.ResponseWriter, r *http.Request) {
	studies, err := storage.GetAllStudySummaries()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(struct {
		StudySummaries []StudySummary `json:"study_summaries"`
	}{
		StudySummaries: toStudySummaries(studies),
	})
}

func handleGetTrials(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	studyID, err := strconv.Atoi(urlVars["study_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid study id")
		return
	}
	trials, err := storage.GetAllTrials(studyID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(struct {
		Trials []FrozenTrial `json:"trials"`
	}{
		Trials: toFrozenTrials(trials),
	})
}
