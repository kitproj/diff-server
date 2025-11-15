package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	workspaceDir := flag.String("C", ".", "Directory to scan for git repositories")
	flag.Parse()

	if err := os.Chdir(*workspaceDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change to directory %s: %v\n", *workspaceDir, err)
		os.Exit(1)
	}

	http.HandleFunc("/", diffsHandler)

	server := &http.Server{
		Addr: ":" + *port,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting server on http://localhost:%s\n", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-stop

	fmt.Println("\nShutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server shutdown error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped")
}
