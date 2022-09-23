clear

echo Building Linux edition
env GOOS=linux GOARCH=amd64 go build -o builds/linux/ng

echo Building Windows edition
env GOOS=windows GOARCH=amd64 go build -o builds/windows/ng.exe

echo Building MacOS edition - Intel
env GOOS=darwin GOARCH=amd64 go build -o builds/macos-intel/ng

echo Building MacOS edition - ARM - Apple Silicon
env GOOS=darwin GOARCH=arm64 go build -o builds/macos-arm/ng
