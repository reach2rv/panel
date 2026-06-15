#!/bin/bash
# OrnaVerse Panel Install Script (Using GitHub Releases)
# This script downloads the latest pre-compiled binaries from your GitHub Releases
# and configures the systemd service to run the panel.

set -e

# Configuration
# Change this to your actual GitHub username/repo
REPO=${REPO:-"reach2rv/panel"}
INSTALL_DIR="/opt/ornaverse/panel"
TMP_DIR="/tmp/ornaverse_install"

echo "================================================="
echo "   OrnaVerse Panel Installation Script (via Releases)   "
echo "================================================="

# 1. Check Root
if [ "$EUID" -ne 0 ]; then
  echo "Error: Please run this script as root."
  exit 1
fi

# 2. Check Prerequisites
echo "==> Checking prerequisites..."
command -v curl >/dev/null 2>&1 || { echo >&2 "Error: 'curl' is required but not installed."; exit 1; }
command -v unzip >/dev/null 2>&1 || { echo >&2 "Error: 'unzip' is required but not installed."; exit 1; }

# 3. Determine Architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: Unsupported architecture $ARCH"
        exit 1
        ;;
esac

# 4. Fetch Latest Release URL
echo "==> Fetching latest release info from $REPO..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")
VERSION=$(echo "$LATEST_RELEASE" | grep -m 1 '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Error: Could not find the latest release for $REPO. Have you published a release tag yet?"
    exit 1
fi

# Strip leading 'v' from version if present
CLEAN_VERSION="${VERSION#v}"
FILENAME="ornaverse-panel_${CLEAN_VERSION}_linux_${ARCH}.zip"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "==> Found version $VERSION for $ARCH"
echo "==> Downloading $DOWNLOAD_URL..."

# 5. Download and Extract
rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR"
cd "$TMP_DIR"

curl -sSL "$DOWNLOAD_URL" -o "$FILENAME"
unzip -q "$FILENAME"

# 6. Install to /opt/ornaverse
echo "==> Installing to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
cp ace "$INSTALL_DIR/"
cp cli "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/ace" "$INSTALL_DIR/cli"

# Optional storage directory if packaged by goreleaser
if [ -d "storage" ]; then
    cp -r storage "$INSTALL_DIR/"
fi

# Create default configuration if it doesn't exist
if [ ! -f "$INSTALL_DIR/config.yml" ]; then
    echo "==> Creating default configuration..."
    if [ -f "config.example.yml" ]; then
        cp config.example.yml "$INSTALL_DIR/config.yml"
    else
        echo "Warning: config.example.yml not found in release archive."
    fi
else
    echo "==> Existing config.yml found, skipping overwrite."
fi

# Cleanup
rm -rf "$TMP_DIR"

# 7. Setup Systemd Service
echo "==> Setting up systemd service..."
cat > /etc/systemd/system/ornaverse.service <<EOF
[Unit]
Description=OrnaVerse Panel Service
After=network.target

[Service]
Type=simple
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/ace
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable ornaverse
systemctl restart ornaverse

# 8. Create symlinks for CLI tool
echo "==> Creating system symlinks..."
ln -sf "$INSTALL_DIR/cli" /usr/local/sbin/acepanel
ln -sf "$INSTALL_DIR/cli" /usr/local/sbin/ornaverse
chmod +x /usr/local/sbin/acepanel /usr/local/sbin/ornaverse

# 9. Initialize database/settings and print connection details
echo "==> Initializing panel database and settings..."
/usr/local/sbin/ornaverse init
/usr/local/sbin/ornaverse info

echo "================================================="
echo "==> Installation complete!"
echo "==> OrnaVerse Panel is now installed and running."
echo "==> Check status with: systemctl status ornaverse"
echo "==> Manage the panel using the 'ornaverse' command."
echo "================================================="
