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
    <div className="
      flex flex-col justify-between 
      bg-white/80 dark:bg-gray-900/60 
      backdrop-blur-xl
      shadow-lg hover:shadow-xl 
      rounded-2xl p-6 
      transition-all duration-300 hover:-translate-y-1 
      ring-1 ring-gray-200/50 dark:ring-gray-800/60
      border border-transparent hover:border-indigo-400/30
    ">
      
      <div className="flex items-center justify-between mb-4">
        <span className="text-xl font-semibold text-gray-900 dark:text-gray-100 tracking-tight">
          {ticker.symbol}
        </span>

        <div className="
          flex items-center gap-1.5 
          px-3 py-1.5 
          rounded-full 
          bg-gradient-to-r from-green-100 to-green-200 dark:from-green-900/40 dark:to-green-700/30 
          text-green-700 dark:text-green-400 
          text-xs font-semibold
          shadow-sm
          animate-pulse
        ">
          <Activity size={12} className="animate-ping" />
          <span>LIVE</span>
        </div>
      </div>

      <div className="text-4xl font-extrabold text-gray-900 dark:text-white tracking-tight mb-2">
        {formatPrice(ticker.price)}
      </div>

      <div className="text-sm text-gray-500 dark:text-gray-400 mt-1">
        Updated <span className="font-medium">{getTimeAgo(ticker.timestamp)}</span>
      </div>
    </div>
  );
};

export default LivePriceCard;
