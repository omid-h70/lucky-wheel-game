package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omid-h70/lucky-wheel-game/db"
	"log"
	"os"
	"time"
)

type AppConfig struct {
	serverConfig  string
	redisDbConfig string
	gameLimit     time.Duration
	gameDuration  time.Duration
	gameOverTime  time.Duration
}

var TestConfig AppConfig = AppConfig{
	serverConfig:  "0.0.0.0:8000",
	redisDbConfig: "0.0.0.0:6379",
	gameLimit:     3,
	gameDuration:  1 * time.Second,
	gameOverTime:  24 * time.Hour,
}

var DockerConfig AppConfig = AppConfig{
	serverConfig:  fmt.Sprintf("%s:%s", os.Getenv("APP_SERVER_ADDR"), os.Getenv("APP_HOST_PORT")),
	redisDbConfig: fmt.Sprintf("%s:%s", os.Getenv("REDIS_CONTAINER_NAME"), os.Getenv("REDIS_HOST_PORT")),
	gameLimit:     3,
	gameDuration:  1 * time.Second,
	gameOverTime:  24 * time.Hour,
}

var testPrizeList map[string]string = map[string]string{
	"A": "0.1",
	"B": "0.3",
	"C": "0.2",
	"D": "0.15",
	"E": "0.25",
}

var ConfigPrizeList map[string]string = map[string]string{
	"A": os.Getenv("AAA"),
	"B": os.Getenv("BBB"),
	"C": os.Getenv("CCC"),
	"D": os.Getenv("DDD"),
	"E": os.Getenv("EEE"),
}

func main() {
	//Test
	//appConfig := TestConfig
	//prizeList := testPrizeList

	//Operational
	appConfig := DockerConfig
	prizeList := ConfigPrizeList

	fmt.Println("Lucky Wheel Has Started")
	fmt.Println("Test Your Luck !")
	fmt.Println("Config \n", appConfig)

	redisClient := db.NewRedisHandler()
	err := redisClient.ConnectToDB(appConfig.redisDbConfig)
	if err != nil {
		log.Fatalln("No db Available")
	}

	fmt.Println("Prize List \n", prizeList)

	err = redisClient.SetDBSeeds(prizeList)
	if err != nil {
		log.Fatalln("DB Insertion Error")
	}

	router := mux.NewRouter()
	NewApp(redisClient).
		SetAddr(appConfig.serverConfig).
		SetDuration(appConfig.gameLimit).
		SetCalcTime(appConfig.gameDuration).
		SetGameOverTime(appConfig.gameOverTime).
		SetAppHandlers(router).
		Listen(router).
		WatchServer()
}
