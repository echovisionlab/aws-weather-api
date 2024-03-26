package app

import (
	"context"
	"errors"
	"github.com/echovisionlab/aws-weather-api/pkg/constants"
	"github.com/echovisionlab/aws-weather-api/pkg/database"
	"github.com/echovisionlab/aws-weather-api/pkg/handler"
	"github.com/echovisionlab/aws-weather-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"strconv"
)

type App struct {
	Service  *service.Service
	Validate *validator.Validate
}

func New() (*App, error) {
	log.SetOutput(os.Stdout)
	validate := validator.New()

	if svc, err := GetService(validate); err != nil {
		return nil, err
	} else {
		return &App{svc, validate}, nil
	}
}

func GetService(v *validator.Validate) (*service.Service, error) {
	db, err := database.New(database.WithEnvConfig())
	if err != nil {
		return nil, err
	}
	mps, err := maxPageSize()
	if err != nil {
		return nil, err
	}

	config := &service.Config{
		DB:          db,
		MaxPageSize: mps,
	}

	if err = v.Struct(config); err != nil {
		return nil, err
	}

	return service.New(config), nil
}

func maxPageSize() (int, error) {
	v, ok := os.LookupEnv(constants.MaxPageSize)
	if !ok {
		return 100, nil
	}
	return strconv.Atoi(v)
}

func (app *App) Run() func() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api/v1/record", handler.GetRecord(app.Service, app.Validate))
	router.GET("/api/v1/station", handler.GetStation(app.Service, app.Validate))

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
