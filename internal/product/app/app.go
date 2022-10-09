package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/thangchung/go-coffeeshop/cmd/product/config"
	"github.com/thangchung/go-coffeeshop/internal/product/entity"
	log "github.com/thangchung/go-coffeeshop/pkg/logger"
)

var (
	router *mux.Router
	logger *log.Logger
)

func Run(cfg *config.Config) {
	var wait time.Duration

	logger = log.New(cfg.Log.Level)
	logger.Info("Init %s %s\n", cfg.Name, cfg.Version)

	// Repository
	// ...

	// Use case
	// ...

	// HTTP Server
	srv, err := initHTTPServer(cfg)
	if err != nil {
		logger.Fatal("%s", "cannot start server.")
	}

	logger.Info("%s %s.", "server start at", fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}

func getItemTypes(w http.ResponseWriter, r *http.Request) {
	logger.Info("%s", "GET: getItemTypes")

	itemTypeDtos := []entity.ItemTypeDto{
		{
			Name: "CAPPUCCINO",
			Type: 0,
		},
		{
			Name: "COFFEE_BLACK",
			Type: 1,
		},
		{
			Name: "COFFEE_WITH_ROOM",
			Type: 2,
		},
		{
			Name: "ESPRESSO",
			Type: 3,
		},
		{
			Name: "ESPRESSO_DOUBLE",
			Type: 4,
		},
		{
			Name: "LATTE",
			Type: 5,
		},
		{
			Name: "CAKEPOP",
			Type: 6,
		},
		{
			Name: "CROISSANT",
			Type: 7,
		},
		{
			Name: "MUFFIN",
			Type: 8,
		},
		{
			Name: "CROISSANT_CHOCOLATE",
			Type: 9,
		},
	}

	responseWithJSON(w, http.StatusOK, itemTypeDtos)
}

// func getItemByTypes(w http.ResponseWriter, r *http.Request) {
// 	logger.Info("%s", "GET: getItemByTypes")

// 	responseWithJson(w, http.StatusOK)
// }

func responseWithError(w http.ResponseWriter, code int, message string) {
	responseWithJSON(w, code, map[string]string{"error": message})
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		logger.Error("%s", "couldn't marshal object.")
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(code)
	w.Write(res)
}

func initHTTPServer(cfg *config.Config) (*http.Server, error) {
	router = mux.NewRouter()
	router.HandleFunc("/v1/api/item-types", getItemTypes).Methods("GET")
	// router.HandleFunc("/v1/api/item-by-types", getItemByTypes).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("%s", err)
		}
	}()

	return srv, nil
}