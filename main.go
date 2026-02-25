package main

import "net/http"

func main() {

	serveMux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	serveMux.HandleFunc("/healthz", handleHealthz)
	server.ListenAndServe()

}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
