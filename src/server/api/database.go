package api

import (
	"encoding/json"
	"fmt"
	"modelhelper/cli/app"
	"net/http"

	"github.com/gorilla/mux"
)

var context *app.Context

func LoadDataSourceRoutes(router *mux.Router, ctx *app.Context) {
	context = ctx
	router.HandleFunc("/api/sources", sourcesHandler)
	router.HandleFunc("/api/sources/{source}/entities", entitiesHandler)
	router.HandleFunc("/api/sources/{source}/entities/{entity}", entityHandler)

}

func sourcesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	c := app.Configuration
	output, err := json.Marshal(c)

	if err != nil {
		panic("Could not marshal json.")
	}

	fmt.Fprintf(responseWriter, string(output))
}

func entitiesHandler(responseWriter http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	source := vars["source"]

	con := context.Connections[source]

	input := con.LoadSource()

	items, _ := input.Entities("")

	output, err := json.Marshal(items)

	if err != nil {
		panic("Could not marshal json.")
	}

	fmt.Fprintf(responseWriter, string(output))
}
func entityHandler(responseWriter http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	source := vars["source"]
	entity := vars["entity"]

	con := context.Connections[source]

	input := con.LoadSource()

	items, _ := input.Entity(entity)

	output, err := json.Marshal(items)

	if err != nil {
		panic("Could not marshal json.")
	}

	fmt.Fprintf(responseWriter, string(output))
}
