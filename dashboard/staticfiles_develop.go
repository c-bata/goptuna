// +build develop

package dashboard

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func registerStaticFileRoutes(r *mux.Router, prefix string) {
	staticRoot := os.Getenv("GOPTUNA_DASHBOARD_STATIC_ROOT")
	if staticRoot == "" {
		staticRoot = "dashboard/public"
	}

	wf := func(filepath string, info os.FileInfo, err error) error {
		if filepath == prefix {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(filepath, ".png") {
			return nil
		}
		urlpath := filepath[len(staticRoot):]
		urlpath = path.Join(prefix, urlpath)

		var contentType string
		if strings.HasSuffix(filepath, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(filepath, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(filepath, ".woff") {
			contentType = "application/font-woff"
		} else if strings.HasSuffix(filepath, ".ttf") {
			contentType = "application/x-font-ttf"
		} else if strings.HasSuffix(filepath, ".otf") {
			contentType = "application/x-font-otf"
		} else if strings.HasSuffix(filepath, ".svg") {
			contentType = "image/svg+xml"
		} else if strings.HasSuffix(filepath, ".eot") {
			contentType = "image/vnd.ms-fontobject"
		}
		handler := func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadFile(filepath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Internal server error")
				return
			}
			w.Header().Set("Content-Type", contentType)
			w.Write(data)
		}
		r.Path(urlpath).Methods("Get").HandlerFunc(handler)
		return nil
	}
	err := filepath.Walk(staticRoot, wf)
	if err != nil {
		panic(err)
	}
}
