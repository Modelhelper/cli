package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"modelhelper/cli/app"
	"modelhelper/cli/server/api"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/browser"

	"github.com/gorilla/mux"
)

//go:embed static/*
var src embed.FS

type sourceApi struct {
	Name       string
	Connection string
}

type configApi struct {
	Version string
}

var sources [2]sourceApi = [2]sourceApi{
	{Name: "lab", Connection: "source info for lab"},
	{Name: "raw", Connection: "source info for raw"},
}

func Serve(port int, open bool, appCtx *app.Context) {
	wait := time.Second * 15

	r := mux.NewRouter()

	host := fmt.Sprintf("0.0.0.0:%v", port)
	// serve api
	// r.HandleFunc("/api/projects", projectsHandler)             // v3.1
	// r.HandleFunc("/api/sources/{name}/tables", sourcesHandler) // v3.0

	getRouter := r.Methods(http.MethodGet).Subrouter()

	api.LoadProjectRoutes(getRouter)
	api.LoadDataSourceRoutes(getRouter, appCtx)
	mime.AddExtensionType(".mjs", "text/javascript")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".svg", "image/svg+xml")
	// getRouter.HandleFunc("/api/projects", apiprojectsHandler)
	// serve website
	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      r,
		Addr:         host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//srv.ListenAndServe()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}

		if open {
			browser.OpenURL(host)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func getFileSystem(useOS bool, root string, f embed.FS) http.FileSystem {
	if useOS {
		return http.FS(os.DirFS(root))
	}

	fsys, err := fs.Sub(f, root)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// func projectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
// 	output, err := json.Marshal(projects)

// 	if err != nil {
// 		panic("Could not marshal json.")
// 	}

// 	fmt.Fprintf(responseWriter, string(output))
// }
func sourcesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	output, err := json.Marshal(sources)

	if err != nil {
		panic("Could not marshal json.")
	}

	fmt.Fprintf(responseWriter, string(output))
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Content-Type", "text/css; charset=utf-8")
	useOS := len(os.Args) > 1 && os.Args[1] == "live"

	http.FileServer(getFileSystem(useOS, "static", src)).ServeHTTP(w, r)
}
