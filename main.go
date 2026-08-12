package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/pkg/browser"
)

func main() {
	port := flag.String("p", "3844", "Port to listen on")
	workspaceDir := flag.String("C", ".", "Directory to scan for git repositories")
	noBrowser := flag.Bool("B", false, "Disable auto-opening browser")
	flag.Parse()

	if err := os.Chdir(*workspaceDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change to directory %s: %v\n", *workspaceDir, err)
		os.Exit(1)
	}

	http.HandleFunc("/", diffsHandler)

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}

	url := fmt.Sprintf("http://localhost:%s", *port)
	fmt.Printf("Starting server on %s\n", url)

	if !*noBrowser {
		if err := browser.OpenURL(url); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open browser: %v\n", err)
		}
	}

	if err := http.Serve(listener, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}
