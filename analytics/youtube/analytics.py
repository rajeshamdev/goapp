import os
import json
import pprint as pretty
import googleapiclient.discovery

## TBD(Raj) : Fix this
#  import ssl
#  ssl._create_default_https_context = ssl._create_unverified_context

import nltk
from nltk.sentiment.vader import SentimentIntensityAnalyzer
# TBD (Raj): Download only once
#nltk.download('vader_lexicon')

class YoutubeChannel:

    def __init__(self, id):
        self.channelID = id
        self.APIKEY = <>
        self.youtubeResource =  googleapiclient.discovery.build('youtube',
            'v3', developerKey=self.APIKEY)
        self.sid = SentimentIntensityAnalyzer()

    def formatNumber(self, num: int) -> str:
        magnitude = 0
        while abs(num) >= 1000:
            magnitude += 1
            num /= 1000.0
    
        if magnitude == 0:
            return f"{num:.0f}"
        else:
            return f"{num:.2f}{['', 'K', 'M', 'B', 'T'][magnitude]}"

    def getVideoMetrics(self, videoID: str) -> dict:
        """
        Returns video metrics

        This method calls youtube.videos().list() API, which is subjective to API
        quota limitations. If possible, look for optimizations to minimize this call.

        Refer: https://developers.google.com/youtube/v3/docs/videos/list
        Quota impact: A call to this method has a quota cost of 1 unit.

        Args:
            videoID (str): video id.

        Returns:
            dict: returns dict with following keys :
                  - 'commentCount'
                  - 'favoriteCount'
                  - 'likeCount'
                  - 'viewCount'
        """

        resp = self.youtubeResource.videos().list(part='statistics', id=videoID
        ).execute()
        return resp['items'][0]['statistics']

    def channelSummary(self) -> None:
        """
        summarize channel metrics:
         - total subscriber count
         - total videos uploaded
         - combined views of all videos uploaded in the channel

        This method calls youtube.channels().list() API, which is subjective to API quota
        limitations. If possible, look for optimizations to minimize this call.
        
        Refer: https://developers.google.com/youtube/v3/docs/channels/list
        Quota impact: A call to this method has a quota cost of 1 unit.

        Args:
            None

        Returns:
            None 
        """

        resp = self.youtubeResource.channels().list(part='statistics', id=self.channelID).execute()
        stats = resp['items'][0]['statistics']
        channelMetrics = {
            'subscriberCount'      : self.formatNumber(int(stats['subscriberCount'])),
            'videoCount'           : int(stats['videoCount']),
            'viewCount'            : self.formatNumber(int(stats['viewCount'])),
        }
        pretty.pprint(channelMetrics)

    def videoSummary(self, maxVideos: int) -> None:
        """
        summarize videos metrics:
         - num of likes
         - num of views
         - num of comments
         - favorite count
         - published date

        This method may call youtube.search().list() API multiple times, which is 
        subjective to API quota limitations. If possible, look for optimizations to
        minimize calls.

        Refer: https://developers.google.com/youtube/v3/docs/search/list
        Quota impact: A call to this method has a quota cost of 100 units.

        Args:
            maxVideos (int): num of videos

        Returns:
            None
        """
        videos = []
        nextPageToken = None
        videoCnt = 0

        # TBD (Raj): get exact video count requested by user
        while videoCnt <= maxVideos:
            resp = self.youtubeResource.search().list(
                    part       = 'snippet',
                    channelId  = self.channelID,
                    type       = 'video',
                    order      = 'viewCount',  # Sort by view count
                    maxResults = 50, # set to max 50 to minimize API calls
                    pageToken=nextPageToken
            ).execute()

            for result in resp.get('items', []):
                videoID = result['id']['videoId']
                videoTitle = result['snippet']['title']
                publishTime = result['snippet']['publishTime']
                channelTitle = result['snippet']['channelTitle']

                videoMetrics = self.getVideoMetrics(videoID)

                videoInfo = {
                    'videoID'            : videoID,
                    'videoTitle'         : videoTitle,
                    #'channelTitle'       : channelTitle,
                    'published datetime' : publishTime,
                    'commentCount'       : self.formatNumber(int(videoMetrics['commentCount'])),
                    'favoriteCount'      : self.formatNumber(int(videoMetrics['favoriteCount'])),
                    'likeCount'          : self.formatNumber(int(videoMetrics['likeCount'])),
                    'viewCount'          : self.formatNumber(int(videoMetrics['viewCount'])), 
                }
                videos.append(videoInfo)
                videoCnt += 1

            nextPageToken = resp.get('nextPageToken')
            if not nextPageToken:
                break

            print("total videos: ", len(videos), videoCnt)
        pretty.pprint(videos)

    def getVideoComments(self, videoID: str, maxComments: int = 500) -> list:
        """
        get vidoe comments.

        This method may call youtube.commentThreads().list() API multiple times,
        which is subjective to API quota limitations. If possible, look for
        optimizations to minimize calls.

        Refer: https://developers.google.com/youtube/v3/docs/commentThreads/list
        Quota impact: A call to this method has a quota cost of 1 unit.

        Args:
            videoID (str): video id
            maxComments (int): number of comments to retrieve

        Returns:
            list: comments list
        """
        comments = []

        request = self.youtubeResource.commentThreads().list(
            part='snippet',
            maxResults=maxComments,
            videoId=videoID,
            textFormat='plainText'
        )

        while request:
            response = request.execute()
            for item in response['items']:
                commenter = item['snippet']['topLevelComment']['snippet']['authorDisplayName']
                comment = item['snippet']['topLevelComment']['snippet']['textDisplay']
                comments.append(commenter + ': ' + comment)
            request = self.youtubeResource.commentThreads().list_next(request, response)
        return comments

    def analyzeSentiment(self, videoID: str) -> list:
        """
        Analyse comments sentiment.

        Args:
            videoID (str): video id
            maxComments (int): number of comments to retrieve

        Returns:
            list: list of comments with sentiments
        """
        comments = self.getVideoComments(videoID)
        results = []
        for comment in comments:
            sentiment_scores = self.sid.polarity_scores(comment)
            sentiment = 'positive' if sentiment_scores['compound'] >= 0 else 'negative'
            results.append((comment, sentiment))
        return results


# Example usage
if __name__ == '__main__':

    channelID  = 'UCz8QaiQxApLq8sLNcszYyJw'  # firstpost
    channelID  = 'UCdc6ObxhdQ8eZIFquU2xolA'  # Sahil Bhadviya
    channelID  = 'UCqW8jxh4tH1Z1sWPbkGWL4g'  # Akshat Shrivastava
    channelID  = 'UCZN6X0ldwi-2W4TV-ab5M_g'  # Thulasi Chandu
    yt = YoutubeChannel(channelID)
    yt.channelSummary()
    yt.videoSummary(51)

    videoID = 'vrOttI2cgAM'
    sentimentResults = yt.analyzeSentiment(videoID)
    for idx, (comment, sentiment) in enumerate(sentimentResults, start=1):
        print(f'Comment {idx}: {comment}')
        print(f'Sentiment: {sentiment}')
        print('---')
    
