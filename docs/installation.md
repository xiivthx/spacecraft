# Installation Guide

This guide covers all supported installation methods for Spacecraft.

## Requirements

- **Go 1.21 or later** (for building from source)
- **macOS or Linux** (Windows not yet supported)
- **OpenCode CLI** installed and configured
- **Git** for version control operations

## Method 1: Build from Source

Recommended for developers who want the latest version or plan to contribute. This method builds from source using the Go toolchain.

```sh
# Clone the repository
git clone <repo-url>
cd spacecraft

# Build the Go binary
make build

# Verify the build
scripts/spacecraft help

# Optional: Install globally
make install
```

The `make build` command compiles the Go helper binary to `scripts/spacecraft`. When you build from source, you get the latest changes. The `make install` command creates a symlink in `~/.local/bin/` and generates the OpenCode configuration.

After running `make install`, restart OpenCode to load the new configuration.

## Method 2: Binary Download

Pre-built binaries are available for macOS and Linux.

1. Download the latest release from the releases page
2. Extract the archive
3. Move the `spacecraft` binary to a directory in your PATH:

```sh
# macOS / Linux
chmod +x spacecraft
sudo mv spacecraft /usr/local/bin/

# Or to user-local bin
mkdir -p ~/.local/bin
mv spacecraft ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

Verify the installation:

```sh
spacecraft help
```

## Method 3: Homebrew (macOS/Linux)

Homebrew installation provides automatic updates and dependency management.

```sh
# Add the tap (if using a custom tap)
brew tap <org>/spacecraft

# Install
brew install spacecraft

# Verify
spacecraft help
```

Update to the latest version:

```sh
brew upgrade spacecraft
```

## Post-Installation Setup

After installing, initialize Spacecraft in your project:

```sh
cd your-project
spacecraft init
```

This creates the `.space/` directory structure for mission state, artifacts, and evidence.

### Verify Installation

Run these commands to verify everything is working:

```sh
# Check CLI is accessible
spacecraft help

# Initialize in a test project
mkdir test-project && cd test-project
spacecraft init

# Create a test mission
spacecraft new "Test mission"

# List missions
spacecraft missions

# Check status
spacecraft status
```

If all commands succeed, Spacecraft is installed correctly.

## Troubleshooting

### "command not found: spacecraft"

- Verify the binary is in your PATH: `which spacecraft`
- If built from source, ensure `scripts/spacecraft` exists and is executable
- If installed via `make install`, check `~/.local/bin/` is in your PATH

### "permission denied" when running scripts

```sh
chmod +x scripts/spacecraft
```

### Build fails with Go errors

- Verify Go version: `go version` (requires 1.21+)
- Clear Go cache: `go clean -cache`
- Rebuild: `make build`

### OpenCode doesn't recognize slash commands

- Restart OpenCode after running `make install`
- Verify config exists: `cat ~/.config/opencode/opencode.jsonc`
- Check that `.engine/` directory is present in your project

## Uninstallation

```sh
# If installed via make install
make uninstall

# Or manually
rm ~/.local/bin/spacecraft
rm -rf ~/.config/opencode

# Remove from project
rm -rf .space/
```

The `make uninstall` command removes the symlink but leaves the config file intact. Delete it manually if you want a complete removal.
