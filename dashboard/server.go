package dashboard

import (
	"encoding/json"
	"errors"
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

	errNotFound = errors.New("not found")
)

func NewServer(s goptuna.Storage) (http.Handler, error) {
	storageMutex.Lock()
	defer storageMutex.Unlock()
	storage = s

	router := mux.NewRouter()

	// Redirect to /dashboard for react-router BrowserRouter
	router.Handle("/", http.RedirectHandler("/dashboard", http.StatusFound)).Methods("GET")
	router.PathPrefix("/dashboard").HandlerFunc(handleGetIndex).Methods("GET")

	// JSON API
	router.HandleFunc("/api/studies", handleGetAllStudySummary).Methods("GET")
	router.HandleFunc("/api/studies", handleCreateStudy).Methods("POST")
	router.HandleFunc("/api/studies/{study_id:[0-9]+}", handleGetStudyDetail).Methods("GET")

	// Static files
	err := registerStaticFileRoutes(router, "static")
	if err != nil {
		return nil, err
	}
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
    <title>Goptuna Dashboard</title>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
    body {
        min-height: 100vh;
        margin: 0;
        padding: 0;
    }
    
    h1, h2, h3 {
        font-weight: 600;
        letter-spacing: 1px;
        line-height: 1.3;
    }
    </style>
    <script defer src="/static/bundle.js"></script>
</head>

<body>
    <noscript>You need to enable JavaScript to run this dashboard.</noscript>
    <div id="dashboard">
         <p>Now loading...</p>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
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

func handleCreateStudy(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"study_name"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		// TODO(c-bata): Return bad request if study already exist
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	studyID, err := storage.CreateNewStudy(req.Name)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	studySummary, err := getStudySummary(studyID)
	if err != errNotFound {
		writeErrorResponse(w, http.StatusNotFound, "Not found")
		return
	} else if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(struct {
		StudySummary StudySummary `json:"study_summary"`
	}{
		StudySummary: studySummary,
	})
}

func handleGetStudyDetail(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	studyID, err := strconv.Atoi(urlVars["study_id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid study id")
		return
	}

	studySummary, err := getStudySummary(studyID)
	if err != errNotFound {
		writeErrorResponse(w, http.StatusNotFound, "Not found")
		return
	} else if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	trials, err := storage.GetAllTrials(studyID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(struct {
		Name          string        `json:"name"`
		Direction     string        `json:"direction"`
		DatetimeStart string        `json:"datetime_start"`
		BestTrial     FrozenTrial   `json:"best_trial"`
		Trials        []FrozenTrial `json:"trials"`
	}{
		Name:          studySummary.Name,
		Direction:     studySummary.Direction,
		DatetimeStart: studySummary.DatetimeStart,
		BestTrial:     studySummary.BestTrial,
		Trials:        toFrozenTrials(trials),
	})
}

func getStudySummary(studyID int) (StudySummary, error) {
	studies, err := storage.GetAllStudySummaries()
	if err != nil {
		return StudySummary{}, err
	}
	for _, s := range studies {
		if s.ID == studyID {
			return toStudySummary(s), nil
		}
	}
	return StudySummary{}, errNotFound
}
