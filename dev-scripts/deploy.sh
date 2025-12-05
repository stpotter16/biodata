#!/usr/bin/env bash

# Fail on first error
set -e

# Fail on unset variable
set -u

# Configuration
NIXOS_CONFIG_DIR="${NIXOS_CONFIG_DIR:-$HOME/development/nixos-config}"
TARGET_HOST="${TARGET_HOST:porphyrion}"
HOSTNAME="${HOSTNAME:-porphyrion}"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}==> Building biodata package...${NC}"
nix build

echo -e "${BLUE}==> Updating flake.lock in nixos-config...${NC}"
cd "$NIXOS_CONFIG_DIR"
nix flake lock --update-input biodata

if [ -z "$TARGET_HOST" ]; then
    echo -e "${YELLOW}No TARGET_HOST set, deploying locally${NC}"
    echo -e "${BLUE}==> Deploying to local system...${NC}"
    sudo nixos-rebuild switch --flake ".#$HOSTNAME"
else
    echo -e "${BLUE}==> Deploying to $TARGET_HOST...${NC}"
    nixos-rebuild switch \
        --flake ".#$HOSTNAME" \
        --target-host "$TARGET_HOST" \
        --build-host localhost \
        --use-remote-sudo
fi

echo -e "${GREEN}==> Deployment complete!${NC}"
