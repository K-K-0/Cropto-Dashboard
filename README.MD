📈 Crypto Market Backend — Golang (Gin)

A production-ready backend built using Golang + Gin to fetch live trade data, ticker data, and historical candlestick (OHLCV) data directly from Binance.
Includes Docker, Docker Compose, and CI/CD pipeline support for easy deployment to AWS/GCP/Render/Railway.

🚀 Features

Live Crypto Trade Data (real-time market updates)

Ticker Data (24-hour price, volume, change)

Historical Candlestick Data (1m, 5m, 1h, 4h, 1d intervals)

Multiple Symbols Historical API

REST API using Gin

Binance Integration

Production-ready Docker Setup

CI/CD Ready

Clean Project Structure

📁 Project Structure
/api            → All API route handlers
/services       → Binance service & data fetch logic
/utils          → Helper utilities
/main.go        → App entrypoint
/Dockerfile     → Production container
/go.mod         → Dependencies

🛠️ Run Locally
1. Clone the project
git clone <repo-url>
cd crypto-backend

2. Install dependencies
go mod tidy

3. Start server
go run main.go


Server will run on:

http://localhost:8080

🐳 Run With Docker
Build image
docker build -t crypto-backend .

Run container
docker run -p 8080:8080 crypto-backend

🔌 API Endpoints
1. Live Trades (POST)

URL: /api/trade

Body:

{ "symbol": "BTCUSDT" }

2. Ticker Data (POST)

URL: /api/ticker

Body:

{ "symbol": "ETHUSDT" }

3. Historical Data (POST)

Fetch OHLCV candlestick data.

URL: /api/historical

Body:

{
  "symbol": "BTCUSDT",
  "interval": "1h",
  "limit": 100
}

4. Multiple Symbols Historical (POST)

URL: /api/chart/multiple

Body:

{
  "symbols": ["BTCUSDT", "ETHUSDT", "SOLUSDT"],
  "interval": "1h",
  "limit": 50
}

❗ Important Notes

All symbols must be uppercase
BTCUSDT, ETHUSDT, SOLUSDT, etc.

Binance returns "Invalid symbol" if:

lowercase is used

wrong pair (e.g., BTC/USD)

extra characters (spaces, slashes)

👨‍🔧 Testing With Postman

Select POST method

Set URL (example):

http://localhost:8080/api/historical


Open Body → raw → JSON

Paste request body

Send the request