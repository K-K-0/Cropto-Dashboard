import { TrendingDown, TrendingUp } from "lucide-react";
import type { TickerStats } from "../types";

interface TickerStatsCardsProps {
    ticker: TickerStats
}

const TickerStatsCard: React.FC<TickerStatsCardsProps> = ({ ticker }) => {
    const formatPrice = (price: number) => {
         return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(price);
    }

    const formatNumber = (num: number) => {
        return new Intl.NumberFormat('en-US').format(Math.floor(num))
    }

    const getTimeAgo = (timestamp: number) => {
        const seconds = Math.floor((Date.now() - timestamp) / 1000)
        return `${seconds}s ago`
    }

    const isPositive = ticker.changePercent >= 0

    return (
        <div className="ticker-card">
            <div className="card-header">
                <span className="symbol">{ticker.symbol}</span>
            </div>

            <div className="price">{formatPrice(ticker.price)}</div>
            <div className={`change ${isPositive ? 'positive' : 'negative'}`}>
                {isPositive ? <TrendingUp size={16} /> : <TrendingDown size={16}/>}
                {isPositive ? '+' : ''}{ticker.changePercent.toFixed(2)}%
            </div>
            <div className="stats-grid-mini">
                <div className="stas-mini">
                    <span className="label">24h High</span>
                    <span className="value">{formatPrice(ticker.high)}</span>
                </div>
                <div className="stas-mini">
                    <span className="label">24h Low</span>
                    <span className="value">{formatPrice(ticker.low)}</span>
                </div>
                <div className="stas-mini">
                    <span className="label">Volume</span>
                    <span className="value">{formatNumber(ticker.volume)}</span>
                </div>
                <div className="stas-mini">
                    <span className="label">change</span>
                    <span className={`value ${isPositive ? 'positive' : 'negative'}`}>
                        {formatPrice(Math.abs(ticker.change))}
                    </span>
                </div>
            </div>
            <div className="timestamp"> {getTimeAgo(ticker.timestamp)}</div>
        </div>
    )

}

export default TickerStatsCard