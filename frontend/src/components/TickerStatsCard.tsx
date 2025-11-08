import { TrendingDown, TrendingUp } from "lucide-react"
import type { TickerStats } from "../types"
import React from "react"

interface TickerStatsCardProps {
  ticker: TickerStats
}

const TickerStatsCard: React.FC<TickerStatsCardProps> = ({ ticker }) => {
  const formatPrice = (price: number) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(price)

  const formatNumber = (num: number) =>
    new Intl.NumberFormat("en-US").format(Math.floor(num))

  const getTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000)
    if (seconds < 60) return `${seconds}s ago`
    const minutes = Math.floor(seconds / 60)
    if (minutes < 60) return `${minutes}m ago`
    const hours = Math.floor(minutes / 60)
    return `${hours}h ago`
  }

  const isPositive = ticker.changePercent >= 0

  return (
    <div
      className={`rounded-2xl shadow-md p-5 transition-all duration-300 bg-white dark:bg-gray-900 hover:shadow-xl border flex flex-col justify-between ${
        isPositive
          ? "border-green-300 hover:border-green-400"
          : "border-red-300 hover:border-red-400"
      }`}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-3">
        <span className="text-lg font-semibold text-gray-800 dark:text-gray-200">
          {ticker.symbol}
        </span>
        <span
          className={`flex items-center gap-1 px-2 py-0.5 rounded-full text-sm font-medium whitespace-nowrap ${
            isPositive
              ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
              : "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
          }`}
        >
          {isPositive ? (
            <TrendingUp size={16} className="inline-block" />
          ) : (
            <TrendingDown size={16} className="inline-block" />
          )}
          {isPositive ? "+" : ""}
          {ticker.changePercent.toFixed(2)}%
        </span>
      </div>

      {/* Price */}
      <div className="text-2xl sm:text-3xl font-bold mb-5 text-gray-900 dark:text-gray-100 break-words">
        {/* {formatPrice(ticker.price)} */}
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
        <div className="flex flex-col">
          <span className="text-gray-500 dark:text-gray-400">24h High</span>
          <span className="font-medium text-gray-900 dark:text-gray-200">
            {formatPrice(ticker.high)}
          </span>
        </div>
        <div className="flex flex-col">
          <span className="text-gray-500 dark:text-gray-400">24h Low</span>
          <span className="font-medium text-gray-900 dark:text-gray-200">
            {formatPrice(ticker.low)}
          </span>
        </div>
        <div className="flex flex-col">
          <span className="text-gray-500 dark:text-gray-400">Volume</span>
          <span className="font-medium text-gray-900 dark:text-gray-200 truncate">
            {formatNumber(ticker.volume)}
          </span>
        </div>
        <div className="flex flex-col">
          <span className="text-gray-500 dark:text-gray-400">Change</span>
          <span
            className={`font-medium ${
              isPositive
                ? "text-green-600 dark:text-green-400"
                : "text-red-600 dark:text-red-400"
            }`}
          >
            {formatPrice(Math.abs(ticker.change))}
          </span>
        </div>
      </div>

      {/* Timestamp */}
      <div className="mt-4 text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap">
        Updated {getTimeAgo(ticker.timestamp)}
      </div>
    </div>
  )
}

export default TickerStatsCard
