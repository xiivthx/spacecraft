#!/bin/sh
# bootstrap.sh — install spacecraft into any project
# Usage: curl -fsSL https://raw.githubusercontent.com/xiivthx/spacecraft/main/bootstrap.sh | sh
#    or: ./bootstrap.sh [project-dir]

set -e

TARGET="${1:-.}"
REPO="https://raw.githubusercontent.com/xiivthx/spacecraft/main"

echo "Spacecraft bootstrap"
echo "===================="

# Create directory structure
mkdir -p "$TARGET/.cursor/agents"
mkdir -p "$TARGET/.cursor/rules"
mkdir -p "$TARGET/.space/missions"
mkdir -p "$TARGET/.space/archive"

# Download rules
echo "Fetching rules..."
for rule in 050-style 000-spacecraft 100-conventions 200-workflow 300-security 400-performance 500-database 600-firmware 610-firmware-peripherals 620-firmware-testing; do
  curl -fsSL "$REPO/.cursor/rules/${rule}.mdc" -o "$TARGET/.cursor/rules/${rule}.mdc"
  echo "  $rule.mdc"
done

# Download agents
echo "Fetching agents..."
for agent in sc-coder sc-tester sc-planner sc-reviewer sc-designer sc-adviser sc-firmware; do
  curl -fsSL "$REPO/.cursor/agents/${agent}.md" -o "$TARGET/.cursor/agents/${agent}.md"
  echo "  $agent"
done

# Download MCP config
echo "Fetching mcp.json..."
curl -fsSL "$REPO/.cursor/mcp.json" -o "$TARGET/.cursor/mcp.json"
echo "  mcp.json"

# Download CLI binary (macOS arm64)
echo "Downloading spacecraft CLI..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$OS" = "darwin" ] && [ "$ARCH" = "arm64" ]; then
  curl -fsSL "$REPO/spacecraft" -o "$TARGET/spacecraft"
  chmod +x "$TARGET/spacecraft"
  echo "  spacecraft (darwin/arm64)"
else
  echo "  skipped — prebuilt binary only for macOS arm64. Build from source: cd cmd/spacecraft && go build ."
fi

echo ""
echo "Done. $TARGET is now spacecraft-ready."
echo ""
echo "To install globally:"
echo "  ln -sf \$(pwd)/spacecraft ~/.local/bin/spacecraft"
echo "  cp .cursor/agents/sc-*.md ~/.cursor/agents/"
echo ""
echo "Restart Cursor to pick up config."
