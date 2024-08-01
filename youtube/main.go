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

	"github.com/gin-contrib/cors"
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
	HTTPServerShutdownWaitSeconds int
)

var (
	GCPAPIKey    string
	allowOrigins string // // SHOULD be your frontend URL
)

const (
	// Refer: https://pkg.go.dev/github.com/gin-contrib/cors#Config
	allowOriginsDefault = "*"
)

func corsMiddleware() gin.HandlerFunc {
	// CORS middleware configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{allowOrigins}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}

	// Create middleware handler
	return cors.New(config)
}

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

	// Apply CORS middleware
	ginRouter.Use(corsMiddleware())

	ginRouter.GET("/v1/api/channel/:id/insights", utube.GetChannelInsights)
	ginRouter.GET("/v1/api/channel/:id/videos", utube.GetChannelVideos)
	ginRouter.GET("/v1/api/video/:id/insights", utube.GetVideoInsights)
	ginRouter.GET("/v1/api/video/:id/sentiments", utube.VideoSentiment)
	ginRouter.GET("/v1/api/health", utube.HealthCheck)

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

func readConfig() (*Config, error) {

	file, err := os.Open("./.config")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode JSON into Config struct
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func initEnv() {
	// pass GCP_APIKEY as env variable or in .config file in cwd
	GCPAPIKey = os.Getenv("GCP_APIKEY")
	if GCPAPIKey == "" {
		conf, err := readConfig()
		if err != nil || conf.APIKey == "" {
			fmt.Printf("pass GCP_APIKEY as env variable or set in .config\n")
			os.Exit(1)
		}
		GCPAPIKey = conf.APIKey
	}

	allowOrigins = os.Getenv("ALLOW_CORS_ORIGINS")
	if allowOrigins == "" {
		fmt.Printf("You must set ALLOW_CORS_ORIGINS in production environment\n")
		// this is not advisable in prod environment, but ok in dev/test
		allowOrigins = allowOriginsDefault
	}
}

func main() {

	initEnv()

	ctx := context.Background()
	var err error
	utube.YoutubeService, err = youtube.NewService(ctx, option.WithAPIKey(GCPAPIKey))
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
