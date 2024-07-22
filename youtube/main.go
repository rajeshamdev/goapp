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
	"github.com/rajeshamdev/analytics/youtube/utube"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	sigChan                       chan os.Signal
	HTTPServer                    *http.Server
	goRoutinesCnt                 int
	goRoutinesWG                  sync.WaitGroup
	apiKey                        string
	HTTPServerShutdownWaitSeconds int
)

func serverInit() {

	// block all async signals to this server
	signal.Ignore()

	// create buffered signal channel and register below signals:
	//   - SIGINT (Ctrl+C)
	//   - SIGHUP (reload config)
	//   - SIGTERM (graceful shutdown)
	//   - SIGCHLD (handle child processes sending signal to parent) - ignore for now
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGCHLD)

	ginRouter := gin.Default()
	ginRouter.Use(gin.Logger())
	ginRouter.GET("/v1/api/channel/:id/insights", utube.GetChannelInsights)
	ginRouter.GET("/v1/api/channel/:id/videos", utube.GetChannelVideos)
	ginRouter.GET("/v1/api/video/:id/insights", utube.GetVideoInsights)
	ginRouter.GET("/v1/api/video/:id/sentiments", utube.VideoSentiment)

	HTTPServer = &http.Server{
		Addr:    ":8080",
		Handler: ginRouter,
	}

	HTTPServerShutdownWaitSeconds = 5
}

func serverStart() {

	fmt.Printf("serverStart starting\n")
	err := HTTPServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("serverStart error: %v. Exiting serverStart\n", err)
		os.Exit(1)
	}

	goRoutinesWG.Done()
}

func signalHandler() {

	for sig := range sigChan {

		if sig == syscall.SIGCHLD {
			continue
		} else if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			fmt.Printf("signal received: %v\n", sig)

			ctx, cancel := context.WithTimeout(context.Background(),
				time.Duration(HTTPServerShutdownWaitSeconds)*time.Second)

			err := HTTPServer.Shutdown(ctx)
			if err != nil {
				fmt.Printf("forced shutdown: %v", err)
			} else {
				fmt.Printf("graceful shutdown")
			}
			cancel()
		}
	}
}

type Config struct {
	APIKey string `json:"GCP_APIKEY"`
}

func readConfig() {

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

	apiKey = config.APIKey
}

func main() {

	readConfig()

	ctx := context.Background()
	var err error
	utube.YoutubeService, err = youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Error creating new YouTube client: %v. Exiting main", err)
		os.Exit(1)
	}

	serverInit()
	go signalHandler()

	goRoutinesCnt++

	go serverStart()

	goRoutinesWG.Add(goRoutinesCnt)
	goRoutinesWG.Wait()
	close(sigChan)
}
