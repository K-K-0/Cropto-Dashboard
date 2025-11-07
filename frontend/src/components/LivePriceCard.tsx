import type { LivePrice } from "../types";
import React from "react";

interface LivePriceCardProps {
    ticker: LivePrice
}

const LivePriceCard: React.FC<LivePriceCardProps> = ({ ticker }) => {
    const formatPrice = (price: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(price)
    }

    const getTimeAgo = (timestamp: number) => {
        const seconds = Math.floor((Date.now() - timestamp) / 1000)
        return `${seconds}s ago`
    }

    return (
        <div className="price-card live">
            <div className="card-header">
                <span className="symbol">{ticker.symbol}</span>
                <span className="live-badge">LIVE</span>
            </div>
            <div className="price">{formatPrice(ticker.price)}</div>
            <div className="timestamp">{getTimeAgo(ticker.timestamp)}</div>
            </div>
    )
}

export default LivePriceCard