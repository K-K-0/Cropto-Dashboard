export interface webSocketMessage {
    symbol: string
    price: string
    change: string
    changePercent: string
    volume: string
    high: string
    low: string
    timeStamp: number
    eventType: "trade" | "ticker"
}

export interface LivePrice {
    symbol: string
    price: number
    timestamp: number
}

export interface TickerStats {
    symbol: string
    price: number
    change: number
    changePercent: number
    high: number
    low: number
    volume: number
    timestamp: number
}

export interface Stats {
    totalMessages: number
    tradeCount: number
    tickerCount: number
    updateRate: number
}

export type LivePriceMap = Record<string, LivePrice>
export type TickerStatsMap = Record<string, TickerStats>
