package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rajeshamdev/analytics/yutube/gocode/channelInsights"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var sigChan chan os.Signal
var insightsHTTPServer *http.Server
var insightsGoRoutinesCnt int
var insightsGoRoutinesWG sync.WaitGroup
var apiKey string

func insightsInit() {

	// block all async signals to this server
	signal.Ignore()

	// create buffered signal channel and register below signals:
	//   - SIGINT (Ctrl+C)
	//   - SIGHUP (reload config)
	//   - SIGTERM (graceful shutdown)
	//   - SIGCHLD (handle child processes sending signal to parent) - ignore for now
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGCHLD)

	insightsRouter := gin.Default()
	insightsRouter.Use(gin.Logger())
	insightsRouter.GET("/v1/api/channel/:id/metrics", channelInsights.GetChannelMetrics)
	insightsRouter.GET("/v1/api/channel/:id/videos", channelInsights.GetChannelVideos)
	insightsRouter.GET("/v1/api/video/:id/metrics", channelInsights.GetVideoMetrics)
	insightsRouter.GET("/v1/api/video/:id/sentiments", channelInsights.VideoSentiment)

	insightsHTTPServer = &http.Server{
		Addr:    ":8080",
		Handler: insightsRouter,
	}
}

func insightsServerStart() {

	fmt.Printf("insightsServerStart starting\n")
	err := insightsHTTPServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("insightsServerStart: %v\n", err)
	}

	insightsGoRoutinesWG.Done()
}

func signalHandler() {

	for sig := range sigChan {

		if sig == syscall.SIGCHLD {
			continue
		} else if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			fmt.Printf("signal received: %v\n", sig)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := insightsHTTPServer.Shutdown(ctx)
			if err != nil {
				fmt.Printf("forced shutdown: %v", err)
			} else {
				fmt.Printf("graceful shutdown")
			}
		}
	}
}

// Config struct to match the JSON structure
type Config struct {
	APIKey string `json:"GCP_APIKEY"`
}

func readConfig() {
	// Open the configuration file
	file, err := os.Open("./.config")
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode JSON into Config struct
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("Error decoding config file: %v\n", err)
		os.Exit(1)
	}

	// Read environment variables from the Config struct
	apiKey = config.APIKey
	fmt.Printf("GCP_APIKEY: %v\n", apiKey)
}

func main() {

	readConfig()
	ctx := context.Background()

	var err error
	channelInsights.YoutubeService, err = youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Error creating new YouTube client: %v", err)
	}

	insightsInit()
	go signalHandler()

	insightsGoRoutinesCnt++

	go insightsServerStart()

	insightsGoRoutinesWG.Add(insightsGoRoutinesCnt)
	insightsGoRoutinesWG.Wait()

	//videoID := "vrOttI2cgAM"
	// channelID := "UCZN6X0ldwi-2W4TV-ab5M_g" // Thulasi
}
