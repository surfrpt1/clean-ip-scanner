#!/bin/bash

set -e

GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}    Clean IP Scanner - Installer        ${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

echo -e "${YELLOW}This installer is for iSH on iPad (Alpine Linux)${NC}"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Go not found. Installing Go...${NC}"
    echo -e "${CYAN}Run these commands in iSH first:${NC}"
    echo ""
    echo "  apk add go"
    echo "  apk add git"
    echo "  apk add build-base"
    echo ""
    echo -e "${RED}After installing Go, run this script again.${NC}"
    exit 1
fi

echo -e "${GREEN}Go found: $(go version)${NC}"

# Create directory structure
echo -e "${CYAN}Setting up project...${NC}"
mkdir -p clean-ip-scanner
cd clean-ip-scanner

# Initialize Go module if needed
if [ ! -f "go.mod" ]; then
    echo -e "${YELLOW}Initializing Go module...${NC}"
    go mod init github.com/surfrpt1/clean-ip-scanner
fi

echo -e "${CYAN}Downloading dependencies...${NC}"
go mod tidy

echo -e "${CYAN}Building...${NC}"
go build -o clean-ip-scanner .

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}Build successful!${NC}"
    echo ""
    echo -e "${CYAN}To run:${NC}"
    echo "  cd clean-ip-scanner"
    echo "  chmod +x clean-ip-scanner"
    echo "  ./clean-ip-scanner -save"
    echo ""
    echo -e "${CYAN}For Xray mode:${NC}"
    echo "  mkdir -p xray"
    echo "  # Download xray binary for Linux AMD64 to ./xray/xray"
    echo "  ./clean-ip-scanner -xray -save"
    echo ""
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
