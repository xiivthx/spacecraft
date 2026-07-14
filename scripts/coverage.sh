#!/bin/sh
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

usage() {
    cat <<EOF
Usage: $(basename "$0") [options]

Options:
    --html          Generate HTML coverage report (default)
    --text          Print text coverage summary
    --open          Open HTML report in browser
    -h, --help      Show this help

Examples:
    $(basename "$0")
    $(basename "$0") --text
    $(basename "$0") --html --open
EOF
}

OPEN_BROWSER=false
OUTPUT_FORMAT="html"

while [ $# -gt 0 ]; do
    case "$1" in
        --html)
            OUTPUT_FORMAT="html"
            ;;
        --text)
            OUTPUT_FORMAT="text"
            ;;
        --open)
            OPEN_BROWSER=true
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

cd "$PROJECT_ROOT/.engine/scripts/src"

echo "Running tests with coverage..."
COVERAGE_FILE="coverage.out"
go test -coverprofile="$COVERAGE_FILE" ./...

if [ "$OUTPUT_FORMAT" = "text" ]; then
    echo ""
    echo "Coverage summary:"
    go tool cover -func="$COVERAGE_FILE" | tail -1
else
    HTML_FILE="coverage.html"
    echo "Generating HTML coverage report..."
    go tool cover -html="$COVERAGE_FILE" -o "$HTML_FILE"
    echo "Report generated: $HTML_FILE"

    if [ "$OPEN_BROWSER" = true ]; then
        if command -v open >/dev/null 2>&1; then
            open "$HTML_FILE"
        elif command -v xdg-open >/dev/null 2>&1; then
            xdg-open "$HTML_FILE"
        else
            echo "WARNING: Could not open browser. Open $HTML_FILE manually."
        fi
    fi
fi

echo ""
echo "Coverage report complete."
