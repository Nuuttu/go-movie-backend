package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	mystructs "example.com/mysctructs"
)

type Movie = mystructs.Movie

var MovieList = mystructs.MovieList

func Movies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Movies sent")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MovieList) // Find a way to get movielist, or start using DB
}
