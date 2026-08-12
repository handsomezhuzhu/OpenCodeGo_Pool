@echo off
setlocal

set "ROOT=%~dp0"
set "OUT_DIR=%ROOT%build"
set "APP_NAME=opencode-pool"
set "MAIN_PKG=./cmd/server"
set "OUTPUT=%OUT_DIR%\%APP_NAME%-linux-amd64"
set "WEB_DIR=%ROOT%web"
set "FRONTEND_DIST=%ROOT%internal\frontend\dist"

pushd "%ROOT%" >nul
if errorlevel 1 (
    echo ERROR: Failed to enter project directory
    exit /b 1
)

if not exist "%OUT_DIR%" (
    mkdir "%OUT_DIR%"
    if errorlevel 1 (
        echo ERROR: Failed to create build directory
        popd >nul
        exit /b 1
    )
)

echo Installing frontend dependencies...
cd /d "%WEB_DIR%"
call npm ci
if errorlevel 1 (
    echo ERROR: npm ci failed
    cd /d "%ROOT%"
    popd >nul
    exit /b 1
)

echo Building frontend...
call npm run build
if errorlevel 1 (
    echo ERROR: Frontend build failed
    cd /d "%ROOT%"
    popd >nul
    exit /b 1
)
cd /d "%ROOT%"

echo Copying frontend build to internal\frontend\dist...
if exist "%FRONTEND_DIST%" (
    rmdir /s /q "%FRONTEND_DIST%"
    if errorlevel 1 (
        echo ERROR: Failed to remove old frontend dist
        popd >nul
        exit /b 1
    )
)
xcopy "%WEB_DIR%\dist" "%FRONTEND_DIST%\" /e /i /y >nul
if errorlevel 1 (
    echo ERROR: Failed to copy frontend build
    popd >nul
    exit /b 1
)

echo Running tests...
go test ./...
if errorlevel 1 (
    echo ERROR: Go tests failed
    popd >nul
    exit /b 1
)

echo Running go vet...
go vet ./...
if errorlevel 1 (
    echo ERROR: go vet failed
    popd >nul
    exit /b 1
)

echo Building Go binary for Linux amd64...
set "GOOS=linux"
set "GOARCH=amd64"
set "CGO_ENABLED=0"
go build -v -trimpath -ldflags="-s -w" -o "%OUTPUT%" "%MAIN_PKG%"
if errorlevel 1 (
    echo ERROR: Go build failed
    popd >nul
    exit /b 1
)
echo Done: %OUTPUT%

popd >nul
endlocal
