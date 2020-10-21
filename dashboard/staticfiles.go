// +build !develop

package dashboard

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	//nolint: golint
	_ "github.com/c-bata/goptuna/dashboard/statik"
)

func registerStaticFileRoutes(r *mux.Router, prefix string) error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	files := []string{
		"/bundle.js",
		"/bundle.js.LICENSE.txt",
	}
	for _, filepath := range files {
		file, err := statikFS.Open(filepath)
		if err != nil {
			return err
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		var contentType string
		if strings.HasSuffix(filepath, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(filepath, ".js") {
			contentType = "application/javascript"
		} else {
			contentType = "text/plain"
		}
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			w.Write(data)
		}
		urlpath := filepath[len("/"):]
		urlpath = path.Join(prefix, urlpath)
		if urlpath[0] != '/' {
			urlpath = "/" + urlpath
		}
		r.Path(urlpath).Methods("Get").HandlerFunc(handler)
	}
	return nil
}
