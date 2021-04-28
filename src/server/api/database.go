package api

import (
	"encoding/json"
	"fmt"
	"modelhelper/cli/app"
	"modelhelper/cli/config"
	"modelhelper/cli/input"
	"net/http"

	"github.com/gorilla/mux"
)

type routeLoader struct {
	config *config.Config
}

func LoadDataSourceRoutes(router *mux.Router) {
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

	input := input.GetSource(source, *app.Configuration)

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

	input := input.GetSource(source, *app.Configuration)

	items, _ := input.Entity(entity)

	output, err := json.Marshal(items)

	if err != nil {
		panic("Could not marshal json.")
	}

	fmt.Fprintf(responseWriter, string(output))
}

func getSourceName(config *config.Config) string {
	defaultSource := config.DefaultSource

	if len(defaultSource) == 0 {
		if len(config.Sources) == 0 {
			defaultSource = ""
		} else {
			for _, s := range config.Sources {

				defaultSource = s.Name
				break
			}
		}

	}

	return defaultSource
}
