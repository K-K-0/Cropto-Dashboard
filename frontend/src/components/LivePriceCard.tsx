import React from "react";
import type { LivePrice } from "../types";
import { Activity } from "lucide-react";

interface LivePriceCardProps {
  ticker: LivePrice;
}

const LivePriceCard: React.FC<LivePriceCardProps> = ({ ticker }) => {
  const formatPrice = (price: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(price);

  const getTimeAgo = (timestamp: number) => {
    if (!timestamp) return "—";
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (isNaN(seconds)) return "—";
    if (seconds < 60) return `${seconds}s ago`;
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    return `${hours}h ago`;
  };

  return (
    <div
      className="
      flex flex-col sm:flex-row sm:items-center sm:justify-between
      gap-3 sm:gap-4
      bg-white/5 dark:bg-gray-900
      rounded-xl p-4
      shadow-sm hover:shadow-md
      border border-gray-800/60
      transition-all duration-200 hover:bg-gray-800/40
      w-full
      overflow-hidden
      "
    >

      <div className="flex items-center gap-3 flex-shrink-0 min-w-0">
        <div
          className="
          flex items-center justify-center 
          w-10 h-10 sm:w-12 sm:h-12 
          rounded-lg bg-indigo-100/10 dark:bg-indigo-900/40 
          text-indigo-400 
          font-semibold text-lg sm:text-xl flex-shrink-0
        "
        >
          {ticker.symbol.charAt(0)}
        </div>

        <div className="min-w-0">
          <p className="text-base sm:text-lg font-semibold text-gray-100 truncate">
            {ticker.symbol}
          </p>
          <p className="text-xs text-gray-500 truncate">
            Updated {getTimeAgo(ticker.timestamp)}
          </p>
        </div>
      </div>


      <div className="flex flex-col items-end justify-center sm:items-center gap-1 sm:gap-2 text-right">
        <p
          className="
          text-lg sm:text-xl font-bold text-gray-100 
          break-all max-w-[110px] sm:max-w-none
          "
        >
          {formatPrice(ticker.price)}
        </p>

        <div
          className="
          flex items-center gap-1 
          px-2 py-0.5 rounded-full 
          bg-green-900/30 text-green-400 
          text-xs sm:text-sm font-semibold
          shadow-sm whitespace-nowrap
        "
        >
          <Activity size={12} className="animate-pulse" />
          <span>LIVE</span>
        </div>
      </div>
    </div>
  );
};

export default LivePriceCard;
