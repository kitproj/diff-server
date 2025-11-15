# diffs-cli

A tiny Go program that shows git diffs in the current directory and all git repositories in subdirectories.

## Usage

```bash
# Build
go build -o diffs-cli .

# Run (default port 8080)
./diffs-cli

# Run on custom port
./diffs-cli 9000
```

Then open http://localhost:8080 (or your custom port) in a browser to view the diffs.
