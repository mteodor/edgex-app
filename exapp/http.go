package exapp

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func AppStatus(w http.ResponseWriter, r *http.Request) {
	// Check if app is working
	fmt.Fprintf(w, "%q", html.EscapeString("edgex-app iw working\n"))
}

func InitHTTP(port string) error {
	fmt.Println("initializing server")
	r := mux.NewRouter()
	r.HandleFunc("/", AppStatus)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()

}
