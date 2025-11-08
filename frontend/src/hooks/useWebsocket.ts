import { useCallback, useEffect, useRef, useState } from 'react'
import type { LivePriceMap, Stats, TickerStatsMap, webSocketMessage } from '../types'


export const useWebsocket = () => {
    const [ connected, setConnected ] = useState(false)
    const [ livePrices, setlivePrices ] = useState<LivePriceMap>({})
    const [ TickerStats, setTickerStats ] = useState<TickerStatsMap>({})
    const [ stats, setStats ] = useState<Stats>({
        totalMessages: 0,
        tradeCount: 0,
        tickerCount: 0,
        updateRate: 0
    })

    const wsRef = useRef<WebSocket | null>(null)
    const lastSecondCountRef = useRef(0)
    const reconnectTimeoutRef = useRef<number | null>(null)

    const connect = useCallback(() => {
        try {
            const ws = new WebSocket("ws://localhost:8080/ws")

            ws.onopen = () => {
                setConnected(true)
                if (reconnectTimeoutRef.current) {
                    clearTimeout(reconnectTimeoutRef.current)
                }
            }

            ws.onmessage = (event) => {
                try {
                    const data: webSocketMessage = JSON.parse(event.data)
                    lastSecondCountRef.current ++

                    setStats((prev) => ({
                        ...prev,
                        totalMessages: prev.totalMessages + 1
                    }))

                    if (data.eventType == "trade") {
                        setStats((prev) => ({
                            ...prev, tradeCount: prev.tradeCount + 1
                        }))
                        setlivePrices((prev) => ({
                            ...prev, [data.symbol]: {
                                symbol: data.symbol,
                                price: parseFloat(data.price),
                               timestamp: data.timeStamp
                            }
                        }))
                    } else if (data.eventType == "ticker") {
                        setStats((prev) => ({ ...prev, tickerCount: prev.tickerCount + 1}))
                        setTickerStats((prev) => ({
                            ...prev, [data.symbol]: {
                                symbol: data.symbol,
                                price: parseFloat(data.price),
                                change: parseFloat(data.change),
                                changePercent: parseFloat(data.changePercent),
                                high: parseFloat(data.high),
                                low: parseFloat(data.low),
                                volume: parseFloat(data.volume),
                                timestamp: data.timeStamp
                            },
                        }))
                    }
                } catch (error) {
                    console.error('Error parsing message:', error);
                }
            }

              ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };

            ws.onclose = () => {
                console.log('🔌 WebSocket disconnected');
                setConnected(false);
                wsRef.current = null;

                reconnectTimeoutRef.current = setTimeout(() => {
                console.log('Attempting to reconnect...');
                connect();
                }, 3000);
            };

            wsRef.current = ws;

        } catch (error) {
            console.error(' Connection error:', error);
        }
    }, [])

    useEffect(() => {
        connect()

        const rateInterval = setInterval(() => {
           setStats((prev) => ({
            ...prev, 
            updateRate: lastSecondCountRef.current,
           })) 
           lastSecondCountRef.current = 0
        }, 1000);
        return () => {
            clearInterval(rateInterval)
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current)
            }
            if (wsRef.current) {
                wsRef.current.close()
            }
        }
    }, [connect])

    return { connected, livePrices, TickerStats, stats}

}