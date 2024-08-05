package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var gcpAPIKey string
var ytConn *youtube.Service

var httpRespHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  "*", // update with your origins
	"Access-Control-Allow-Headers": "Content-Type",
}

func lambdaSetup() error {

	gcpAPIKey = os.Getenv("GCP_APIKEY")
	if gcpAPIKey == "" {
		return fmt.Errorf("GCP_APIKEY is required")
	}

	var err error
	ctx := context.Background()
	ytConn, err = youtube.NewService(ctx, option.WithAPIKey(gcpAPIKey))
	if err != nil {
		return err
	}
	return nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("DEBUG: entire req: %+v\n", request)

	err := lambdaSetup()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    httpRespHeaders,
		}, err
	}

	switch request.Resource {

	case "/v1/api/channel/{id}/insights":
		return getChannelInsights(request)

	case "/v1/api/channel/{id}/videos":
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotImplemented,
		}, fmt.Errorf("TBD")

	case "/v1/api/video/{id}/insights":
		return getVideoInsights(request)

	case "/v1/api/video/{id}/sentiments":
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotImplemented,
		}, fmt.Errorf("TBD")

	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    httpRespHeaders,
		}, fmt.Errorf("invalid request")
	}
}

// Define a struct for the JSON response
type Response struct {
	Status   string                     `json:"status"`
	Insights *youtube.ChannelStatistics `json:"insights"`
}

func getChannelInsights(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// channelID := request.QueryStringParameters["id"]
	channelID := request.PathParameters["id"]
	fmt.Printf("getChannelInsights : DEBUG: entire channelID: %+v\n", channelID)

	ChannelStatsCall := ytConn.Channels.List([]string{"statistics"}).Id(channelID)
	resp, err := ChannelStatsCall.Do()
	if err != nil {
		code := http.StatusInternalServerError
		if apiErr, ok := err.(*googleapi.Error); ok {
			code = apiErr.Code
		}
		return events.APIGatewayProxyResponse{
			StatusCode: code,
			Headers:    httpRespHeaders,
		}, err
	}

	if resp.Items == nil || resp.Items[0] == nil {
		err = fmt.Errorf("channel id: %v not found", channelID)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    httpRespHeaders,
		}, err
	}

	channel := resp.Items[0]
	Resp := Response{
		Status:   "success",
		Insights: channel.Statistics,
	}

	respBody, err := json.Marshal(Resp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    httpRespHeaders,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
		Headers:    httpRespHeaders,
	}, nil
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

// Define a struct for the JSON response
type ResponseVideoInsignts struct {
	Status string        `json:"status"`
	Vm     *VideoMetrics `json:"insights"`
}

func getVideoInsights(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// videoID := request.QueryStringParameters["id"]
	videoID := request.PathParameters["id"]
	fmt.Printf("getChannelInsights : DEBUG: entire channelID: %+v\n", videoID)

	videoListApiCall := ytConn.Videos.List([]string{"snippet,statistics"}).Id(videoID)
	resp, err := videoListApiCall.Do()
	if err != nil {
		code := http.StatusInternalServerError
		if apiErr, ok := err.(*googleapi.Error); ok {
			code = apiErr.Code
		}
		return events.APIGatewayProxyResponse{
			StatusCode: code,
			Headers:    httpRespHeaders,
		}, err
	}

	// Check if there are any videos matching the request
	if len(resp.Items) == 0 || resp.Items[0] == nil {
		err := fmt.Errorf("video id: %v not found", videoID)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    httpRespHeaders,
		}, err
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

	Resp := ResponseVideoInsignts{
		Status: "success",
		Vm:     &vm,
	}

	respBody, err := json.Marshal(Resp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    httpRespHeaders,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
		Headers:    httpRespHeaders,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
