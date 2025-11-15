# diffs-cli

A tiny Go program that shows git diffs in the current directory and all git repositories in subdirectories.

## Usage

```bash
# Build
go build -o diffs-cli .

# Run (default random port, current directory)
./diffs-cli

# Run on custom port
./diffs-cli -p 9000

# Scan a different directory
./diffs-cli -C /path/to/workspace
```

Then open the URL shown in the output (e.g., http://localhost:52341) in a browser to view the diffs.
