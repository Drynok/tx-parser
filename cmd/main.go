package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Drynok/tx-parser/internal/api"
	"github.com/Drynok/tx-parser/internal/parser"
	rpc "github.com/Drynok/tx-parser/internal/rpc"
	"github.com/Drynok/tx-parser/internal/storage"
	"github.com/Drynok/tx-parser/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://ethereum-rpc.publicnode.com"
	}

	storage := storage.NewMemoryStorage()
	client := rpc.NewClient(rpcURL)
	logger := logger.NewLogger()
	ctx := context.Background()

	parser := parser.NewEthereumParser(client, storage, logger)

	// Parser start.
	parser.Start(ctx)

	// API endpoints.
	handler := api.NewHandler(parser)

	// Gin router
	router := gin.Default()

	router.GET("/current-block", handler.GetCurrentBlock)
	router.POST("/subscribe", handler.Subscribe)
	router.GET("/transactions", handler.GetTransactions)

	// Get the port from the environment variable or default to ":8080"
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	logger.Info("Starting server on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error("Failed to start server", "error", err)
	}

	srv := &http.Server{
		Addr: ":" + port,
	}

	go func() {
		logger.Info("Starting server on port " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	// Gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}
