package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PalmerTurley34/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	numServerHits int
	database      *database.DB
	jwtSecret     string
	polkaAPIKey   string
}

func main() {
	ex, _ := os.Executable()
	currPath := filepath.Dir(ex)

	godotenv.Load(currPath + "/.env")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_API_KEY")

	dbFile := currPath + "/database.json"

	debugPtr := flag.Bool("debug", false, "")
	flag.Parse()
	if *debugPtr {
		os.Remove(dbFile)
	}

	mainRouter := chi.NewRouter()
	corsMux := middlewareCors(mainRouter)
	httpServer := &http.Server{
		Handler: corsMux,
		Addr:    ":8080",
	}
	db, err := database.NewDB(dbFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	apiCfg := apiConfig{
		numServerHits: 0,
		database:      db,
		jwtSecret:     jwtSecret,
		polkaAPIKey:   polkaKey,
	}

	fileHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(currPath))))
	mainRouter.Handle("/app/*", fileHandler)
	mainRouter.Handle("/app", fileHandler)

	apiRouter := getApiRouter(&apiCfg)
	adminRouter := getAdminRouter(&apiCfg)

	mainRouter.Mount("/api", apiRouter)
	mainRouter.Mount("/admin", adminRouter)

	fmt.Println("Start up server...")
	httpServer.ListenAndServe()
}
