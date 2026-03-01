$ErrorActionPreference = "Stop"

# Build Windows binary
New-Item -ItemType Directory -Force -Path "dist" | Out-Null
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o dist/server-windows-amd64.exe ./

Write-Host "Built dist/server-windows-amd64.exe"
