package main

import (
	"chirpy/internal/database"
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("JWT_SECRET")
	polka_api_key := os.Getenv("POLKA_API_KEY")
	if secret == "" {
		panic("JWT_SECRET is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)

	serveMux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	apiConfig := &apiConfig{}
	apiConfig.queries = dbQueries
	apiConfig.jwtSecret = secret
	apiConfig.polka_api_key = polka_api_key
	serveMux.Handle("/app/", http.StripPrefix("/app/", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /api/healthz", handleHealthz)
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiConfig.handleReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	serveMux.HandleFunc("POST /api/users", apiConfig.handleCreateUser)
	serveMux.HandleFunc("POST /api/chirps", apiConfig.handleCreateChirp)
	serveMux.HandleFunc("GET /api/chirps", apiConfig.getChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{id}", apiConfig.getChirpHandler)
	serveMux.HandleFunc("POST /api/login", apiConfig.handleLogin)
	serveMux.HandleFunc("POST /api/refresh", apiConfig.handleRefreshToken)
	serveMux.HandleFunc("POST /api/revoke", apiConfig.handleRevokeToken)
	serveMux.HandleFunc("PUT /api/users", apiConfig.handleUpdateUser)
	serveMux.HandleFunc("DELETE /api/chirps/{id}", apiConfig.handleDeleteChirp)
	serveMux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlePolkaWebhook)
	server.ListenAndServe()

}

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	jwtSecret      string
	polka_api_key  string
}
