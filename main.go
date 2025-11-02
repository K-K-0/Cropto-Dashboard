package main

import (
	"context"
	"crypto-dashboard/exchange"
	"crypto-dashboard/server/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := websocket.NewHub()
	log.Println("Starting websocket Hub...")
	go hub.Run()

	symbols := []string{"btcusdt", "ethusdt", "bnbusdt", "solusdt"}
	binanceClient := *exchange.NewBinanceClient(symbols)
	binanceClient.Start()

	go bridgeExchangeToHub(ctx, binanceClient, hub)

	router := setupRouter(hub)

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("server has been started on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	binanceClient.Close()

	cancel()

	log.Println("Server exited")
}

func bridgeExchangeToHub(ctx context.Context, client exchange.BinanceClient, hub *websocket.Hub) {
	log.Println("Starting  exchange to hub bridge....")

	for {
		select {
		case <-ctx.Done():
			log.Println("Bridge shutting down...")
			return
		case message, ok := <-client.GetMessageChannel():
			if !ok {
				log.Println("Echange client channel closed")
				return
			}
			hub.Broadcast(message)
		}
	}
}

func setupRouter(hub *websocket.Hub) *gin.Engine {

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Status":    "healthy",
			"Clients":   hub.GetClientCount(),
			"timestamp": time.Now().Unix(),
		})
	})

	router.GET("/ws", func(ctx *gin.Context) {
		websocket.ServeWS(hub, ctx.Writer, ctx.Request)
	})

	api := router.Group("/api")
	{
		api.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "pong"})
		})

		api.GET("/stats", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"connected clients": hub.GetClientCount(),
				"uptime in second":  time.Since(startTime).Seconds(),
			})
		})
	}
	return router
}

var startTime = time.Now()
