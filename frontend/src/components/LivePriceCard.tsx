import React from "react";
import type { LivePrice } from "../types";
import { Activity } from "lucide-react";

interface LivePriceCardProps {
  ticker: LivePrice;
}

const LivePriceCard: React.FC<LivePriceCardProps> = ({ ticker }) => {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(price);
  };

  const getTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds}s ago`;
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    return `${hours}h ago`;
  };

  return (
    <div className="flex flex-col justify-between bg-white dark:bg-gray-900 shadow-sm hover:shadow-md rounded-2xl p-5 transition-all duration-300 hover:-translate-y-1 ring-1 ring-gray-200 dark:ring-gray-800">
      {/* Header */}
      <div className="flex items-center justify-between mb-3">
        <span className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {ticker.symbol}
        </span>
        <div className="flex items-center gap-1 px-2 py-1 rounded-full bg-green-100 dark:bg-green-900/40 text-green-700 dark:text-green-400 text-xs font-medium">
          <Activity size={12} />
          <span>LIVE</span>
        </div>
      </div>

      {/* Price */}
      <div className="text-3xl font-bold text-gray-900 dark:text-white tracking-tight">
        {formatPrice(ticker.price)}
      </div>

      {/* Timestamp */}
      <div className="text-sm text-gray-500 dark:text-gray-400 mt-3">
        Updated {getTimeAgo(ticker.timestamp)}
      </div>
    </div>
  );
};

export default LivePriceCard;
