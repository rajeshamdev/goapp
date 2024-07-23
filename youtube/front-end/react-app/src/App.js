import React, { useState } from 'react';
import './App.css';

function App() {
  const [channelId, setChannelId] = useState('');
  const [videoId, setVideoId] = useState('');
  const [channelInsights, setChannelInsights] = useState(null);
  const [videoInsights, setVideoInsights] = useState(null);
  const [videoSentiments, setVideoSentiments] = useState(null);
  const [error, setError] = useState('');

  const fetchChannelInsights = async (id) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/api/channel/${encodeURIComponent(id)}/insights`);

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data = await response.json();
      setChannelInsights(data);
      setError('');
    } catch (error) {
      console.error('Error:', error);
      setError('Error fetching channel insights');
      setChannelInsights(null);
    }
  };

  const fetchVideoInsights = async (id) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/api/video/${encodeURIComponent(id)}/insights`);

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data = await response.json();
      setVideoInsights(data);
      setError('');
    } catch (error) {
      console.error('Error:', error);
      setError('Error fetching video insights');
      setVideoInsights(null);
    }
  };

  const fetchVideoSentiments = async (id) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/api/video/${encodeURIComponent(id)}/sentiments`);

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data = await response.json();
      setVideoSentiments(data);
      setError('');
    } catch (error) {
      console.error('Error:', error);
      setError('Error fetching video sentiments');
      setVideoSentiments(null);
    }
  };

  const handleChannelSubmit = async (event) => {
    event.preventDefault();
    await fetchChannelInsights(channelId);
  };

  const handleVideoSubmit = async (event) => {
    event.preventDefault();
    await fetchVideoInsights(videoId);
    await fetchVideoSentiments(videoId); // Fetch sentiments after fetching insights
  };

  const handleClear = () => {
    setChannelId('');
    setVideoId('');
    setChannelInsights(null);
    setVideoInsights(null);
    setVideoSentiments(null);
    setError('');
  };

  
  return (
    <div className="App">
      <h1>Fetch Channel and Video Insights</h1>

      <div className="fetch-form">
        <form onSubmit={handleChannelSubmit}>
          <label>
            Enter Channel ID:
            <input
              type="text"
              value={channelId}
              onChange={(e) => setChannelId(e.target.value)}
              required
            />
          </label>
          <button type="submit">Fetch Channel Insights</button>
        </form>

        <form onSubmit={handleVideoSubmit}>
          <label>
            Enter Video ID:
            <input
              type="text"
              value={videoId}
              onChange={(e) => setVideoId(e.target.value)}
              required
            />
          </label>
          <button type="submit">Fetch Video Insights</button>
        </form>
      </div>

      <button className="clear-button" onClick={handleClear}>Clear</button>

      {error && <p className="error">{error}</p>}

      {channelInsights && (
        <div className="insights">
          <h2>Channel Insights:</h2>
          <p>Subscriber Count: {channelInsights.subscriberCount}</p>
          <p>Video Count: {channelInsights.videoCount}</p>
          <p>View Count: {channelInsights.viewCount}</p>
        </div>
      )}

      {videoInsights && (
        <div className="insights">
          <h2>Video Insights:</h2>
          <p>Channel ID: {videoInsights.channelId}</p>
          <p>Channel Title: {videoInsights.channelTitle}</p>
          <p>Default Audio Language: {videoInsights.defaultAudioLanguage}</p>
          <p>Default Language: {videoInsights.defaultLanguage}</p>
          <p>Published At: {videoInsights.publishedAt}</p>
          <p>Comment Count: {videoInsights.statistics.commentCount}</p>
          <p>Like Count: {videoInsights.statistics.likeCount}</p>
          <p>View Count: {videoInsights.statistics.viewCount}</p>
        </div>
      )}

      {videoSentiments && (
        <div className="insights">
          <h2>Video Sentiments:</h2>
          <p>Positive Comments: {videoSentiments.positivecomments}</p>
          <p>Negative Comments: {videoSentiments.negativecomments}</p>
        </div>
      )}
    </div>
  );
}

export default App;
