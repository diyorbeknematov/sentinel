#!/bin/bash
set -e

BINARY_NAME="sentinel"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/sentinel"
CONFIG_FILE="$CONFIG_DIR/config.yaml"
SERVICE_FILE="/etc/systemd/system/sentinel.service"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

log() { echo -e "${GREEN}[INFO]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Root check
[ "$EUID" -ne 0 ] && error "Iltimos sudo bilan ishga tushiring: sudo $0"

# Binary check
[ ! -f "./$BINARY_NAME" ] && error "'$BINARY_NAME' binary topilmadi (current folderda bo‘lishi kerak)"

log "Sentinel o'rnatilmoqda..."

# =========================
# 1. Install binary
# =========================
install -m 755 "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
log "Binary o'rnatildi: $INSTALL_DIR/$BINARY_NAME"

# =========================
# 2. Config setup
# =========================
mkdir -p "$CONFIG_DIR"

if [ ! -f "$CONFIG_FILE" ]; then
    cp config.yaml "$CONFIG_FILE"
    log "Config yaratildi: $CONFIG_FILE"
else
    cp config.yaml "$CONFIG_DIR/config.yaml.new"
    log "Config mavjud, yangi variant saqlandi: config.yaml.new"
fi

# =========================
# 3. Systemd service
# =========================
if [ ! -f "$SERVICE_FILE" ]; then
    cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Sentinel Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/$BINARY_NAME --config $CONFIG_FILE
Restart=always
RestartSec=3
User=root

[Install]
WantedBy=multi-user.target
EOF

    log "Systemd service yaratildi"
else
    log "Systemd service allaqachon mavjud"
fi

# =========================
# 4. Reload + start service
# =========================
systemctl daemon-reload
systemctl enable sentinel
systemctl restart sentinel

log "Service ishga tushirildi"

# =========================
# DONE
# =========================
echo ""
echo "================================="
echo "   ✅ INSTALLATION COMPLETE"
echo "================================="
echo "Binary : $INSTALL_DIR/$BINARY_NAME"
echo "Config : $CONFIG_FILE"
echo "Service: sentinel (systemctl)"
echo ""
echo "Commands:"
echo "  systemctl status sentinel"
echo "  journalctl -u sentinel -f"
echo "================================="