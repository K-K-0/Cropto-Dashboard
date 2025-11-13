import { TrendingDown, TrendingUp } from "lucide-react"
import type { TickerStats } from "../types"
import React from "react"

interface TickerStatsCardProps {
  ticker: TickerStats
}

const TickerStatsCard: React.FC<TickerStatsCardProps> = ({ ticker }) => {
  const formatPrice = (price: number) => {
  if (price >= 1_000_000_000) return `$${(price / 1_000_000_000).toFixed(2)}B`;
  if (price >= 1_000_000) return `$${(price / 1_000_000).toFixed(2)}M`;
  if (price >= 1_000) return `$${(price / 1_000).toFixed(2)}K`;

  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(price);
};

const formatNumber = (num: number) => {
  if (num >= 1_000_000_000) return (num / 1_000_000_000).toFixed(2) + "B";
  if (num >= 1_000_000) return (num / 1_000_000).toFixed(2) + "M";
  if (num >= 1_000) return (num / 1_000).toFixed(2) + "K";
  return Math.floor(num).toString();
};




  const isPositive = ticker.changePercent >= 0

  return (
    <div
      className={`rounded-2xl shadow-lg p-6 bg-white dark:bg-gray-900 border transition-all duration-300 
      hover:shadow-2xl hover:scale-[1.02] 
      ${isPositive ? "border-green-300" : "border-red-300"}
      `}
    >
  
      <div className="flex justify-between items-center mb-4">
        <span className="text-xl font-semibold tracking-wide text-gray-900 dark:text-gray-100">
          {ticker.symbol}
        </span>

        <span
          className={`flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium shadow-sm
          ${isPositive
            ? "bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400"
            : "bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400"
          }`}
        >
          {isPositive ? <TrendingUp size={18} /> : <TrendingDown size={18} />}
          {isPositive ? "+" : ""}
          {ticker.changePercent.toFixed(2)}%
        </span>
      </div>

      <div className="grid grid-cols-2 gap-4 text-sm">
        <Stat label="24h High" value={formatPrice(ticker.high)} />
        <Stat label="24h Low" value={formatPrice(ticker.low)} />
        <Stat label="Volume" value={formatNumber(ticker.volume)} />

        <Stat
          label="Change"
          value={formatPrice(Math.abs(ticker.change))}
          highlight
          positive={isPositive}
        />
      </div>

    
    </div>
  )
}

const Stat = ({
  label,
  value,
  highlight = false,
  positive = true,
}: {
  label: string
  value: string | number
  highlight?: boolean
  positive?: boolean
}) => (
  <div className="flex flex-col gap-1">
    <span className="text-gray-500 dark:text-gray-400">{label}</span>
    <span
      className={`text-base font-semibold 
      ${highlight ? (positive ? "text-green-600 dark:text-green-400" : "text-red-600 dark:text-red-400") : "text-gray-900 dark:text-gray-200"}
      `}
    >
      {value}
    </span>
  </div>
)

export default TickerStatsCard