@echo off
echo ========================================================
echo   ARCHIVAL VAULT // MONOSPACED RAW LEDGER & REGISTRY
echo ========================================================

cd /d "%~dp0backend"
if not exist vault-server.exe (
    echo [BUILD] Compiling Go backend binary...
    go build -o vault-server.exe main.go
)

echo [START] Launching Golang REST Server on port 8085...
start "" vault-server.exe

echo [HERD] Opening Laravel Herd Site: http://archival-vault.test
timeout /t 2 >nul
start http://archival-vault.test

echo [READY] Server active. Press any key to stop.
pause
