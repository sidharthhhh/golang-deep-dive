#!/bin/bash

# Go Auth Service - Test Runner Script (Linux/Mac)

echo "========================================="
echo "Go Auth Service - Running Tests"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Check if MySQL is running
if ! command -v mysql &> /dev/null; then
    echo -e "${YELLOW}Warning: MySQL client not found. Skipping database check.${NC}"
else
    echo "Checking MySQL connection..."
    if mysql -u root -phjkl -e "SELECT 1" &> /dev/null; then
        echo -e "${GREEN}✓ MySQL is running${NC}"
    else
        echo -e "${RED}✗ Cannot connect to MySQL${NC}"
        echo "Please ensure MySQL is running and credentials are correct"
        exit 1
    fi
fi

# Navigate to project root
cd "$(dirname "$0")/.."

# Install dependencies
echo ""
echo "Installing dependencies..."
go mod download
go get github.com/stretchr/testify/assert

# Run tests
echo ""
echo "Running tests..."
echo "========================================="

if [ "$1" == "coverage" ]; then
    echo "Running tests with coverage..."
    go test ./tests/... -v -cover -coverprofile=coverage.out
    
    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ All tests passed!${NC}"
        echo ""
        echo "Generating coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        echo -e "${GREEN}Coverage report generated: coverage.html${NC}"
    else
        echo ""
        echo -e "${RED}✗ Some tests failed${NC}"
        exit 1
    fi
elif [ "$1" == "verbose" ]; then
    echo "Running tests in verbose mode..."
    go test ./tests/... -v
    
    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ All tests passed!${NC}"
    else
        echo ""
        echo -e "${RED}✗ Some tests failed${NC}"
        exit 1
    fi
else
    go test ./tests/... -v
    
    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ All tests passed!${NC}"
    else
        echo ""
        echo -e "${RED}✗ Some tests failed${NC}"
        exit 1
    fi
fi

echo ""
echo "========================================="
echo "Test run complete!"
echo "========================================="
