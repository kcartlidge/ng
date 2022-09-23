clear
@echo off

@echo Building Linux edition
@set GOOS=linux
@set GOARCH=amd64
@go build -o builds/linux/ng

@echo Building Mac edition - Intel
@set GOOS=darwin
@set GOARCH=amd64
@go build -o builds/macos-intel/ng

@echo Building Mac edition - ARM - Apple Silicon
@set GOOS=darwin
@set GOARCH=arm64
@go build -o builds/macos-arm/ng

@echo Building Windows edition
@set GOOS=windows
@set GOARCH=amd64
@go build -o builds/windows/ng.exe
