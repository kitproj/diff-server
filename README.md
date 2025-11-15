# diff-server

A tiny Go program that shows git diffs in the current directory and all git repositories in subdirectories.

## Usage

```bash
# Build
go build -o diff-server .

# Run (default port 8080, current directory)
./diff-server

# Run on custom port
./diff-server -port 9000

# Scan a different directory
./diff-server -C /path/to/workspace
```

Then open http://localhost:8080 (or your custom port) in a browser to view the diffs.
# This is a test change
