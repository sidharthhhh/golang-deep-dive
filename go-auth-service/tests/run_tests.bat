@echo off
REM Go Auth Service - Test Runner Script (Windows)

echo =========================================
echo Go Auth Service - Running Tests
echo =========================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed
    exit /b 1
)

REM Navigate to project root
cd /d "%~dp0\.."

REM Install dependencies
echo Installing dependencies...
go mod download
go get github.com/stretchr/testify/assert

REM Run tests
echo.
echo Running tests...
echo =========================================

if "%1"=="coverage" (
    echo Running tests with coverage...
    go test ./tests/... -v -cover -coverprofile=coverage.out
    
    if %ERRORLEVEL% EQU 0 (
        echo.
        echo [32m✓ All tests passed![0m
        echo.
        echo Generating coverage report...
        go tool cover -html=coverage.out -o coverage.html
        echo [32mCoverage report generated: coverage.html[0m
    ) else (
        echo.
        echo [31m✗ Some tests failed[0m
        exit /b 1
    )
) else if "%1"=="verbose" (
    echo Running tests in verbose mode...
    go test ./tests/... -v
    
    if %ERRORLEVEL% EQU 0 (
        echo.
        echo [32m✓ All tests passed![0m
    ) else (
        echo.
        echo [31m✗ Some tests failed[0m
        exit /b 1
    )
) else (
    go test ./tests/... -v
    
    if %ERRORLEVEL% EQU 0 (
        echo.
        echo [32m✓ All tests passed![0m
    ) else (
        echo.
        echo [31m✗ Some tests failed[0m
        exit /b 1
    )
)

echo.
echo =========================================
echo Test run complete!
echo =========================================
