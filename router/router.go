package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ishanjain28/pogo/auth"
	"github.com/ishanjain28/pogo/common"
)

type NewConfig struct {
	Name        string
	Host        string
	Email       string
	Description string
	Image       string
	PodcastURL  string
}

// Handle takes multiple Handler and executes them in a serial order starting from first to last.
// In case, Any middle ware returns an error, The error is logged to console and sent to the user, Middlewares further up in chain are not executed.
func Handle(handlers ...common.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rc := &common.RouterContext{}
		for _, handler := range handlers {
			err := handler(rc, w, r)
			if err != nil {
				log.Printf("%v", err)

				w.Write([]byte(http.StatusText(err.StatusCode)))

				return
			}
		}
	})
}

func Init() *mux.Router {

	r := mux.NewRouter()

	// "Static" paths
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/web/static"))))
	r.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("podcasts"))))

	// Paths that require specific handlers
	r.Handle("/", Handle(
		rootHandler(),
	)).Methods("GET")

	r.Handle("/rss", Handle(
		rootHandler(),
	)).Methods("GET")

	r.Handle("/json", Handle(
		rootHandler(),
	)).Methods("GET")

	// Authenticated endpoints should be passed to BasicAuth()
	// first
	r.Handle("/admin", Handle(
		auth.RequireAuthorization(),
		adminHandler(),
	))
	// r.HandleFunc("/admin/publish", BasicAuth(CreateEpisode))
	// r.HandleFunc("/admin/delete", BasicAuth(RemoveEpisode))
	// r.HandleFunc("/admin/css", BasicAuth(CustomCss))

	r.Handle("/setup", Handle(
		serveSetup(),
	)).Methods("GET", "POST")

	return r
}

// Handles /, /feed and /json endpoints
func rootHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		var file string
		switch r.URL.Path {
		case "/rss":
			w.Header().Set("Content-Type", "application/rss+xml")
			file = "assets/web/feed.rss"
		case "/json":
			w.Header().Set("Content-Type", "application/json")
			file = "assets/web/feed.json"
		case "/":
			w.Header().Set("Content-Type", "text/html")
			file = "assets/web/index.html"
		default:
			return &common.HTTPError{
				Message:    fmt.Sprintf("%s: Not Found", r.URL.Path),
				StatusCode: http.StatusNotFound,
			}
		}

		return common.ReadAndServeFile(file, w)
	}
}

func adminHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		return common.ReadAndServeFile("assets/web/admin.html", w)
	}
}

// Serve setup.html and config parameters
func serveSetup() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		if r.Method == "GET" {
			return common.ReadAndServeFile("assets/web/setup.html", w)
		}
		r.ParseMultipartForm(32 << 20)

		// Parse form and convert to JSON
		cnf := NewConfig{
			strings.Join(r.Form["podcastname"], ""),  // Podcast name
			strings.Join(r.Form["podcasthost"], ""),  // Podcast host
			strings.Join(r.Form["podcastemail"], ""), // Podcast host email
			"", // Podcast image
			"", // Podcast location
			"", // Podcast location
		}

		b, err := json.Marshal(cnf)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("assets/config/config.json", b, 0644)
		w.Write([]byte("Done"))
		return nil
	}
}

func redirectHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		return nil
	}
}
