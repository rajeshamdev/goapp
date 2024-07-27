import React, { useState } from 'react';
import './App.css';
import config from './config'

import dega from './assets/images/dega.png'

const backendURL = config.backendURL;

function App() {
  const [channelId, setChannelId] = useState('');
  const [videoId, setVideoId] = useState('');
  const [channelInsights, setChannelInsights] = useState(null);
  const [videoInsights, setVideoInsights] = useState(null);
  const [videoSentiments, setVideoSentiments] = useState(null);
  const [error, setError] = useState('');

  const fetchChannelInsights = async (id) => {
    try {
      const response = await fetch(`${backendURL}/v1/api/channel/${encodeURIComponent(id)}/insights`);

      if (!response.ok) {
        const errorMessage = await response.json()
        throw new Error(errorMessage.message)
      }

      const data = await response.json();
      setChannelInsights(data);
      setError('');
    } catch (error) {
      console.error('Error:', error);
      setError(`${error.message}`);
      setChannelInsights(null);
    }
  };

  const fetchVideoInsights = async (id) => {
    try {
      const response = await fetch(`${backendURL}/v1/api/video/${encodeURIComponent(id)}/insights`);

      if (!response.ok) {
        const errorMessage = await response.json()
        throw new Error(errorMessage.message)
      }

      const data = await response.json();
      setVideoInsights(data);
      setError('');
    } catch (error) {
      console.error('Error:', error);
      setError(`${error.message}`);
      setVideoInsights(null);
    }
  };

  const fetchVideoSentiments = async (id) => {
    try {
      const response = await fetch(`${backendURL}/v1/api/video/${encodeURIComponent(id)}/sentiments`);

      if (!response.ok) {
        const errorMessage = await response.json()
        throw new Error(errorMessage.message)
      }

      const data = await response.json();
      setVideoSentiments(data);
      setError('');
    } catch (error) {
      console.error('Error:', error);
      setError(`${error.message}`);
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
      <h1>
        <div className="text-overlay">Unlock the Power of Data: Revealing Actionable Insights </div>
      </h1>
      <div className="fetch-form">
        <form onSubmit={handleChannelSubmit}>
          <label>
            Enter YouTube Channel ID:
            <input
              type="text"
              value={channelId}
              onChange={(e) => setChannelId(e.target.value)}
              required
            />
          </label>
          <button type="submit">Get Insights</button>
        </form>

        <form onSubmit={handleVideoSubmit}>
          <label>
            Enter YouTube Video ID:
            <input
              type="text"
              value={videoId}
              onChange={(e) => setVideoId(e.target.value)}
              required
            />
          </label>
          <button type="submit">Get Insights</button>
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

      <footer className="App-footer">
        <img src={dega} alt="Scoop your insights from ocean of data" />
      </footer>
    </div>
  );
}

export default App;
