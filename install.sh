#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="${PREFIX:-$HOME/.local/bin}"
BINARY_NAME="ghostwriter"

echo ""
echo "  ██████╗ ██╗  ██╗ ██████╗ ███████╗████████╗██╗    ██╗██████╗ ██╗████████╗███████╗██████╗ "
echo "  ██╔════╝██║  ██║██╔═══██╗██╔════╝╚══██╔══╝██║    ██║██╔══██╗██║╚══██╔══╝██╔════╝██╔══██╗"
echo "  ██║  ███╗███████║██║   ██║███████╗   ██║   ██║ █╗ ██║██████╔╝██║   ██║   █████╗  ██████╔╝"
echo "  ██║   ██║██╔══██║██║   ██║╚════██║   ██║   ██║███╗██║██╔══██╗██║   ██║   ██╔══╝  ██╔══██╗"
echo "  ╚██████╔╝██║  ██║╚██████╔╝███████║   ██║   ╚███╔███╔╝██║  ██║██║   ██║   ███████╗██║  ██║"
echo "   ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝   ╚═╝    ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝"
echo ""

# Check Go
if ! command -v go &>/dev/null; then
    echo "[ERROR] Go не найден. Установи Go 1.21+: https://go.dev/dl/"
    exit 1
fi

# Build & install
echo "[INFO] Собираю $BINARY_NAME..."
cd "$SCRIPT_DIR"
go build -o "$BINARY_NAME" ./cmd/ghostwriter/

mkdir -p "$BIN_DIR"
echo "[INFO] Устанавливаю $BINARY_NAME -> $BIN_DIR/$BINARY_NAME"
mv "$BINARY_NAME" "$BIN_DIR/$BINARY_NAME"
chmod +x "$BIN_DIR/$BINARY_NAME"

echo ""
echo "[SUCCESS] ghostwriter установлен!"
echo ""
echo "  Использование:"
echo "    ghostwriter               # случайный язык, средняя скорость"
echo "    ghostwriter --lang python # только Python"
echo "    ghostwriter --lang go     # только Go"
echo "    ghostwriter --lang js     # только JavaScript/TypeScript"
echo "    ghostwriter --speed fast  # быстрее (slow/medium/fast/turbo)"
echo "    Ctrl+C                    # остановить"
echo ""
