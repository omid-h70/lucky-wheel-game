package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omid-h70/lucky-wheel-game/db"
	"github.com/omid-h70/lucky-wheel-game/domain"
	"github.com/omid-h70/lucky-wheel-game/internal/response"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	hDB          db.DBHandler
	Addr         string
	stopSignal   chan os.Signal
	server       *http.Server
	gameLimit    time.Duration
	gameDuration time.Duration
	gameOverTime time.Duration
}

type UserRequest struct {
	UUID string `json:"UUID" validate:"required"`
}

type UserResponse struct {
	Prize      string `json:"prize"`
	Message    string `json:"message"`
	UUID       string `json:"uuid"`
	DailyCount int64  `json:"daily_count"`
}

var (
	ErrInvalidBodyRequestType = errors.New("Invalid Body Request Type")
	ErrInvalidApplicationType = errors.New("Request Application Type Must be json")
)

func NewApp(db db.DBHandler) *App {
	return &App{
		hDB:          db,
		stopSignal:   make(chan os.Signal, 1),
		gameDuration: 5 * time.Second,
	}
}

func (a *App) defaultHandler(w http.ResponseWriter, _ *http.Request) {
	response.NewError(errors.New("Invalid request"), http.StatusBadRequest).Send(w)
}

func (a *App) healthCheck(w http.ResponseWriter, _ *http.Request) {
	response.NewSuccess("Yo I'm up", http.StatusOK).Send(w)
}

func (a *App) testYourLuck(w http.ResponseWriter, r *http.Request) {

	var err error
	var req UserRequest
	var resp UserResponse

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NewError(ErrInvalidBodyRequestType, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" &&
		r.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		response.NewError(ErrInvalidApplicationType, http.StatusBadRequest).Send(w)
		return
	}

	time.Sleep(a.gameDuration)

	resp.DailyCount, err = a.hDB.InsertOrUpdateUUID(req.UUID, a.gameOverTime)
	fmt.Printf("Count %d \n", resp.DailyCount)

	if err != nil {
		response.NewError(err, 200).Send(w)
		return
	}

	log.Printf("uuid %s cnt %d \n"+req.UUID, string(resp.DailyCount))

	seeds, _ := a.hDB.GetDBSeeds()
	sortedSeeds := domain.SortMapBasedOnValue(seeds)
	prize, _ := domain.GetPrizeV2(domain.GetRandomNumberByTime(), sortedSeeds)

	resp.Prize = prize
	resp.UUID = req.UUID
	resp.Message = "you win"

	data, _ := json.Marshal(resp)
	response.NewSuccess(string(data), 200).Send(w)
}

func (a *App) SetAddr(addr string) *App {
	a.Addr = addr
	return a
}

func (a *App) SetDuration(time time.Duration) *App {
	a.gameLimit = time
	return a
}

func (a *App) SetCalcTime(time time.Duration) *App {
	a.gameDuration = time
	return a
}

func (a *App) SetGameOverTime(time time.Duration) *App {
	a.gameOverTime = time
	return a
}

func (a *App) SetAppHandlers(router *mux.Router) *App {

	api := router.PathPrefix("/v1").Subrouter()
	api.NotFoundHandler = http.HandlerFunc(a.defaultHandler)
	api.HandleFunc("/health", a.healthCheck)
	api.HandleFunc("/lottery", a.testYourLuck)
	return a
}

func (a *App) Listen(router *mux.Router) *App {

	a.server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         a.Addr,
		Handler:      router,
	}

	signal.Notify(a.stopSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Fatalln("Error starting HTTP server")
		}
	}()
	fmt.Printf("Server is Listening On %s", a.Addr)
	return a
}

func (a *App) WatchServer() *App {
	<-a.stopSignal

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Failed")
	}

	log.Fatal("Service is down")
	return a
}
