import { Activity, Users, Zap, RefreshCw } from "lucide-react";
import { useWebsocket } from "./hooks/useWebsocket";
import StatsCard from "./components/statsCard";
import LivePriceCard from "./components/livefeed";
import TickerStatsCard from "./components/TickerStatsCard";

function App() {
  const { connected, livePrices, TickerStats, stats } = useWebsocket();

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950 text-gray-900 dark:text-gray-100 flex flex-col">
    
      <header className="sticky top-0 z-30 backdrop-blur-md bg-white/70 dark:bg-gray-900/70 border-b border-gray-200 dark:border-gray-800 shadow-sm">
        <div className="max-w-7xl mx-auto flex items-center justify-between p-5">
          <h1 className="flex items-center gap-2 text-2xl font-bold">
            <Activity className="text-blue-600 dark:text-blue-400" size={26} />
            <span>Crypto Live Dashboard</span>
          </h1>

          <div className="flex items-center gap-2 text-sm">
            <div
              className={`w-3 h-3 rounded-full ${
                connected ? "bg-emerald-500" : "bg-red-500"
              }`}
            ></div>
            <span
              className={`font-medium ${
                connected ? "text-emerald-600 dark:text-emerald-400" : "text-red-500"
              }`}
            >
              {connected ? "Connected" : "Disconnected"}
            </span>
          </div>
        </div>
      </header>

  
      <main className="flex-1 max-w-7xl mx-auto w-full p-6 space-y-10">
        
        <div>
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2 text-gray-800 dark:text-gray-200">
            <Zap className="text-yellow-500" />
            Realtime Stats
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
            <StatsCard
              icon={Activity}
              value={stats.totalMessages}
              label="Total Messages"
              variant="message"
            />
            <StatsCard
              icon={Zap}
              value={stats.tradeCount}
              label="Trade Events"
              variant="trades"
            />
            <StatsCard
              icon={RefreshCw}
              value={stats.tickerCount}
              label="Ticker Events"
              variant="tickers"
            />
            <StatsCard
              icon={Users}
              value={stats.updateRate}
              label="Updates/sec"
              variant="rate"
            />
          </div>
        </div>


        <section>
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2 text-gray-800 dark:text-gray-200">
            <Zap className="text-emerald-500" />
            Live Prices (Real-time Trades)
          </h2>
          {Object.keys(livePrices).length === 0 && connected ? (
            <div className="text-center py-12 text-gray-500 dark:text-gray-400 italic">
              Waiting for live price data...
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
              {Object.values(livePrices).map((ticker) => (
                <LivePriceCard key={ticker.symbol} ticker={ticker} />
              ))}
            </div>
          )}
        </section>


        <section>
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2 text-gray-800 dark:text-gray-200">
            <Activity className="text-purple-500" />
            24-Hour Statistics
          </h2>
          {Object.keys(TickerStats).length === 0 && connected ? (
            <div className="text-center py-12 text-gray-500 dark:text-gray-400 italic">
              Waiting for ticker statistics...
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
              {Object.values(TickerStats).map((ticker) => (
                <TickerStatsCard key={ticker.symbol} ticker={ticker} />
              ))}
            </div>
          )}
        </section>
      </main>


      <footer className="py-6 border-t border-gray-200 dark:border-gray-800 text-center text-sm text-gray-500 dark:text-gray-400">
        Powered by <span className="font-medium text-blue-600">WebSocket</span> •{" "}
        <span className="italic"></span>
      </footer>
    </div>
  );
}

export default App;