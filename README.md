# diffs-cli

A tiny Go program that shows git diffs in the current directory and all git repositories in subdirectories.

## Usage

```bash
# Build
go build -o diffs-cli .

# Run (default port 3844, current directory)
./diffs-cli

# Run on custom port
./diffs-cli -p 9000

# Scan a different directory
./diffs-cli -C /path/to/workspace
```

Then open http://localhost:3844 (or your custom port) in a browser to view the diffs.
