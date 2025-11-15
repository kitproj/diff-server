package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	randomPort := strconv.Itoa(rand.Intn(65535-49152) + 49152)

	port := flag.String("p", randomPort, "Port to listen on")
	workspaceDir := flag.String("C", ".", "Directory to scan for git repositories")
	flag.Parse()

	if err := os.Chdir(*workspaceDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change to directory %s: %v\n", *workspaceDir, err)
		os.Exit(1)
	}

	http.HandleFunc("/", diffsHandler)

	fmt.Printf("Starting server on http://localhost:%s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}
