package channelInsights

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jonreiter/govader"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"
)

var YoutubeService *youtube.Service

// To list channel videos, we utilize youtube.search().list() API, which is subjective
// to API quota limitations. If possible, look for optimizations to minimize calls.
// Refer: https://developers.google.com/youtube/v3/docs/search/list
// Quota impact: A call to this method has a quota cost of 100 units.
//
// Then, for each video metrics, we call youtube.videos().list() API, which has a quota
// cost of 1 unit. Refer: https://developers.google.com/youtube/v3/docs/videos/list

func GetChannelVideos(c *gin.Context) {

	channelID := c.Param("id")

	// this call has quota cost of 100 units.
	videosSearchApiCall := YoutubeService.Search.List([]string{"id,snippet"}).ChannelId(channelID).Type("video").MaxResults(50)
	resp, err := videosSearchApiCall.Do()
	if err != nil {
		fmt.Printf("Error calling Search.List() API call: %v", err)
		return
	}

	for _, item := range resp.Items {

		switch item.Id.Kind {

		case "youtube#video":

			videoID := item.Id.VideoId
			// This method calls youtube.videos().list() API, which has a quota cost of 1 unit.
			getVideoMetrics(videoID)
		}
	}
}

// This method calls youtube.videos().list() API, which is subjective to API
// quota limitations.
// Refer: https://developers.google.com/youtube/v3/docs/videos/list
// Quota impact: A call to this method has a quota cost of 1 unit.

// Metrics collected include:
//  - num of likes
//  - num of views
//  - num of comments
//  - favorite count
//  - published date

func GetVideoMetrics(c *gin.Context) {

	videoID := c.Param("id")
	vm := getVideoMetrics((videoID))
	//c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, vm)

}

type VideoMetrics struct {
	ChannelId            string                   `json:"channelId,omitempty"`
	ChannelTitle         string                   `json:"channelTitle,omitempty"`
	DefaultAudioLanguage string                   `json:"defaultAudioLanguage,omitempty"`
	DefaultLanguage      string                   `json:"defaultLanguage,omitempty"`
	PublishedAt          string                   `json:"publishedAt,omitempty"`
	Title                string                   `json:"title,omitempty"`
	Statistics           *youtube.VideoStatistics `json:"statistics,omitempty"`
}

func getVideoMetrics(videoID string) *VideoMetrics {

	videoListApiCall := YoutubeService.Videos.List([]string{"snippet,statistics"}).Id(videoID)
	resp, err := videoListApiCall.Do()
	if err != nil {
		fmt.Printf("Error calling Videos.List() API call: %v", err)
		// return json
	}

	// Check if there are any videos matching the request
	if len(resp.Items) == 0 {
		fmt.Println("No videos found.")
		// return json
	}

	video := resp.Items[0]
	var vm VideoMetrics
	/*

		DefaultLanguage string `json:"defaultLanguage,omitempty"`
		// Localized: Localized title and description, read-only.
		Localized   *ChannelLocalization     `json:"localized,omitempty"`
		PublishedAt string                   `json:"publishedAt,omitempty"`
		Title       string
	*/

	vm.ChannelId = video.Snippet.ChannelId
	vm.ChannelTitle = video.Snippet.ChannelTitle
	vm.DefaultAudioLanguage = video.Snippet.DefaultAudioLanguage
	vm.DefaultLanguage = video.Snippet.DefaultLanguage
	vm.PublishedAt = video.Snippet.PublishedAt
	// TBD(Raj): utf-8
	// vm.Title = video.Snippet.Title
	vm.Statistics = video.Statistics

	fmt.Printf("Video Title: %v\n", video.Snippet.Title)
	fmt.Printf("Video ID: %v\n", video.Id)
	fmt.Printf("View Count: %v\n", video.Statistics.ViewCount)
	fmt.Printf("Like Count: %v\n", video.Statistics.LikeCount)
	fmt.Printf("Comment Count: %v\n", video.Statistics.CommentCount)
	fmt.Printf("Dislike Count: %v\n", video.Statistics.DislikeCount)
	fmt.Printf("Published on: %v\n", video.Snippet.PublishedAt)
	return &vm
	// fmt.Printf("Description: %v\n", video.Snippet.Description)
	// return json
}

// This method calls youtube.channels().list() API, which is subjective to API quota
// limitations.
// Refer: https://developers.google.com/youtube/v3/docs/channels/list
// Quota impact: A call to this method has a quota cost of 1 unit.

func GetChannelMetrics(c *gin.Context) {

	// *youtube.ChannelStatistics
	channelID := c.Param("id")

	ChannelStatsCall := YoutubeService.Channels.List([]string{"statistics"}).Id(channelID)
	resp, err := ChannelStatsCall.Do()
	if err != nil {
		fmt.Printf("Error making channel list API call: %v", err)
		// return json
	}

	if resp.Items[0] == nil {
		fmt.Printf("channel id: %v not found", channelID)
		// return nil, fmt.Errorf("channel id: %v not found", channelID)
		// // return json
	}

	channel := resp.Items[0]
	c.JSON(http.StatusOK, channel.Statistics)
}

// This method calls youtube.commentThreads().list() API multiple times,
// which is subjective to API quota limitations. If possible, look for
// optimizations to minimize calls.
//
// Refer: https://developers.google.com/youtube/v3/docs/commentThreads/list
// Quota impact: each call has a quota cost of 1 unit.

func GetVideoComments(service *youtube.Service, videoID string) ([]*youtube.CommentThread, error) {

	var comments []*youtube.CommentThread
	nextPageToken := ""
	commentsCall := service.CommentThreads.List([]string{"snippet"}).VideoId(videoID).MaxResults(100)

	for {
		if nextPageToken != "" {
			commentsCall = commentsCall.PageToken(nextPageToken)
		}

		resp, err := commentsCall.Do()
		if err != nil {

			if apiErr, ok := err.(*googleapi.Error); ok {

				if apiErr.Code == 403 {

					fmt.Printf("API rate limit exceeded. Wait for quota to be replenished...")
					// Handle rate limit gracefully if needed
					return comments, err
				}
			}
			return comments, err
		}

		comments = append(comments, resp.Items...)

		nextPageToken = resp.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return comments, nil
}

// This method calls youtube.commentThreads().list() API multiple times,
// which is subjective to API quota limitations.
//
// Refer: https://developers.google.com/youtube/v3/docs/commentThreads/list
// Quota impact: each call has a quota cost of 1 unit.

func VideoSentiment(c *gin.Context) {

	videoID := c.Param("id")
	comments, _ := GetVideoComments(YoutubeService, videoID)

	for _, comment := range comments {

		commenterID := comment.Snippet.TopLevelComment.Snippet.AuthorDisplayName
		analyzer := govader.NewSentimentIntensityAnalyzer()
		polarityScore := analyzer.PolarityScores(comment.Snippet.TopLevelComment.Snippet.TextDisplay)
		sentStr := "negative"
		if polarityScore.Compound >= 0 {
			sentStr = "positive"
		}
		fmt.Printf("%v: %v: %s\n", commenterID, sentStr, comment.Snippet.TopLevelComment.Snippet.TextDisplay)
	}
}
