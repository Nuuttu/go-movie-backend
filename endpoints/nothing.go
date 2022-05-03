package endpoints

import (
	"net/http"

	"github.com/darahayes/go-boom"
)

func Nothing(w http.ResponseWriter, r *http.Request) {
	boom.NotFound(w, "Sorry, there's nothing here.")
	// https://pkg.go.dev/github.com/darahayes/go-boom?utm_source=godoc
}
