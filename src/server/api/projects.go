package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type projectApi struct {
	Name string
	Path string
}

// LoadProjectRoutes loads all routes for
func LoadProjectRoutes(router *mux.Router) {
	router.HandleFunc("/api/projects", projectsHandler)

}

func projectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	projects := [2]projectApi{
		{Name: "PatoLab", Path: "<full-path>/.modelhelper"},
		{Name: "PatoLink", Path: "<full-path>/.modelhelper"},
	}

	output, err := json.Marshal(projects)

	if err != nil {
		panic("Could not marshal json.")
	}

	fmt.Fprintf(responseWriter, string(output))
}
