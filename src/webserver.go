/* webserver.go
 * 
 * This is the webserver handler for White Rabbit, and handles
 * all incoming connections, including authentication. 
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/subtle"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// Prints out contects of feed.rss
func RssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml")
	w.WriteHeader(http.StatusOK)
	data, err := ioutil.ReadFile("feed.rss")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, string(data))
}

// Does the same as above method, only with the JSON feed data
func JsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := ioutil.ReadFile("feed.json")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, string(data))
}

// Serve up homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("assets/index.html")

	if err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
	}
}

// Authenticate user using basic webserver authentication
/*
 * Code from stackoverflow by user Timmmm
 * https://stackoverflow.com/questions/21936332/idiomatic-way-of-requiring-http-basic-auth-in-go/39591234#39591234
*/
func BasicAuth(handler http.HandlerFunc,) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
 	username := viper.GetString("AdminUsername")
 	password := viper.GetString("AdminPassword")
 	realm := "Login to White Rabbit admin interface"
        user, pass, ok := r.BasicAuth()

        if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
            w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
            w.WriteHeader(401)
            w.Write([]byte("Unauthorised.\n"))
            return
        }

        handler(w, r)
    }
}

// Handler for serving up admin page
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("assets/admin.html")

	if err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
	}
}

// Main function that defines routes
func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Start the watch() function in generate_rss.go, which
	// watches for file changes and regenerates the feed 
	go watch()

	// Define routes
	r := mux.NewRouter()

	// "Static" paths
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/static"))))
	r.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("podcasts"))))

	// Paths that require specific handlers
	http.Handle("/", r) // Pass everything to gorilla mux
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/rss", RssHandler)
	r.HandleFunc("/json", JsonHandler)

	// Authenticated endpoints should be passed to BasicAuth()
	// first
	r.HandleFunc("/admin", BasicAuth(AdminHandler))
	r.HandleFunc("/admin/publish", BasicAuth(CreateEpisode))
	r.HandleFunc("/admin/delete", BasicAuth(RemoveEpisode))
	r.HandleFunc("/admin/css", BasicAuth(CustomCss))

	// We're live!
	log.Fatal("Live at localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
