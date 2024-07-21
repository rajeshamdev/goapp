package main

import (
	"context"
	"log"

	"github.com/rajeshamdev/analytics/yutube/gocode/api"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {

	apiKey := ""

	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	videoID := "vrOttI2cgAM" // Replace with the ID of the video you want to analyze

	api.VideoSentiment(service, videoID)

	//channelID := "UCZN6X0ldwi-2W4TV-ab5M_g" // Thulasi
	// api.ListChannelVideos(service, channelID)
	//api.GetChannelMetrics(service, channelID)
}
