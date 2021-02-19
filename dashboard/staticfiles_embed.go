// +build !develop

package dashboard

import (
	_ "embed"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

//go:embed public/bundle.js
var jsbytes []byte

func registerStaticFileRoutes(r *mux.Router, prefix string) error {
	r.Path(path.Join(prefix, "bundle.js")).Methods("Get").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(jsbytes)
	})
	return nil
}
