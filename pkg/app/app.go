package app

import (
	"context"
	"errors"
	"github.com/echovisionlab/aws-weather-api/pkg/database"
	"github.com/echovisionlab/aws-weather-api/pkg/record"
	"github.com/echovisionlab/aws-weather-api/pkg/service"
	"github.com/echovisionlab/aws-weather-api/pkg/station"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
)

type App struct {
	Service  *service.Service
	Validate *validator.Validate
}

func New() (*App, error) {
	log.SetOutput(os.Stdout)
	db, err := database.New(database.WithEnvConfig())
	if err != nil {
		return nil, err
	}
	return &App{service.New(db), validator.New()}, nil
}

func (app *App) Run() func() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api/v1/record", record.Get(app.Service, app.Validate))
	router.GET("/api/v1/station", station.Get(app.Service, app.Validate))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	stop := func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error: %s\n", err)
		}
	}

	return stop
}
