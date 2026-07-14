#!/bin/sh
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

usage() {
    cat <<EOF
Usage: $(basename "$0") [options]

Options:
    --all           Run all tests (Go + Node)
    --go            Run Go unit tests only
    --node          Run Node integration tests only
    --watch         Run Go tests in watch mode (requires entr)
    --verbose       Show detailed output
    -h, --help      Show this help

Examples:
    $(basename "$0") --all
    $(basename "$0") --go --verbose
    $(basename "$0") --node
EOF
}

if [ $# -eq 0 ]; then
    set -- --all
fi

VERBOSE=""
RUN_GO=false
RUN_NODE=false
WATCH=false

while [ $# -gt 0 ]; do
    case "$1" in
        --all)
            RUN_GO=true
            RUN_NODE=true
            ;;
        --go)
            RUN_GO=true
            ;;
        --node)
            RUN_NODE=true
            ;;
        --watch)
            WATCH=true
            ;;
        --verbose)
            VERBOSE="-v"
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
    shift
done

cd "$PROJECT_ROOT"

if [ "$WATCH" = true ]; then
    if ! command -v entr >/dev/null 2>&1; then
        echo "ERROR: entr not found. Install with: brew install entr"
        exit 1
    fi
    echo "Watching for changes... (Ctrl+C to stop)"
    find .engine/scripts/src -name '*.go' | entr -r go test $VERBOSE ./...
    exit 0
fi

if [ "$RUN_GO" = true ]; then
    echo "Running Go unit tests..."
    cd .engine/scripts/src
    go test $VERBOSE ./...
    cd "$PROJECT_ROOT"
    echo "Go tests complete."
    echo ""
fi

if [ "$RUN_NODE" = true ]; then
    if [ ! -d tests ]; then
        echo "WARNING: tests/ directory not found. Skipping Node tests."
    elif [ ! -f tests/package.json ]; then
        echo "WARNING: tests/package.json not found. Skipping Node tests."
    else
        echo "Running Node integration tests..."
        cd tests
        if [ ! -d node_modules ]; then
            echo "Installing Node dependencies..."
            npm install
        fi
        npm test
        cd "$PROJECT_ROOT"
        echo "Node tests complete."
        echo ""
    fi
fi

echo "All tests passed."
