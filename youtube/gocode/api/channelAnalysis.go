package api

import (
	"fmt"
	"log"

	"github.com/jonreiter/govader"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"
)

// This method calls youtube.search().list() API multiple times, which is
// subjective to API quota limitations. If possible, look for optimizations to
// minimize calls.
//
// Refer: https://developers.google.com/youtube/v3/docs/search/list
// Quota impact: A call to this method has a quota cost of 100 units.

func ListChannelVideos(service *youtube.Service, channelID string) {

	videosSearchApiCall := service.Search.List([]string{"id,snippet"}).ChannelId(channelID).Type("video").MaxResults(50)
	resp, err := videosSearchApiCall.Do()
	if err != nil {
		fmt.Printf("Error calling Search.List() API call: %v", err)
		return
	}

	for _, item := range resp.Items {

		switch item.Id.Kind {

		case "youtube#video":
			// This method calls youtube.videos().list() API, which has a quota cost of 1 unit.
			videoID := item.Id.VideoId
			GetVideoMetrics(service, videoID)
		}
	}
}

// This method calls youtube.videos().list() API, which is subjective to API
// quota limitations. If possible, look for optimizations to minimize this call.
//
// Refer: https://developers.google.com/youtube/v3/docs/videos/list
// Quota impact: A call to this method has a quota cost of 1 unit.

// Metrics collected include:
//  - num of likes
//  - num of views
//  - num of comments
//  - favorite count
//  - published date

func GetVideoMetrics(service *youtube.Service, videoID string) {

	videoListApiCall := service.Videos.List([]string{"snippet,statistics"}).Id(videoID)
	resp, err := videoListApiCall.Do()
	if err != nil {
		fmt.Printf("Error calling Videos.List() API call: %v", err)
		return
	}

	// Check if there are any videos matching the request
	if len(resp.Items) == 0 {
		fmt.Println("No videos found.")
		return
	}

	video := resp.Items[0]
	fmt.Printf("Video Title: %v\n", video.Snippet.Title)
	fmt.Printf("Video ID: %v\n", video.Id)
	fmt.Printf("View Count: %v\n", video.Statistics.ViewCount)
	fmt.Printf("Like Count: %v\n", video.Statistics.LikeCount)
	fmt.Printf("Comment Count: %v\n", video.Statistics.CommentCount)
	fmt.Printf("Dislike Count: %v\n", video.Statistics.DislikeCount)
	fmt.Printf("Published At: %v\n", video.Snippet.PublishedAt)
	// fmt.Printf("Description: %v\n", video.Snippet.Description)
}

func GetChannelMetrics(service *youtube.Service, channelID string) (*youtube.ChannelStatistics, error) {

	// Define the parameters for the channel list request
	call := service.Channels.List([]string{"snippet,statistics"}).Id(channelID)
	resp, err := call.Do()
	if err != nil {
		log.Fatal("Error making channel list API call: %v", err)
	}

	// Print the title and ID of the channel
	channel := resp.Items[0]
	fmt.Printf("Channel Title: %v\n", channel.Snippet.Title)
	fmt.Printf("Channel subscribers count: %v\n", channel.Statistics.SubscriberCount)
	fmt.Printf("Channel Video count: %v\n", channel.Statistics.VideoCount)
	fmt.Printf("Channel total views count: %v\n", channel.Statistics.ViewCount)

	return channel.Statistics, nil
}

func GetVideoComments(service *youtube.Service, videoID string) ([]*youtube.CommentThread, error) {

	call := service.CommentThreads.List([]string{"snippet"}).VideoId(videoID).MaxResults(100) // Adjust as needed, maximum 100 per request

	var comments []*youtube.CommentThread
	nextPageToken := ""

	for {
		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		response, err := call.Do()
		if err != nil {
			if apiErr, ok := err.(*googleapi.Error); ok {
				if apiErr.Code == 403 {
					log.Println("API rate limit exceeded. Waiting for quota to be replenished...")
					// Handle rate limit gracefully if needed
					return comments, err
				}
			}
			return comments, err
		}

		comments = append(comments, response.Items...)

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return comments, nil
}

func VideoSentiment(service *youtube.Service, videoID string) {

	comments, _ := GetVideoComments(service, videoID)

	for _, comment := range comments {
		commenter := comment.Snippet.TopLevelComment.Snippet.AuthorDisplayName
		analyzer := govader.NewSentimentIntensityAnalyzer()
		polarityScore := analyzer.PolarityScores(comment.Snippet.TopLevelComment.Snippet.TextDisplay)
		sentStr := "negative"
		if polarityScore.Compound >= 0 {
			sentStr = "positive"
		}
		fmt.Printf("%v: %v: %s\n", commenter, sentStr, comment.Snippet.TopLevelComment.Snippet.TextDisplay)
	}
}
