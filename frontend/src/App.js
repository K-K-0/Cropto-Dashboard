import React, { useState, useEffect, useRef } from 'react';
import { TrendingUp, TrendingDown, Activity, Users, Zap, RefreshCw } from 'lucide-react';
import './App.css';

function App() {
  const [ws, setWs] = useState(null);
  const [connected, setConnected] = useState(false);
  const [livePrices, setLivePrices] = useState({});
  const [tickerStats, setTickerStats] = useState({});
  const [stats, setStats] = useState({
    totalMessages: 0,
    tradeCount: 0,
    tickerCount: 0,
    updateRate: 0
  });
  const [lastSecondCount, setLastSecondCount] = useState(0);
  const wsRef = useRef(null);

  useEffect(() => {
    connectWebSocket();
    
    const interval = setInterval(() => {
      setStats(prev => ({
        ...prev,
        updateRate: lastSecondCount
      }));
      setLastSecondCount(0);
    }, 1000);

    return () => {
      clearInterval(interval);
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  const connectWebSocket = () => {
    const websocket = new WebSocket('ws://localhost:8080/ws');
    
    websocket.onopen = () => {
      console.log('Connected to WebSocket');
      setConnected(true);
    };

    websocket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        
        setLastSecondCount(prev => prev + 1);
        setStats(prev => ({
          ...prev,
          totalMessages: prev.totalMessages + 1
        }));

        if (data.eventType === 'trade') {
          setStats(prev => ({ ...prev, tradeCount: prev.tradeCount + 1 }));
          setLivePrices(prev => ({
            ...prev,
            [data.symbol]: {
              symbol: data.symbol,
              price: parseFloat(data.price),
              timestamp: data.timestamp
            }
          }));
        } else if (data.eventType === 'ticker') {
          setStats(prev => ({ ...prev, tickerCount: prev.tickerCount + 1 }));
          setTickerStats(prev => ({
            ...prev,
            [data.symbol]: {
              symbol: data.symbol,
              price: parseFloat(data.price),
              change: parseFloat(data.change),
              changePercent: parseFloat(data.changePercent),
              high: parseFloat(data.high),
              low: parseFloat(data.low),
              volume: parseFloat(data.volume),
              timestamp: data.timestamp
            }
          }));
        }
      } catch (error) {
        console.error('Error parsing message:', error);
      }
    };

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    websocket.onclose = () => {
      console.log('WebSocket disconnected');
      setConnected(false);
      // Auto-reconnect after 3 seconds
      setTimeout(connectWebSocket, 3000);
    };

    wsRef.current = websocket;
    setWs(websocket);
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(price);
  };

  const formatNumber = (num) => {
    return new Intl.NumberFormat('en-US').format(Math.floor(num));
  };

  const getTimeAgo = (timestamp) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    return `${seconds}s ago`;
  };

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <h1>
            <Activity className="icon" />
            Crypto Live Dashboard
          </h1>
          <div className="connection-status">
            <div className={`status-dot ${connected ? 'connected' : 'disconnected'}`}></div>
            <span>{connected ? 'Connected' : 'Disconnected'}</span>
          </div>
        </div>
      </header>

      <div className="container">
        {/* Stats Cards */}
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-icon messages">
              <Activity size={24} />
            </div>
            <div className="stat-content">
              <div className="stat-value">{formatNumber(stats.totalMessages)}</div>
              <div className="stat-label">Total Messages</div>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-icon trades">
              <Zap size={24} />
            </div>
            <div className="stat-content">
              <div className="stat-value">{formatNumber(stats.tradeCount)}</div>
              <div className="stat-label">Trade Events</div>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-icon tickers">
              <RefreshCw size={24} />
            </div>
            <div className="stat-content">
              <div className="stat-value">{formatNumber(stats.tickerCount)}</div>
              <div className="stat-label">Ticker Events</div>
            </div>
          </div>

          <div className="stat-card">
            <div className="stat-icon rate">
              <Users size={24} />
            </div>
            <div className="stat-content">
              <div className="stat-value">{stats.updateRate}</div>
              <div className="stat-label">Updates/sec</div>
            </div>
          </div>
        </div>

        {/* Live Prices Section */}
        <section className="section">
          <h2>
            <Zap className="section-icon" />
            Live Prices (Real-time Trades)
          </h2>
          <div className="ticker-grid">
            {Object.values(livePrices).map((ticker) => (
              <div key={ticker.symbol} className="price-card live">
                <div className="card-header">
                  <span className="symbol">{ticker.symbol}</span>
                  <span className="live-badge">LIVE</span>
                </div>
                <div className="price">{formatPrice(ticker.price)}</div>
                <div className="timestamp">{getTimeAgo(ticker.timestamp)}</div>
              </div>
            ))}
          </div>
        </section>

        {/* 24h Statistics Section */}
        <section className="section">
          <h2>
            <Activity className="section-icon" />
            24-Hour Statistics
          </h2>
          <div className="ticker-grid">
            {Object.values(tickerStats).map((ticker) => (
              <div key={ticker.symbol} className="ticker-card">
                <div className="card-header">
                  <span className="symbol">{ticker.symbol}</span>
                </div>
                <div className="price">{formatPrice(ticker.price)}</div>
                <div className={`change ${ticker.changePercent >= 0 ? 'positive' : 'negative'}`}>
                  {ticker.changePercent >= 0 ? (
                    <TrendingUp size={16} />
                  ) : (
                    <TrendingDown size={16} />
                  )}
                  {ticker.changePercent >= 0 ? '+' : ''}{ticker.changePercent.toFixed(2)}%
                </div>
                <div className="stats-grid-mini">
                  <div className="stat-mini">
                    <span className="label">24h High</span>
                    <span className="value">{formatPrice(ticker.high)}</span>
                  </div>
                  <div className="stat-mini">
                    <span className="label">24h Low</span>
                    <span className="value">{formatPrice(ticker.low)}</span>
                  </div>
                  <div className="stat-mini">
                    <span className="label">Volume</span>
                    <span className="value">{formatNumber(ticker.volume)}</span>
                  </div>
                  <div className="stat-mini">
                    <span className="label">Change</span>
                    <span className={`value ${ticker.change >= 0 ? 'positive' : 'negative'}`}>
                      {formatPrice(Math.abs(ticker.change))}
                    </span>
                  </div>
                </div>
                <div className="timestamp">{getTimeAgo(ticker.timestamp)}</div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}

export default App;