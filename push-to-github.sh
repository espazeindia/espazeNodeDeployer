#!/bin/bash

# Script to create private GitHub repository and push code
# Run this AFTER authenticating with: gh auth login

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${PURPLE}"
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                                                           ║"
echo "║     Espaze Node Deployer - GitHub Push Script            ║"
echo "║                                                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""

# Check if authenticated
echo -e "${BLUE}🔍 Checking GitHub authentication...${NC}"
if ! gh auth status &>/dev/null; then
    echo -e "${RED}❌ Not authenticated with GitHub${NC}"
    echo ""
    echo -e "${YELLOW}Please run:${NC} gh auth login"
    echo ""
    echo "Then run this script again."
    exit 1
fi

echo -e "${GREEN}✅ GitHub authentication verified${NC}"
echo ""

# Use espazeindia as the organization/username
USERNAME="espazeindia"
echo -e "${BLUE}📝 GitHub Organization: ${GREEN}${USERNAME}${NC}"
echo ""

# Check if repo already exists
REPO_NAME="espazeNodeDeployer"
echo -e "${BLUE}🔍 Checking if repository exists...${NC}"

if gh repo view "${USERNAME}/${REPO_NAME}" &>/dev/null; then
    echo -e "${YELLOW}⚠️  Repository '${REPO_NAME}' already exists${NC}"
    echo ""
    read -p "Do you want to push to the existing repository? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Cancelled."
        exit 0
    fi
    
    # Add remote if not exists
    if ! git remote get-url origin &>/dev/null; then
        git remote add origin "https://github.com/${USERNAME}/${REPO_NAME}.git"
    fi
else
    echo -e "${BLUE}📦 Creating private repository '${REPO_NAME}'...${NC}"
    
    # Create repository under espazeindia organization
    gh repo create "${USERNAME}/${REPO_NAME}" \
        --private \
        --source=. \
        --description="Kubernetes deployment management platform with centralized node tracking, GitHub integration, and real-time observability. Built with Go + React + MongoDB." \
        --disable-wiki
    
    echo -e "${GREEN}✅ Repository created successfully!${NC}"
fi

echo ""
echo -e "${BLUE}🚀 Pushing to GitHub...${NC}"

# Push to GitHub
git push -u origin master

echo ""
echo -e "${GREEN}"
echo "┌─────────────────────────────────────────────────────────┐"
echo "│                                                         │"
echo "│  ✅ Successfully pushed to GitHub!                     │"
echo "│                                                         │"
echo "└─────────────────────────────────────────────────────────┘"
echo -e "${NC}"
echo ""

# Show repository info
echo -e "${BLUE}📊 Repository Information:${NC}"
echo "   🔗 URL: https://github.com/${USERNAME}/${REPO_NAME}"
echo "   🔒 Visibility: Private"
echo "   🌿 Branch: master"
echo "   📦 Files: 64"
echo "   📝 Commits: $(git rev-list --count HEAD)"
echo ""

echo -e "${PURPLE}🎉 Next Steps:${NC}"
echo "  1. View repository: gh repo view --web"
echo "  2. Clone elsewhere: git clone https://github.com/${USERNAME}/${REPO_NAME}.git"
echo "  3. Add collaborators: gh repo edit --add-collaborator USERNAME"
echo "  4. Enable GitHub Actions for CI/CD"
echo ""

echo -e "${GREEN}Happy coding! 🚀${NC}"
echo ""

