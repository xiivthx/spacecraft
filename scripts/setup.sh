#!/bin/sh
set -e

echo "Spacecraft Developer Setup"
echo "=========================="

# Check Go version
if ! command -v go >/dev/null 2>&1; then
    echo "ERROR: Go not found. Install Go 1.21+ first."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Go version: $GO_VERSION"

# Check Node.js (optional, for integration tests)
if command -v node >/dev/null 2>&1; then
    NODE_VERSION=$(node --version)
    echo "Node.js version: $NODE_VERSION"
else
    echo "WARNING: Node.js not found. Integration tests won't run."
fi

# Build the binary
echo ""
echo "Building spacecraft binary..."
make build

# Verify binary
if [ ! -x scripts/spacecraft ]; then
    echo "ERROR: Build failed - scripts/spacecraft not found"
    exit 1
fi

echo "Build successful: scripts/spacecraft"

# Run tests
echo ""
echo "Running tests..."
make test

# Optional: install globally
echo ""
printf "Install globally? (y/N) "
read -r answer
case "$answer" in
    y|Y)
        make install
        echo "Installed. Restart OpenCode to apply config."
        ;;
    *)
        echo "Skipped global install."
        ;;
esac

echo ""
echo "Setup complete!"
echo "Run 'scripts/spacecraft help' to get started."
