# diff-server

A lightweight Go web server that displays git diffs from the current directory and all git repositories in subdirectories. It provides a clean, real-time web interface for viewing uncommitted changes across multiple repositories in your workspace.

Like `jq`, it is a single tiny binary without dependencies, making it easy to install and use anywhere.

## Installation

### Supported Platforms

Binaries are available for:
- **Linux**: 386, amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)

### Download and Install

Download the binary for your platform from the [release page](https://github.com/kitproj/diff-server/releases).

#### Linux

**For Linux (amd64):**
```bash
sudo curl -fsL -o /usr/local/bin/diff-server https://github.com/kitproj/diff-server/releases/download/v0.0.1/diff-server_v0.0.1_linux_amd64
sudo chmod +x /usr/local/bin/diff-server
```

**For Linux (arm64):**
```bash
sudo curl -fsL -o /usr/local/bin/diff-server https://github.com/kitproj/diff-server/releases/download/v0.0.1/diff-server_v0.0.1_linux_arm64
sudo chmod +x /usr/local/bin/diff-server
```

**For Linux (386):**
```bash
sudo curl -fsL -o /usr/local/bin/diff-server https://github.com/kitproj/diff-server/releases/download/v0.0.1/diff-server_v0.0.1_linux_386
sudo chmod +x /usr/local/bin/diff-server
```

#### macOS

**For macOS (Apple Silicon/arm64):**
```bash
sudo curl -fsL -o /usr/local/bin/diff-server https://github.com/kitproj/diff-server/releases/download/v0.0.1/diff-server_v0.0.1_darwin_arm64
sudo chmod +x /usr/local/bin/diff-server
```

**For macOS (Intel/amd64):**
```bash
sudo curl -fsL -o /usr/local/bin/diff-server https://github.com/kitproj/diff-server/releases/download/v0.0.1/diff-server_v0.0.1_darwin_amd64
sudo chmod +x /usr/local/bin/diff-server
```

#### Verify Installation

After installing, verify the installation works:
```bash
diff-server -h
```

### Build from Source

If you prefer to build from source:

```bash
git clone https://github.com/kitproj/diff-server.git
cd diff-server
go build -o diff-server .
```

## Usage

### Starting the Server

**Run with default settings (port 3844, current directory):**
```bash
diff-server
```

**Run on a custom port:**
```bash
diff-server -p 9000
```

**Scan a different directory:**
```bash
diff-server -C /path/to/workspace
```

**Combine options:**
```bash
diff-server -p 8080 -C ~/projects
```

### Viewing Diffs

After starting the server, open your web browser and navigate to:
```
http://localhost:3844
```

The web interface will automatically:
- Display all uncommitted changes in the current directory (if it's a git repository)
- Recursively scan subdirectories for git repositories
- Show diffs from all discovered repositories
- Auto-refresh every 10 seconds to show new changes

### Command-Line Options

```bash
Usage of diff-server:
  -C string
    	Directory to scan for git repositories (default ".")
  -p string
    	Port to listen on (default "3844")
```

## Features

- **Multi-repository support**: Automatically discovers and displays diffs from all git repositories in subdirectories
- **Real-time updates**: The web interface polls for changes every 10 seconds
- **Clean UI**: Uses Diff2Html for a beautiful, syntax-highlighted diff view
- **Untracked files**: Shows new files that haven't been added to git yet
- **Lightweight**: Single binary with no runtime dependencies
- **Fast**: Efficient scanning and rendering of diffs

## Use Cases

- **Code review**: Quickly review all changes across multiple projects before committing
- **Workspace monitoring**: Keep an eye on what's changed in your development workspace
- **Pair programming**: Share a link to show your current work to teammates
- **CI/CD**: Run in a pipeline to visualize changes before deployment

## Troubleshooting

### Common Issues

**"Cannot connect" or server not accessible**
- Verify the server is running: check the terminal for the startup message
- Ensure the port is not already in use by another application
- Try a different port: `diff-server -p 8080`

**No diffs showing**
- Ensure you're in a directory that contains git repositories
- Check that you have uncommitted changes: `git status`
- Verify the directory path if using `-C` flag

**Large diffs not fully displayed**
- The server limits output to 5MB per request to prevent memory issues
- Consider committing some changes or viewing specific repositories separately

**Permission denied when installing**
- Use `sudo` when installing to `/usr/local/bin`
- Alternatively, install to a user-owned directory like `~/bin` and add it to your PATH

### Getting Help

- Report issues: https://github.com/kitproj/diff-server/issues
- Check existing issues for solutions and workarounds
