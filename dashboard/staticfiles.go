// +build !develop

package dashboard

import "github.com/gorilla/mux"

func registerStaticFileRoutes(r *mux.Router, prefix string) error {
	panic("not implemented") // TODO(c-bata): embed static files in Go
}
