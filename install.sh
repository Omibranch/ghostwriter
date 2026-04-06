#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo ""
echo "  ██████╗ ██╗  ██╗ ██████╗ ███████╗████████╗██╗    ██╗██████╗ ██╗████████╗███████╗██████╗ "
echo "  ██╔════╝██║  ██║██╔═══██╗██╔════╝╚══██╔══╝██║    ██║██╔══██╗██║╚══██╔══╝██╔════╝██╔══██╗"
echo "  ██║  ███╗███████║██║   ██║███████╗   ██║   ██║ █╗ ██║██████╔╝██║   ██║   █████╗  ██████╔╝"
echo "  ██║   ██║██╔══██║██║   ██║╚════██║   ██║   ██║███╗██║██╔══██╗██║   ██║   ██╔══╝  ██╔══██╗"
echo "  ╚██████╔╝██║  ██║╚██████╔╝███████║   ██║   ╚███╔███╔╝██║  ██║██║   ██║   ███████╗██║  ██║"
echo "   ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝   ╚═╝    ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝"
echo ""

if ! command -v python3 &>/dev/null; then
    echo "[ERROR] Python 3 не найден."
    exit 1
fi

PY_VER=$(python3 -c "import sys; print(sys.version_info.major * 10 + sys.version_info.minor)")
if [ "$PY_VER" -lt 39 ]; then
    echo "[ERROR] Требуется Python 3.9+."
    exit 1
fi

echo "[INFO] Устанавливаю ghostwriter..."
pip3 install -e "$SCRIPT_DIR" --quiet

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
