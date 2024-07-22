package utube

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jonreiter/govader"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"
)

var YoutubeService *youtube.Service

// GET on /v1/api/channel/:id/insights : Fetch channel insights.
//
// It sends stats described at
// https://pkg.go.dev/google.golang.org/api@v0.188.0/youtube/v3#ChannelStatistics
//
// This method calls youtube.channels().list() API, which is subjective to API quota
// limitations.
// Refer: https://developers.google.com/youtube/v3/docs/channels/list
// Quota impact: A call to this method has a quota cost of 1 unit.

func GetChannelInsights(c *gin.Context) {

	channelID := c.Param("id")
	ChannelStatsCall := YoutubeService.Channels.List([]string{"statistics"}).Id(channelID)
	resp, err := ChannelStatsCall.Do()
	if err != nil {
		errMsg := fmt.Sprintf("%v", err)
		code := http.StatusInternalServerError
		if apiErr, ok := err.(*googleapi.Error); ok {
			code = apiErr.Code
		}
		c.JSON(code, gin.H{"message": errMsg})
		return
	}

	if resp.Items[0] == nil {
		errMsg := fmt.Sprintf("channel id: %v not found", channelID)
		c.JSON(http.StatusNotFound, gin.H{"message": errMsg})
		return
	}

	channel := resp.Items[0]
	c.JSON(http.StatusOK, channel.Statistics)
}

// GET on /v1/api/channel/:id/videos": Fetch all videos insights of a channel.
//
// youtube.search().list() API is called to fetch videos info of a channel, which is subjective
// to API quota limitations. TBD for optimizations to minimize calls.
// Refer: https://developers.google.com/youtube/v3/docs/search/list
// Quota impact: A call to this method has a quota cost of 100 units.
//
// Then, youtube.videos().list() API is called to fetch video stats, which has a quota
// cost of 1 unit. Refer: https://developers.google.com/youtube/v3/docs/videos/list

func GetChannelVideos(c *gin.Context) {

	channelID := c.Param("id")

	// this call has quota cost of 100 units.
	videosSearchApiCall := YoutubeService.Search.List([]string{"id,snippet"}).
		ChannelId(channelID).Type("video").MaxResults(50)
	resp, err := videosSearchApiCall.Do()
	if err != nil {
		errMsg := fmt.Sprintf("%v", err)
		code := http.StatusInternalServerError
		if apiErr, ok := err.(*googleapi.Error); ok {
			code = apiErr.Code
		}
		c.JSON(code, gin.H{"message": errMsg})
		return
	}

	// TBD (RAJ): repeat videosSearchApiCal.Do() to fetch all videos of channel.

	var vms []*VideoMetrics

	for _, item := range resp.Items {
		switch item.Id.Kind {
		// TBD (Raj): Is this necessary because we already called with "video" type
		case "youtube#video":
			// This method calls youtube.videos().list() API, which has a quota cost of 1 unit.
			vm, _, err := getVideoInsights((item.Id.VideoId))
			if err == nil {
				vms = append(vms, vm)
			} else {
				fmt.Printf("failed to get metrics for video: %v. error: %v\n", item.Id.VideoId, err)
			}
		}
	}

	c.JSON(http.StatusOK, vms)
}

// GET on /v1/api/video/:id/insights : Fetch insights of a video.
//
// youtube.videos().list() API is called to fetch video info, which is subjective to API
// quota limitations.
// Refer: https://developers.google.com/youtube/v3/docs/videos/list
// Quota impact: A call to this method has a quota cost of 1 unit.

func GetVideoInsights(c *gin.Context) {

	videoID := c.Param("id")
	vm, retCode, err := getVideoInsights((videoID))
	if err != nil {
		c.JSON(retCode, gin.H{"message": err})
		return
	}

	// TBD (Raj): work on local languages
	// c.Header("Content-Type", "application/json; charset=utf-8")
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

func getVideoInsights(videoID string) (*VideoMetrics, int, error) {

	videoListApiCall := YoutubeService.Videos.List([]string{"snippet,statistics"}).Id(videoID)
	resp, err := videoListApiCall.Do()
	if err != nil {
		code := http.StatusInternalServerError
		if apiErr, ok := err.(*googleapi.Error); ok {
			code = apiErr.Code
		}
		return nil, code, err
	}

	// Check if there are any videos matching the request
	if len(resp.Items) == 0 {
		errMsg := fmt.Sprintf("video id: %v not found", videoID)
		return nil, http.StatusNotFound, errors.New(errMsg)
	}

	video := resp.Items[0]
	var vm VideoMetrics

	vm.ChannelId = video.Snippet.ChannelId
	vm.ChannelTitle = video.Snippet.ChannelTitle
	vm.DefaultAudioLanguage = video.Snippet.DefaultAudioLanguage
	vm.DefaultLanguage = video.Snippet.DefaultLanguage
	vm.PublishedAt = video.Snippet.PublishedAt

	// TBD(Raj): utf-8
	// vm.Title = video.Snippet.Title
	fmt.Printf("Video Title: %v\n", video.Snippet.Title)
	// fmt.Printf("Description: %v\n", video.Snippet.Description)
	vm.Statistics = video.Statistics
	return &vm, http.StatusOK, nil
}

// Read comments of video using youtube.commentThreads().list() API multiple times.
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

type VideoSentiments struct {
	PositiveComments int `json:"positivecomments"`
	NegativeComments int `json:"negativecomments"`
}

// GET on /v1/api/video/:id/sentiments : Get sentiments on video
//
// Read comments of video and then find the sentiment score with
// VADER (Valence Aware Dictionary and sEntiment Reasoner).
// Refer: https://github.com/jonreiter/govader

// This method calls youtube.commentThreads().list() API multiple times,
// which is subjective to API quota limitations.
//
// Refer: https://developers.google.com/youtube/v3/docs/commentThreads/list
// Quota impact: each call has a quota cost of 1 unit.

func VideoSentiment(c *gin.Context) {

	videoID := c.Param("id")
	comments, _ := GetVideoComments(YoutubeService, videoID)

	positiveCommentsCnt := 0
	negativeCommentsCnt := 0
	for _, comment := range comments {

		//commenterID := comment.Snippet.TopLevelComment.Snippet.AuthorDisplayName
		analyzer := govader.NewSentimentIntensityAnalyzer()
		polarityScore := analyzer.PolarityScores(comment.Snippet.TopLevelComment.Snippet.TextDisplay)
		//sentStr := "negative"
		if polarityScore.Compound >= 0 {
			//sentStr = "positive"
			positiveCommentsCnt += 1
		} else {
			negativeCommentsCnt += 1
		}
		//fmt.Printf("%v: %v: %s\n", commenterID, sentStr, comment.Snippet.TopLevelComment.Snippet.TextDisplay)
	}
	fmt.Printf("positive comments: %v, negative comments: %v\n", positiveCommentsCnt, negativeCommentsCnt)

	sentiments := VideoSentiments{
		PositiveComments: positiveCommentsCnt,
		NegativeComments: negativeCommentsCnt,
	}

	c.JSON(http.StatusOK, sentiments)
}
