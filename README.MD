📈 Crypto Market Backend — Go (Gin)

A production-ready backend built using Golang + Gin that provides real-time crypto market data, ticker data, and historical candlestick data.
Integrated with Docker, Docker Compose, and CI/CD pipelines, ready for deployment.

🚀 Features
✅ Market Data (Live Trades)

Fetches real-time trade information for any supported crypto pair.

✅ Ticker Data

Get ticker price updates from Binance (24hr statistics etc.).

✅ Historical Candle Data

Supports historical OHLCV (Open, High, Low, Close, Volume) data for any interval.

🔧 Tech Stack

Go (Golang)

Gin Framework

Binance API

Docker + Docker Compose

CI/CD Pipeline (GitHub Actions / GitLab / etc.)


🛠 Setup Instructions
1. Clone the repository
git clone https://github.com/yourusername/crypto-backend.git
cd crypto-backend

▶️ Run Locally (Without Docker)
Install dependencies
go mod tidy

Start the server
go run main.go

Local API Base URL
http://localhost:8080

🐳 Run with Docker
Build Docker image
docker build -t crypto-backend .

Run container
docker run -p 8080:8080 crypto-backend

📌 API Endpoints
1. Live Trade Data
POST /api/live


Request Body:

{
  "symbol": "BTCUSDT"
}

2. Ticker Data
POST /api/ticker


Request Body:

{
  "symbol": "BTCUSDT"
}

3. Historical Candlestick Data
POST /api/historical


Request Body:

{
  "symbol": "BTCUSDT",
  "interval": "1h",
  "limit": 200
}

4. Multiple Symbols (Charts)
POST /api/chart/multiple


Request Example:

{
  "symbols": ["BTCUSDT", "ETHUSDT"],
  "interval": "1h",
  "limit": 100
}

🧪 Testing With Postman
Local API URL
http://localhost:8080/api/<endpoint>


Example for historical data:

http://localhost:8080/api/historical


Set Content-Type: application/json
Add correct symbol → e.g. "BTCUSDT" (uppercase, no spaces)

⚠️ Common Errors & Fixes
❌ Error: "Invalid symbol"

Binance requires:

Uppercase symbols

Without spaces or slashes

Example → BTCUSDT, ETHUSDT, SOLUSDT
❌ “BTC-USDT”
❌ “btc/usdt”
❌ “ btcusdt ”