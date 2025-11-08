import React from "react"
import { type LucideIcon } from "lucide-react"

interface StatsCardProps {
  icon: LucideIcon
  value: number | string
  label: string
  variant: "message" | "trades" | "tickers" | "rate"
}

const StatsCard: React.FC<StatsCardProps> = ({
  icon: Icon,
  value,
  label,
  variant,
}) => {
  const formatNumber = (num: number | string) => {
    if (typeof num === "string") return num
    return new Intl.NumberFormat("en-US").format(num)
  }

  const variantColors: Record<
    StatsCardProps["variant"],
    { bg: string; text: string; ring: string }
  > = {
    message: {
      bg: "bg-blue-100 dark:bg-blue-900/30",
      text: "text-blue-600 dark:text-blue-400",
      ring: "ring-blue-300 dark:ring-blue-700",
    },
    trades: {
      bg: "bg-emerald-100 dark:bg-emerald-900/30",
      text: "text-emerald-600 dark:text-emerald-400",
      ring: "ring-emerald-300 dark:ring-emerald-700",
    },
    tickers: {
      bg: "bg-purple-100 dark:bg-purple-900/30",
      text: "text-purple-600 dark:text-purple-400",
      ring: "ring-purple-300 dark:ring-purple-700",
    },
    rate: {
      bg: "bg-amber-100 dark:bg-amber-900/30",
      text: "text-amber-600 dark:text-amber-400",
      ring: "ring-amber-300 dark:ring-amber-700",
    },
  }

  const colors = variantColors[variant]

  return (
    <div className="flex items-center gap-4 p-5 rounded-2xl bg-white dark:bg-gray-900 shadow-sm hover:shadow-md transition-all duration-300 ring-1 ring-gray-200 dark:ring-gray-800 hover:-translate-y-1">
      <div
        className={`flex items-center justify-center w-12 h-12 rounded-xl ${colors.bg} ${colors.text} ring-2 ${colors.ring} shadow-inner`}
      >
        <Icon size={24} />
      </div>

      <div className="flex flex-col">
        <span className="text-xl font-semibold text-gray-900 dark:text-gray-100">
          {formatNumber(value)}
        </span>
        <span className="text-sm text-gray-500 dark:text-gray-400">
          {label}
        </span>
      </div>
    </div>
  )
}

export default StatsCard
