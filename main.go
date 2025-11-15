package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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

	fmt.Printf("Starting server on http://localhost:%s\n", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}
